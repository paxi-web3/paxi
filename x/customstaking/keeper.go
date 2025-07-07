package customstaking

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"sort"

	addresscodec "cosmossdk.io/core/address"
	store "cosmossdk.io/core/store"
	sdkmath "cosmossdk.io/math"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	"github.com/cosmos/cosmos-sdk/x/staking/types"
	gogotypes "github.com/cosmos/gogoproto/types"
	"github.com/paxi-web3/paxi/utils"
)

type CustomStakingKeeper struct {
	*stakingkeeper.Keeper
	authKeeper            types.AccountKeeper
	bankKeeper            types.BankKeeper
	storeService          store.KVStoreService
	cdc                   codec.BinaryCodec
	validatorAddressCodec addresscodec.Codec
}

func NewCustomKeeper(baseKeeper *stakingkeeper.Keeper, authKeeper types.AccountKeeper, bankKeeper types.BankKeeper, storeService store.KVStoreService, cdc codec.BinaryCodec, validatorAddressCodec addresscodec.Codec) *CustomStakingKeeper {
	return &CustomStakingKeeper{
		Keeper:                baseKeeper,
		authKeeper:            authKeeper,
		bankKeeper:            bankKeeper,
		storeService:          storeService,
		cdc:                   cdc,
		validatorAddressCodec: validatorAddressCodec,
	}
}

type validatorsByAddr map[string]int64

// EndBlocker called at every block, update validator set
func (k *CustomStakingKeeper) EndBlocker(ctx context.Context) ([]abci.ValidatorUpdate, error) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, telemetry.Now(), telemetry.MetricKeyEndBlocker)
	return k.BlockValidatorUpdates(ctx)
}

// BlockValidatorUpdates calculates the ValidatorUpdates for the current block
// Called in each EndBlock
func (k CustomStakingKeeper) BlockValidatorUpdates(ctx context.Context) ([]abci.ValidatorUpdate, error) {
	// Calculate validator set changes.
	//
	// NOTE: ApplyAndReturnValidatorSetUpdates has to come before
	// UnbondAllMatureValidatorQueue.
	// This fixes a bug when the unbonding period is instant (is the case in
	// some of the tests). The test expected the validator to be completely
	// unbonded after the Endblocker (go from Bonded -> Unbonding during
	// ApplyAndReturnValidatorSetUpdates and then Unbonding -> Unbonded during
	// UnbondAllMatureValidatorQueue).
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	// NOTE: ApplyAndReturnValidatorSetUpdates is now Paxi's custom function
	validatorUpdates, err := k.ApplyAndReturnValidatorSetUpdates(sdkCtx)
	if err != nil {
		return nil, err
	}

	// unbond all mature validators from the unbonding queue
	err = k.UnbondAllMatureValidators(ctx)
	if err != nil {
		return nil, err
	}

	// Remove all mature unbonding delegations from the ubd queue.
	matureUnbonds, err := k.DequeueAllMatureUBDQueue(ctx, sdkCtx.BlockHeader().Time)
	if err != nil {
		return nil, err
	}

	for _, dvPair := range matureUnbonds {
		addr, err := k.Keeper.ValidatorAddressCodec().StringToBytes(dvPair.ValidatorAddress)
		if err != nil {
			return nil, err
		}
		delegatorAddress, err := k.authKeeper.AddressCodec().StringToBytes(dvPair.DelegatorAddress)
		if err != nil {
			return nil, err
		}

		balances, err := k.CompleteUnbonding(ctx, delegatorAddress, addr)
		if err != nil {
			continue
		}

		sdkCtx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeCompleteUnbonding,
				sdk.NewAttribute(sdk.AttributeKeyAmount, balances.String()),
				sdk.NewAttribute(types.AttributeKeyValidator, dvPair.ValidatorAddress),
				sdk.NewAttribute(types.AttributeKeyDelegator, dvPair.DelegatorAddress),
			),
		)
	}

	// Remove all mature redelegations from the red queue.
	matureRedelegations, err := k.DequeueAllMatureRedelegationQueue(ctx, sdkCtx.BlockHeader().Time)
	if err != nil {
		return nil, err
	}

	for _, dvvTriplet := range matureRedelegations {
		valSrcAddr, err := k.ValidatorAddressCodec().StringToBytes(dvvTriplet.ValidatorSrcAddress)
		if err != nil {
			return nil, err
		}
		valDstAddr, err := k.ValidatorAddressCodec().StringToBytes(dvvTriplet.ValidatorDstAddress)
		if err != nil {
			return nil, err
		}
		delegatorAddress, err := k.authKeeper.AddressCodec().StringToBytes(dvvTriplet.DelegatorAddress)
		if err != nil {
			return nil, err
		}

		balances, err := k.CompleteRedelegation(
			ctx,
			delegatorAddress,
			valSrcAddr,
			valDstAddr,
		)
		if err != nil {
			continue
		}

		sdkCtx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeCompleteRedelegation,
				sdk.NewAttribute(sdk.AttributeKeyAmount, balances.String()),
				sdk.NewAttribute(types.AttributeKeyDelegator, dvvTriplet.DelegatorAddress),
				sdk.NewAttribute(types.AttributeKeySrcValidator, dvvTriplet.ValidatorSrcAddress),
				sdk.NewAttribute(types.AttributeKeyDstValidator, dvvTriplet.ValidatorDstAddress),
			),
		)
	}

	return validatorUpdates, nil
}

