package keeper

import (
	"encoding/json"
	"fmt"

	sdkmath "cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/paxi-web3/paxi/x/prediction/types"
)

type prc20AllowanceQuery struct {
	Allowance struct {
		Owner   string `json:"owner"`
		Spender string `json:"spender"`
	} `json:"allowance"`
}

type prc20AllowanceResponse struct {
	Allowance string `json:"allowance"`
}

type prc20BalanceQuery struct {
	Balance struct {
		Address string `json:"address"`
	} `json:"balance"`
}

type prc20BalanceResponse struct {
	Balance string `json:"balance"`
}

func (k Keeper) transferCollateralBetweenAccounts(ctx sdk.Context, market *types.Market, from sdk.AccAddress, to sdk.AccAddress, amount sdkmath.Int) error {
	if !amount.IsPositive() {
		return nil
	}

	switch market.CollateralType {
	case types.CollateralType_COLLATERAL_TYPE_NATIVE:
		coin := sdk.NewCoin(market.CollateralDenom, amount)
		return k.bankKeeper.SendCoins(ctx, from, to, sdk.NewCoins(coin))
	case types.CollateralType_COLLATERAL_TYPE_PRC20:
		return k.transferPRC20FromUserToUser(ctx, from, to, market.CollateralContractAddr, amount)
	default:
		return fmt.Errorf("unsupported collateral type")
	}
}

func (k Keeper) transferCollateralFromModule(ctx sdk.Context, market *types.Market, to sdk.AccAddress, amount sdkmath.Int) error {
	if !amount.IsPositive() {
		return nil
	}

	switch market.CollateralType {
	case types.CollateralType_COLLATERAL_TYPE_NATIVE:
		coin := sdk.NewCoin(market.CollateralDenom, amount)
		return k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, to, sdk.NewCoins(coin))
	case types.CollateralType_COLLATERAL_TYPE_PRC20:
		return k.transferPRC20FromModule(ctx, to, market.CollateralContractAddr, amount)
	default:
		return fmt.Errorf("unsupported collateral type")
	}
}

func (k Keeper) transferPRC20FromUserToUser(ctx sdk.Context, from sdk.AccAddress, to sdk.AccAddress, contract string, amount sdkmath.Int) error {
	if k.prc20Keeper == nil {
		return fmt.Errorf("prc20 keeper is not configured")
	}
	if !amount.IsPositive() {
		return nil
	}

	type transferFrom struct {
		Owner     string `json:"owner"`
		Recipient string `json:"recipient"`
		Amount    string `json:"amount"`
	}
	type msgWrapper struct {
		TransferFrom transferFrom `json:"transfer_from"`
	}

	contractAddr, err := sdk.AccAddressFromBech32(contract)
	if err != nil {
		return fmt.Errorf("invalid prc20 contract address: %w", err)
	}
	moduleAddr := k.accountKeeper.GetModuleAddress(types.ModuleName)
	msg := msgWrapper{
		TransferFrom: transferFrom{
			Owner:     from.String(),
			Recipient: to.String(),
			Amount:    amount.String(),
		},
	}

	bz, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal transfer_from message: %w", err)
	}

	const safeGas uint64 = 10_000_000
	if ctx.IsCheckTx() || ctx.IsReCheckTx() {
		_, err = k.prc20Keeper.Execute(ctx, contractAddr, moduleAddr, bz, nil)
		return err
	}

	parent := ctx.GasMeter()
	child := storetypes.NewGasMeter(safeGas)
	execCtx := ctx.WithGasMeter(child)

	_, err = k.prc20Keeper.Execute(execCtx, contractAddr, moduleAddr, bz, nil)
	parent.ConsumeGas(child.GasConsumed(), "prediction prc20 transfer_from")
	return err
}

func (k Keeper) transferPRC20FromModule(ctx sdk.Context, to sdk.AccAddress, contract string, amount sdkmath.Int) error {
	if k.prc20Keeper == nil {
		return fmt.Errorf("prc20 keeper is not configured")
	}
	if !amount.IsPositive() {
		return nil
	}

	type transfer struct {
		Recipient string `json:"recipient"`
		Amount    string `json:"amount"`
	}
	type msgWrapper struct {
		Transfer transfer `json:"transfer"`
	}

	contractAddr, err := sdk.AccAddressFromBech32(contract)
	if err != nil {
		return fmt.Errorf("invalid prc20 contract address: %w", err)
	}
	moduleAddr := k.accountKeeper.GetModuleAddress(types.ModuleName)
	msg := msgWrapper{Transfer: transfer{Recipient: to.String(), Amount: amount.String()}}

	bz, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal transfer message: %w", err)
	}

	const safeGas uint64 = 10_000_000
	if ctx.IsCheckTx() || ctx.IsReCheckTx() {
		_, err = k.prc20Keeper.Execute(ctx, contractAddr, moduleAddr, bz, nil)
		return err
	}

	parent := ctx.GasMeter()
	child := storetypes.NewGasMeter(safeGas)
	execCtx := ctx.WithGasMeter(child)

	_, err = k.prc20Keeper.Execute(execCtx, contractAddr, moduleAddr, bz, nil)
	parent.ConsumeGas(child.GasConsumed(), "prediction prc20 transfer")
	return err
}

