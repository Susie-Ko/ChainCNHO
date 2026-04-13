package params

import sdk "github.com/cosmos/cosmos-sdk/types"

const (
	Bech32PrefixAccAddr = "cnho"
	Bech32PrefixAccPub  = "cnhopub"

	Bech32PrefixValAddr = "cnhovaloper"
	Bech32PrefixValPub  = "cnhovaloperpub"

	Bech32PrefixConsAddr = "cnhovalcons"
	Bech32PrefixConsPub  = "cnhovalconspub"
)

func SetAddressPrefixes() {
	config := sdk.GetConfig()

	config.SetBech32PrefixForAccount(Bech32PrefixAccAddr, Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(Bech32PrefixValAddr, Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(Bech32PrefixConsAddr, Bech32PrefixConsPub)

	config.Seal()
}
