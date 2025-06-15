package customstaking

import (
	"context"

	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/module"
	staking "github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/staking/exported"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	"github.com/cosmos/cosmos-sdk/x/staking/types"
)

type AppModule struct {
	staking.AppModule
	customKeeper *CustomStakingKeeper
}

func NewAppModule(
	cdc codec.Codec,
	customKeeper *CustomStakingKeeper,
	ak types.AccountKeeper,
	bk types.BankKeeper,
	ls exported.Subspace,
) AppModule {
	original := staking.NewAppModule(cdc, customKeeper.Keeper, ak, bk, ls)
	return AppModule{AppModule: original, customKeeper: customKeeper}
}

// EndBlock returns the end blocker for the staking module. It returns no validator
// updates.
func (am AppModule) EndBlock(ctx context.Context) ([]abci.ValidatorUpdate, error) {
	// override the EndBlocker with our custom function
	return am.customKeeper.EndBlocker(ctx)
}

// RegisterServices registers the staking module's message server and query server.
func (am AppModule) RegisterServices(cfg module.Configurator) {
	types.RegisterMsgServer(
		cfg.MsgServer(),
		NewMsgServer(stakingkeeper.NewMsgServerImpl(am.customKeeper.Keeper)),
	)

	types.RegisterQueryServer(cfg.QueryServer(), stakingkeeper.NewQuerier(am.customKeeper.Keeper))
}