// In the default Cosmos SDK staking module, ApplyAndReturnValidatorSetUpdates()
// selects validators purely based on their staking amount, always choosing the
// top "MaxValidators" sorted by voting power.
//
// However, to improve decentralization, fairness, and network security, we
// implement a custom validator selection mechanism:
//
//   - Every specific blocks, the validator set is refreshed.
//   - 50% of the validators are selected based on the highest voting power (top N).
//   - The remaining 50% are randomly selected from the rest of the candidates,
//     using weighted random sampling proportional to their staking amount.
//
// This design maintains consensus safety (by ensuring a strong voting power
// foundation) while giving smaller validators a fair chance to participate,
// thereby improving overall network decentralization and resilience.
//
// By overriding ApplyAndReturnValidatorSetUpdates, we achieve custom validator
// rotation logic without forking the entire staking module, making future
// upgrades easier and keeping the system modular and maintainable.
func (k CustomStakingKeeper) ApplyAndReturnValidatorSetUpdates(ctx sdk.Context) ([]abci.ValidatorUpdate, error) {
	params, err := k.GetParams(ctx)
	if err != nil {
		return nil, err
	}
	maxValidators := params.MaxValidators
	powerReduction := k.PowerReduction(ctx)
	amtFromBondedToNotBonded, amtFromNotBondedToBonded := sdkmath.ZeroInt(), sdkmath.ZeroInt()
	totalPower := sdkmath.ZeroInt()
	var updates []abci.ValidatorUpdate

	// Get all the validators
	allValidators, err := k.getValidatorsAboveThreshold(ctx, sdkmath.NewInt(MinBondedTokens), MaxCandidates)
	if err != nil {
		return nil, fmt.Errorf("failed to get last validator set: %w", err)
	}

	// Retrieve the last validator set.
	// The persistent set is updated later in this function.
	// (see LastValidatorPowerKey).
	last, err := k.getLastValidatorsByAddr(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get last validator set: %w", err)
	}

	// Update validator set every specific blocks
	totalCandidates := uint32(len(allValidators))
	currentBlock := ctx.BlockHeight()
	if currentBlock%BlocksPerUpdate == 0 || currentBlock <= 1 {
		newValidators := []types.Validator{}
		if totalCandidates == 1 {
			newValidators = append(newValidators, allValidators[:totalCandidates]...)
		} else {
			// Top 50% will certainly be selected
			maxValidators = sdkmath.Min(maxValidators, totalCandidates)
			splitter := maxValidators / 2
			newValidators = append(newValidators, allValidators[0:splitter]...)

			// The another 50% validators will be selected randomly from the remaining candidates
			remainingCandidates := allValidators[splitter:]
			appHash := ctx.BlockHeader().AppHash
			if len(appHash) < 8 {
				panic("AppHash too short for seed")
			}
			seed := binary.BigEndian.Uint64(appHash) ^ uint64(ctx.BlockHeader().Time.UnixNano())
			r := rand.New(rand.NewSource(int64(seed)))
			newValidators = append(newValidators, pickWeightedRandomSubsetBinarySearch(r, remainingCandidates, int(maxValidators)-len(newValidators))...)
		}

		// Add to updates
		for _, val := range newValidators {
			// apply the appropriate state change if necessary
			switch {
			case val.IsUnbonded():
				val, err = k.unbondedToBonded(ctx, val)
				if err != nil {
					return nil, err
				}
				amtFromNotBondedToBonded = amtFromNotBondedToBonded.Add(val.GetTokens())
			case val.IsUnbonding():
				val, err = k.unbondingToBonded(ctx, val)
				if err != nil {
					return nil, err
				}
				amtFromNotBondedToBonded = amtFromNotBondedToBonded.Add(val.GetTokens())
			case val.IsBonded():
				// no state change
			default:
				return nil, errors.New("unexpected validator status")
			}

			valAddr, err := sdk.ValAddressFromBech32(val.OperatorAddress)
			if err != nil {
				return nil, err
			}
			valAddrStr := string(valAddr)

			// fetch the old power bytes
			oldPower, found := last[valAddrStr]
			newPower := val.ConsensusPower(powerReduction)

			// update the validator set if power has changed
			if !found || oldPower != newPower {
				update := val.ABCIValidatorUpdate(powerReduction)
				updates = append(updates, update)

				if err = k.SetLastValidatorPower(ctx, valAddr, newPower); err != nil {
					return nil, err
				}
			}

			// Delete from the old list
			delete(last, valAddrStr)

			// Add total power
			totalPower = totalPower.AddRaw(newPower)
		}
	} else {
		// Kick who is jailed or unbonded
		addresses := make([]string, 0, len(last))
		for k := range last {
			addresses = append(addresses, k)
		}
		for idx := range addresses {
			addr := addresses[idx]
			valAddr := sdk.ValAddress([]byte(addr))
			val, err := k.GetValidator(ctx, valAddr)
			if err != nil {
				return nil, fmt.Errorf("validator record not found for address: %X", valAddr)
			}

			// Add to updates
			if val.IsBonded() && !val.Jailed {
				oldPower, found := last[addr]
				newPower := val.ConsensusPower(powerReduction)

				// update the validator set if power has changed
				if !found || oldPower != newPower {
					update := val.ABCIValidatorUpdate(powerReduction)
					updates = append(updates, update)

					if err = k.SetLastValidatorPower(ctx, valAddr, newPower); err != nil {
						return nil, err
					}
				}

				// Delete from the old list
				delete(last, addr)

				// Add total power
				totalPower = totalPower.AddRaw(newPower)
			}
		}
	}

	// Process last validators which will be unboned
	noLongerBonded, err := sortNoLongerBonded(last, k.validatorAddressCodec)
	if err != nil {
		return nil, err
	}

	for _, valAddrBytes := range noLongerBonded {
		validator, err := k.GetValidator(ctx, sdk.ValAddress(valAddrBytes))
		if err != nil {
			return nil, fmt.Errorf("validator record not found for address: %X", sdk.ValAddress(valAddrBytes))
		}
		validator, err = k.bondedToUnbonding(ctx, validator)
		if err != nil {
			return nil, err
		}
		str, err := k.validatorAddressCodec.StringToBytes(validator.GetOperator())
		if err != nil {
			return nil, fmt.Errorf("failed to get validator operator address: %w", err)
		}
		amtFromBondedToNotBonded = amtFromBondedToNotBonded.Add(validator.GetTokens())
		if err = k.DeleteLastValidatorPower(ctx, str); err != nil {
			return nil, err
		}

		updates = append(updates, validator.ABCIValidatorUpdateZero())
	}

	// Update the pools based on the recent updates in the validator set:
	// - The tokens from the non-bonded candidates that enter the new validator set need to be transferred
	// to the Bonded pool.
	// - The tokens from the bonded validators that are being kicked out from the validator set
	// need to be transferred to the NotBonded pool.
	switch {
	// Compare and subtract the respective amounts to only perform one transfer.
	// This is done in order to avoid doing multiple updates inside each iterator/loop.
	case amtFromNotBondedToBonded.GT(amtFromBondedToNotBonded):
		if err = k.notBondedTokensToBonded(ctx, amtFromNotBondedToBonded.Sub(amtFromBondedToNotBonded)); err != nil {
			return nil, err
		}
	case amtFromNotBondedToBonded.LT(amtFromBondedToNotBonded):
		if err = k.bondedTokensToNotBonded(ctx, amtFromBondedToNotBonded.Sub(amtFromNotBondedToBonded)); err != nil {
			return nil, err
		}
	default: // equal amounts of tokens; no update required
	}

	if len(updates) > 0 || len(last) != 0 {
		if err := k.SetLastTotalPower(ctx, totalPower); err != nil {
			return nil, err
		}
	}

	if err := k.SetValidatorUpdates(ctx, updates); err != nil {
		return nil, err
	}

	return updates, nil
}

