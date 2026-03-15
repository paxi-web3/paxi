package keeper

import (
	"context"
	"fmt"
	"strconv"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/paxi-web3/paxi/x/custommint/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(k Keeper) types.MsgServer {
	return &msgServer{Keeper: k}
}

func (k msgServer) UpdateParams(goCtx context.Context, msg *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if msg.Authority != k.authority {
		return nil, sdkerrors.ErrUnauthorized
	}

	// Parse and validate BurnThreshold
	burnThreshold, ok := sdkmath.NewIntFromString(msg.Params.BurnThreshold)
	if !ok {
		return nil, fmt.Errorf("invalid BurnThreshold: %q", msg.Params.BurnThreshold)
	}
	if burnThreshold.IsNegative() {
		return nil, fmt.Errorf("BurnThreshold must be non-negative")
	}

	parsedParams := types.Params{
		BurnThreshold:       burnThreshold,
		BurnRatio:           sdkmath.LegacyMustNewDecFromStr(msg.Params.BurnRatio),
		BlocksPerYear:       msg.Params.BlocksPerYear,
		FirstYearInflation:  sdkmath.LegacyMustNewDecFromStr(msg.Params.FirstYearInflation),
		SecondYearInflation: sdkmath.LegacyMustNewDecFromStr(msg.Params.SecondYearInflation),
		OtherYearInflation:  sdkmath.LegacyMustNewDecFromStr(msg.Params.OtherYearInflation),
		UUSDTAuthority:      msg.Params.UusdtAuthority,
	}
	if parsedParams.UUSDTAuthority == "" {
		parsedParams.UUSDTAuthority = k.GetParams(ctx).UUSDTAuthority
	}

	if err := parsedParams.Validate(); err != nil {
		return nil, err
	}

	k.SetParams(ctx, parsedParams)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"custommint_params_updated",
			sdk.NewAttribute("burn_threshold", parsedParams.BurnThreshold.String()),
			sdk.NewAttribute("burn_ratio", parsedParams.BurnRatio.String()),
			sdk.NewAttribute("blocks_per_year", strconv.FormatInt(parsedParams.BlocksPerYear, 10)),
			sdk.NewAttribute("first_year_inflation", parsedParams.FirstYearInflation.String()),
			sdk.NewAttribute("second_year_inflation", parsedParams.SecondYearInflation.String()),
			sdk.NewAttribute("other_year_inflation", parsedParams.OtherYearInflation.String()),
			sdk.NewAttribute("uusdt_authority", parsedParams.UUSDTAuthority),
		),
	)

	return &types.MsgUpdateParamsResponse{}, nil
}

func (k msgServer) MintUUSDT(goCtx context.Context, msg *types.MsgMintUUSDT) (*types.MsgMintUUSDTResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if msg == nil {
		return nil, sdkerrors.ErrInvalidRequest
	}

	params := k.GetParams(ctx)
	if msg.Authority != params.UUSDTAuthority {
		return nil, sdkerrors.ErrUnauthorized
	}

	toAddr, err := sdk.AccAddressFromBech32(msg.ToAddress)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid to_address: %v", err)
	}

	amount, ok := sdkmath.NewIntFromString(msg.Amount)
	if !ok {
		return nil, sdkerrors.ErrInvalidRequest.Wrap("invalid amount")
	}
	if !amount.IsPositive() {
		return nil, sdkerrors.ErrInvalidRequest.Wrap("amount must be positive")
	}

	coin := sdk.NewCoin(types.UUSDTDenom, amount)
	if err := k.bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(coin)); err != nil {
		return nil, err
	}

	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, toAddr, sdk.NewCoins(coin)); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"custommint_uusdt_minted",
			sdk.NewAttribute("authority", msg.Authority),
			sdk.NewAttribute("to_address", msg.ToAddress),
			sdk.NewAttribute("amount", amount.String()),
		),
	)

	return &types.MsgMintUUSDTResponse{}, nil
}

func (k msgServer) BurnUUSDT(goCtx context.Context, msg *types.MsgBurnUUSDT) (*types.MsgBurnUUSDTResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if msg == nil {
		return nil, sdkerrors.ErrInvalidRequest
	}

	params := k.GetParams(ctx)
	if msg.Authority != params.UUSDTAuthority {
		return nil, sdkerrors.ErrUnauthorized
	}

	fromAddr, err := sdk.AccAddressFromBech32(msg.FromAddress)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid from_address: %v", err)
	}

	amount, ok := sdkmath.NewIntFromString(msg.Amount)
	if !ok {
		return nil, sdkerrors.ErrInvalidRequest.Wrap("invalid amount")
	}
	if !amount.IsPositive() {
		return nil, sdkerrors.ErrInvalidRequest.Wrap("amount must be positive")
	}

	coin := sdk.NewCoin(types.UUSDTDenom, amount)
	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, fromAddr, types.ModuleName, sdk.NewCoins(coin)); err != nil {
		return nil, err
	}

	if err := k.bankKeeper.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(coin)); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"custommint_uusdt_burned",
			sdk.NewAttribute("authority", msg.Authority),
			sdk.NewAttribute("from_address", msg.FromAddress),
			sdk.NewAttribute("amount", amount.String()),
		),
	)

	return &types.MsgBurnUUSDTResponse{}, nil
}
