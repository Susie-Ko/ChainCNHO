package keeper

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"cnho/x/tokenfactory/types"
)

// ConvertToBaseToken converts a fee amount in a whitelisted fee token to the base fee token amount
func (k Keeper) CreateDenom(ctx sdk.Context, creatorAddr string, subdenom string) (newTokenDenom string, err error) {
	// Temporary check until IBC bug is sorted out

	baseSubdenom := strings.ToLower(subdenom)
	symbolDenom := strings.ToUpper(subdenom)

	if k.bankKeeper.HasSupply(ctx, baseSubdenom) {
		return "", fmt.Errorf("Temporary error until IBC bug is sorted out, " +
			"can't create subdenoms that are the same as a native denom.")
	}

	// Send creation fee to community pool
	creationFee := k.GetParams(ctx).DenomCreationFee
	accAddr, err := sdk.AccAddressFromBech32(creatorAddr)
	if err != nil {
		return "", err
	}
	if len(creationFee) > 0 {
		if err := k.distrKeeper.FundCommunityPool(ctx, creationFee, accAddr); err != nil {
			return "", err
		}
	}

	denom, err := types.GetTokenDenom(creatorAddr, baseSubdenom)
	if err != nil {
		return "", err
	}

	_, found := k.bankKeeper.GetDenomMetaData(ctx, denom)
	if found {
		return "", types.ErrDenomExists
	}

	denomMetaData := banktypes.Metadata{
		Description: fmt.Sprintf("%s Stable is a tokenized Real World Asset (RWA) backed by underlying off-chain assets. Assets labeled 'Stable' are pegged to fiat currencies, 'Stock' represents tokenized equity instruments, 'Bond' represents tokenized fixed-income securities, while other asset categories may represent commodities, real estate, ETFs, treasury products, and other real-world financial assets.", symbolDenom),

		Base:    denom,
		Display: baseSubdenom,

		Name:   fmt.Sprintf("%s Stable", symbolDenom),
		Symbol: symbolDenom,

		DenomUnits: []*banktypes.DenomUnit{
			{
				Denom:    denom,
				Exponent: 0,
			},
			{
				Denom:    baseSubdenom,
				Exponent: 6,
			},
		},
	}

	err = denomMetaData.Validate()
	if err != nil {
		return "", err
	}

	k.bankKeeper.SetDenomMetaData(ctx, denomMetaData)

	authorityMetadata := types.DenomAuthorityMetadata{
		Admin: creatorAddr,
	}
	err = k.setAuthorityMetadata(ctx, denom, authorityMetadata)
	if err != nil {
		return "", err
	}

	k.addDenomFromCreator(ctx, creatorAddr, denom)
	return denom, nil
}