func (k CustomStakingKeeper) getValidatorsAboveThreshold(ctx sdk.Context, minTokens sdkmath.Int, maxCount int) ([]types.Validator, error) {
	var validators []types.Validator

	// read validator from the store
	store := k.storeService.OpenKVStore(ctx)

	start, end := utils.PrefixRange(types.ValidatorsByPowerIndexKey)
	iterator, err := store.ReverseIterator(start, end)
	if err != nil {
		return nil, err
	}
	defer iterator.Close()

	vcount := 0
	for ; iterator.Valid(); iterator.Next() {
		valAddr := sdk.ValAddress(iterator.Value())
		validator, err := k.GetValidator(ctx, valAddr)
		if err != nil {
			return nil, fmt.Errorf("validator record not found for address: %X", valAddr)
		}

		if validator.Jailed {
			return nil, errors.New("should never retrieve a jailed validator from the power store")
		}

		// Elimate validator that lower than minimum threshold
		if validator.Tokens.LT(minTokens) {
			break
		}
		validators = append(validators, validator)

		// Check if the validator set reached maximum
		vcount += 1
		if vcount >= maxCount {
			break
		}
	}

	return validators, nil
}

func pickWeightedRandomSubsetBinarySearch(r *rand.Rand, candidates []types.Validator, count int) []types.Validator {
	// Build prefix sum
	n := len(candidates)
	prefixSums := make([]int64, n)
	var total int64 = 0
	for i, val := range candidates {
		w := int64(math.Sqrt(float64(val.Tokens.Int64())))
		if w < 0 {
			w = 0 // safeguard
		}
		total += w
		prefixSums[i] = total
	}

	// Result
	selected := make(map[int]bool)
	result := make([]types.Validator, 0, count)

	for len(result) < count {
		randWeight := r.Int63n(total)
		// Binary search
		i := sort.Search(n, func(i int) bool { return prefixSums[i] > randWeight })

		// Skip if already selected
		if selected[i] {
			continue
		}
		selected[i] = true
		result = append(result, candidates[i])
	}

	return result
}

