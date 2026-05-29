package v3

import (
	storetypes "github.com/cosmos/cosmos-sdk/store/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/types/module"

	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	tokenfactorykeeper "cnho/x/tokenfactory/keeper"
	tokenfactorytypes "cnho/x/tokenfactory/types"
)

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	tokenFactoryKeeper tokenfactorykeeper.Keeper,
) upgradetypes.UpgradeHandler {

	return func(
		ctx sdk.Context,
		plan upgradetypes.Plan,
		vm module.VersionMap,
	) (module.VersionMap, error) {

		ctx.Logger().Info("🚀 RUNNING UPGRADE " + UpgradeName)

		// Run migrations
		newVM, err := mm.RunMigrations(ctx, configurator, vm)
		if err != nil {
			return vm, err
		}

		// Optional param reset
		tokenFactoryKeeper.SetParams(
			ctx,
			tokenfactorytypes.DefaultParams(),
		)

		ctx.Logger().Info("✅ UPGRADE " + UpgradeName + " SUCCESS")

		return newVM, nil
	}
}

func StoreUpgrades() *storetypes.StoreUpgrades {
	return &storetypes.StoreUpgrades{
		Added: []string{},
	}
}
