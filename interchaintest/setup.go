package interchaintest

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	globalfee "github.com/strangelove-ventures/globalfee/x/globalfee/types"
	"github.com/strangelove-ventures/interchaintest/v8"
	"github.com/strangelove-ventures/interchaintest/v8/chain/cosmos"
	"github.com/strangelove-ventures/interchaintest/v8/ibc"

	sdktestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
)

var (
	votingPeriod     = "15s"
	maxDepositPeriod = "10s"

	accAddr     = "cosmos1hj5fveer5cjtn4wd6wstzugjfdxzl0xpxvjjvr"
	accMnemonic = "decorate bright ozone fork gallery riot bus exhaust worth way bone indoor calm squirrel merry zero scheme cotton until shop any excess stage laundry"

	CosmosGovModuleAcc = "cosmos10d07y265gmmuvt4z0w9aw880jnsr700j6zn9kn"

	NumberVals         = 1
	NumberFullNodes    = 0
	Denom              = "token"
	GlobalFeeGasPrices = "0.0025"

	DefaultGenesis = []cosmos.GenesisKV{
		cosmos.NewGenesisKV("app_state.gov.params.voting_period", votingPeriod),
		cosmos.NewGenesisKV("app_state.gov.params.max_deposit_period", maxDepositPeriod),
		cosmos.NewGenesisKV("consensus.params.abci.vote_extensions_enable_height", "1"),
		cosmos.NewGenesisKV("app_state.globalfee.params.minimum_gas_prices", sdk.DecCoins{sdk.NewDecCoinFromDec(Denom, sdkmath.LegacyMustNewDecFromStr(GlobalFeeGasPrices))}),
	}

	// `make local-image`
	LocalChainConfig = ibc.ChainConfig{
		Type:    "cosmos",
		Name:    "globalfee",
		ChainID: "globalfee-2",
		Images: []ibc.DockerImage{
			{
				Repository: "globalfee",
				Version:    "local",
				UidGid:     "1025:1025",
			},
		},
		Bin:            "globald",
		Bech32Prefix:   "cosmos",
		Denom:          Denom,
		GasPrices:      GlobalFeeGasPrices + Denom,
		GasAdjustment:  1.3,
		TrustingPeriod: "508h",
		NoHostMount:    false,
		EncodingConfig: AppEncoding(),
		ModifyGenesis:  cosmos.ModifyGenesis(DefaultGenesis),
	}

	DefaultGenesisAmt = sdkmath.NewInt(10_000_000)

	DefaultChainSpec = interchaintest.ChainSpec{
		Name:          LocalChainConfig.Name,
		ChainName:     LocalChainConfig.Name,
		Version:       LocalChainConfig.Images[0].Version,
		ChainConfig:   LocalChainConfig,
		NumValidators: &NumberVals,
		NumFullNodes:  &NumberFullNodes,
	}
)

func AppEncoding() *sdktestutil.TestEncodingConfig {
	enc := cosmos.DefaultEncoding()

	globalfee.RegisterInterfaces(enc.InterfaceRegistry)

	return &enc
}