// get the last validator set
func (k CustomStakingKeeper) getLastValidatorsByAddr(ctx context.Context) (validatorsByAddr, error) {
	last := make(validatorsByAddr)

	iterator, err := k.LastValidatorsIterator(ctx)
	if err != nil {
		return nil, err
	}
	defer iterator.Close()

	var intVal gogotypes.Int64Value
	for ; iterator.Valid(); iterator.Next() {
		// extract the validator address from the key (prefix is 1-byte, addrLen is 1-byte)
		valAddrStr := string(types.AddressFromLastValidatorPowerKey(iterator.Key()))
		k.cdc.MustUnmarshal(iterator.Value(), &intVal)
		last[valAddrStr] = intVal.GetValue()
	}

	return last, nil
}

func (k CustomStakingKeeper) bondedToUnbonding(ctx context.Context, validator types.Validator) (types.Validator, error) {
	if !validator.IsBonded() {
		return types.Validator{}, fmt.Errorf("bad state transition bondedToUnbonding, validator: %v", validator)
	}

	return k.BeginUnbondingValidator(ctx, validator)
}

func (k CustomStakingKeeper) unbondingToBonded(ctx context.Context, validator types.Validator) (types.Validator, error) {
	if !validator.IsUnbonding() {
		return types.Validator{}, fmt.Errorf("bad state transition unbondingToBonded, validator: %v", validator)
	}

	return k.bondValidator(ctx, validator)
}