func (k Keeper) ensurePRC20Allowance(ctx sdk.Context, contract string, owner sdk.AccAddress, required sdkmath.Int) error {
	if !required.IsPositive() {
		return nil
	}
	allowance, err := k.getPRC20Allowance(ctx, contract, owner)
	if err != nil {
		return err
	}
	if allowance.LT(required) {
		return fmt.Errorf("insufficient prc20 allowance: required=%s allowance=%s", required.String(), allowance.String())
	}

	return nil
}

func (k Keeper) ensurePRC20Balance(ctx sdk.Context, contract string, owner sdk.AccAddress, required sdkmath.Int) error {
	if !required.IsPositive() {
		return nil
	}
	balance, err := k.getPRC20Balance(ctx, contract, owner)
	if err != nil {
		return err
	}
	if balance.LT(required) {
		return fmt.Errorf("insufficient prc20 balance: required=%s balance=%s", required.String(), balance.String())
	}

	return nil
}

func (k Keeper) getPRC20Allowance(ctx sdk.Context, contract string, owner sdk.AccAddress) (sdkmath.Int, error) {
	if k.prc20Query == nil {
		return sdkmath.Int{}, fmt.Errorf("prc20 query keeper is not configured")
	}

	contractAddr, err := sdk.AccAddressFromBech32(contract)
	if err != nil {
		return sdkmath.Int{}, fmt.Errorf("invalid prc20 contract address: %w", err)
	}
	moduleAddr := k.accountKeeper.GetModuleAddress(types.ModuleName)

	query := prc20AllowanceQuery{}
	query.Allowance.Owner = owner.String()
	query.Allowance.Spender = moduleAddr.String()
	queryBz, err := json.Marshal(query)
	if err != nil {
		return sdkmath.Int{}, fmt.Errorf("failed to marshal prc20 allowance query: %w", err)
	}

	respBz, err := k.prc20Query.QuerySmart(ctx, contractAddr, queryBz)
	if err != nil {
		return sdkmath.Int{}, fmt.Errorf("failed to query prc20 allowance: %w", err)
	}

	var resp prc20AllowanceResponse
	if err := json.Unmarshal(respBz, &resp); err != nil {
		return sdkmath.Int{}, fmt.Errorf("failed to decode prc20 allowance response: %w", err)
	}
	allowance, err := parseNonNegativeInt(resp.Allowance, "allowance")
	if err != nil {
		return sdkmath.Int{}, fmt.Errorf("invalid prc20 allowance response: %w", err)
	}

	return allowance, nil
}

func (k Keeper) getPRC20Balance(ctx sdk.Context, contract string, owner sdk.AccAddress) (sdkmath.Int, error) {
	if k.prc20Query == nil {
		return sdkmath.Int{}, fmt.Errorf("prc20 query keeper is not configured")
	}

	contractAddr, err := sdk.AccAddressFromBech32(contract)
	if err != nil {
		return sdkmath.Int{}, fmt.Errorf("invalid prc20 contract address: %w", err)
	}

	query := prc20BalanceQuery{}
	query.Balance.Address = owner.String()
	queryBz, err := json.Marshal(query)
	if err != nil {
		return sdkmath.Int{}, fmt.Errorf("failed to marshal prc20 balance query: %w", err)
	}

	respBz, err := k.prc20Query.QuerySmart(ctx, contractAddr, queryBz)
	if err != nil {
		return sdkmath.Int{}, fmt.Errorf("failed to query prc20 balance: %w", err)
	}

	var resp prc20BalanceResponse
	if err := json.Unmarshal(respBz, &resp); err != nil {
		return sdkmath.Int{}, fmt.Errorf("failed to decode prc20 balance response: %w", err)
	}
	balance, err := parseNonNegativeInt(resp.Balance, "balance")
	if err != nil {
		return sdkmath.Int{}, fmt.Errorf("invalid prc20 balance response: %w", err)
	}

	return balance, nil
}

func (k Keeper) ensureCollateralBalance(ctx sdk.Context, market *types.Market, owner sdk.AccAddress, required sdkmath.Int) error {
	if !required.IsPositive() {
		return nil
	}

	switch market.CollateralType {
	case types.CollateralType_COLLATERAL_TYPE_NATIVE:
		bal := k.bankKeeper.GetBalance(ctx, owner, market.CollateralDenom).Amount
		if bal.LT(required) {
			return fmt.Errorf("insufficient funds: required=%s balance=%s", required.String(), bal.String())
		}
		return nil
	case types.CollateralType_COLLATERAL_TYPE_PRC20:
		if err := k.ensurePRC20Balance(ctx, market.CollateralContractAddr, owner, required); err != nil {
			return err
		}
		if err := k.ensurePRC20Allowance(ctx, market.CollateralContractAddr, owner, required); err != nil {
			return err
		}
		return nil
	default:
		return fmt.Errorf("unsupported collateral type")
	}
}
