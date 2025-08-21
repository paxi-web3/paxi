package app

import (
	"context"

	storetypes "cosmossdk.io/store/types"
	upgradetypes "cosmossdk.io/x/upgrade/types"

	"github.com/cosmos/cosmos-sdk/types/module"
)

// UpgradeName is the name of the upgrade.
// It is used to identify the upgrade in the upgrade handler and
const (
	UpgradeV103 = "v1.0.3"
	UpgradeV104 = "v1.0.4"
	UpgradeV105 = "v1.0.5"
)

func (app *PaxiApp) RegisterUpgradeHandlers() {
	// Common migration handler for both upgrades
	// This will execute all registered module migrations
	migrate := func(ctx context.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		return app.ModuleManager.RunMigrations(ctx, app.Configurator(), fromVM)
	}

	// Register both upgrade handlers so nodes can replay past upgrades
	// NOTE: The upgrade name must match exactly the one on-chain
	app.UpgradeKeeper.SetUpgradeHandler(UpgradeV103, migrate)
	app.UpgradeKeeper.SetUpgradeHandler(UpgradeV104, migrate)
	app.UpgradeKeeper.SetUpgradeHandler(UpgradeV105, migrate)

	// Read upgrade plan information from disk (upgrade-info.json)
	upgradeInfo, err := app.UpgradeKeeper.ReadUpgradeInfoFromDisk()
	if err != nil {
		panic(err)
	}

	// Only apply StoreLoader when we are at the upgrade height and not skipping it
	if !app.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height) {
		switch upgradeInfo.Name {

		case UpgradeV103:
			// Store upgrades for v1.0.3
			// Fill in Added/Deleted stores if this upgrade changed the store layout
			v103Stores := &storetypes.StoreUpgrades{
				Added: []string{"swap"},
				// Deleted: []string{"oldmodule"},
			}
			app.SetStoreLoader(upgradetypes.UpgradeStoreLoader(upgradeInfo.Height, v103Stores))
			app.Logger().Info("Registering upgrade handler for", "upgrade", UpgradeV103)

		case UpgradeV104:
			// Store upgrades for v1.0.4
			// If there are no store changes, this can be nil
			var v104Stores *storetypes.StoreUpgrades = &storetypes.StoreUpgrades{
				// Added: []string{"newmodule"},
				// Deleted: []string{"oldmodule"},
			}
			app.SetStoreLoader(upgradetypes.UpgradeStoreLoader(upgradeInfo.Height, v104Stores))
			app.Logger().Info("Registering upgrade handler for", "upgrade", UpgradeV104)

		case UpgradeV105:
			// Store upgrades for v1.0.5
			// If there are no store changes, this can be nil
			var v105Stores *storetypes.StoreUpgrades = &storetypes.StoreUpgrades{
				// Added: []string{"newmodule"},
				// Deleted: []string{"oldmodule"},
			}
			app.SetStoreLoader(upgradetypes.UpgradeStoreLoader(upgradeInfo.Height, v105Stores))
			app.Logger().Info("Registering upgrade handler for", "upgrade", UpgradeV105)
		}
	}
}