func (k CustomStakingKeeper) unbondedToBonded(ctx context.Context, validator types.Validator) (types.Validator, error) {
	if !validator.IsUnbonded() {
		return types.Validator{}, fmt.Errorf("bad state transition unbondedToBonded, validator: %v", validator)
	}

	return k.bondValidator(ctx, validator)
}

// perform all the store operations for when a validator status becomes bonded
func (k CustomStakingKeeper) bondValidator(ctx context.Context, validator types.Validator) (types.Validator, error) {
	// delete the validator by power index, as the key will change
	if err := k.DeleteValidatorByPowerIndex(ctx, validator); err != nil {
		return types.Validator{}, err
	}

	validator = validator.UpdateStatus(types.Bonded)

	// save the now bonded validator record to the two referenced stores
	if err := k.SetValidator(ctx, validator); err != nil {
		return types.Validator{}, err
	}

	if err := k.SetValidatorByPowerIndex(ctx, validator); err != nil {
		return types.Validator{}, err
	}

	// delete from queue if present
	if err := k.DeleteValidatorQueue(ctx, validator); err != nil {
		return types.Validator{}, err
	}

	// trigger hook
	consAddr, err := validator.GetConsAddr()
	if err != nil {
		return types.Validator{}, err
	}

	str, err := k.validatorAddressCodec.StringToBytes(validator.GetOperator())
	if err != nil {
		return types.Validator{}, fmt.Errorf("failed to get validator operator address: %w", err)
	}

	if err := k.Hooks().AfterValidatorBonded(ctx, consAddr, str); err != nil {
		return types.Validator{}, err
	}

	return validator, nil
}

// given a map of remaining validators to previous bonded power
// returns the list of validators to be unbonded, sorted by operator address
func sortNoLongerBonded(last validatorsByAddr, _ addresscodec.Codec) ([][]byte, error) {
	// sort the map keys for determinism
	noLongerBonded := make([][]byte, len(last))
	index := 0

	for valAddrStr := range last {
		valAddrBytes := []byte(valAddrStr)
		noLongerBonded[index] = valAddrBytes
		index++
	}
	// sorted by address - order doesn't matter
	sort.SliceStable(noLongerBonded, func(i, j int) bool {
		// -1 means strictly less than
		return bytes.Compare(noLongerBonded[i], noLongerBonded[j]) == -1
	})

	return noLongerBonded, nil
}

// bondedTokensToNotBonded transfers coins from the bonded to the not bonded pool within staking
func (k CustomStakingKeeper) bondedTokensToNotBonded(ctx context.Context, tokens sdkmath.Int) error {
	bondDenom, err := k.BondDenom(ctx)
	if err != nil {
		return err
	}

	coins := sdk.NewCoins(sdk.NewCoin(bondDenom, tokens))
	return k.bankKeeper.SendCoinsFromModuleToModule(ctx, types.BondedPoolName, types.NotBondedPoolName, coins)
}

// notBondedTokensToBonded transfers coins from the not bonded to the bonded pool within staking
func (k CustomStakingKeeper) notBondedTokensToBonded(ctx context.Context, tokens sdkmath.Int) error {
	bondDenom, err := k.BondDenom(ctx)
	if err != nil {
		return err
	}

	coins := sdk.NewCoins(sdk.NewCoin(bondDenom, tokens))
	return k.bankKeeper.SendCoinsFromModuleToModule(ctx, types.NotBondedPoolName, types.BondedPoolName, coins)
}
