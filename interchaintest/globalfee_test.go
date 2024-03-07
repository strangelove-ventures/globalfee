package interchaintest

import (
	"context"
	"fmt"
	"testing"

	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/strangelove-ventures/interchaintest/v8"
	"github.com/strangelove-ventures/interchaintest/v8/chain/cosmos"
	"github.com/strangelove-ventures/interchaintest/v8/ibc"
	"github.com/strangelove-ventures/interchaintest/v8/testutil"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

func TestGlobalFee(t *testing.T) {
	cf := interchaintest.NewBuiltinChainFactory(zaptest.NewLogger(t), []*interchaintest.ChainSpec{
		&DefaultChainSpec,
	})

	chains, err := cf.Chains(t.Name())
	require.NoError(t, err)

	chain := chains[0].(*cosmos.CosmosChain)

	ic := interchaintest.NewInterchain().
		AddChain(chain)

	ctx := context.Background()
	client, network := interchaintest.DockerSetup(t)

	require.NoError(t, ic.Build(ctx, nil, interchaintest.InterchainBuildOptions{
		TestName:         t.Name(),
		Client:           client,
		NetworkID:        network,
		SkipPathCreation: true,
	}))
	t.Cleanup(func() {
		_ = ic.Close()
	})

	// faucet funds to the user
	users := interchaintest.GetAndFundTestUsers(t, ctx, "default", DefaultGenesisAmt, chain, chain)
	user := users[0]
	toUser := users[1]

	// balance check
	balance, err := chain.GetBalance(ctx, user.FormattedAddress(), Denom)
	require.NoError(t, err)
	require.True(t, balance.Equal(DefaultGenesisAmt), "user balance should be equal to genesis funds")

	// std := bankSendWithFees(t, ctx, juno, sender, receiver, "1"+nativeDenom, "0"+nativeDenom, 200000)

	token := fmt.Sprintf("1%s", Denom)

	// this should fail for not enough funds
	// res, err := chain.GetNode().ExecTx(ctx, user.KeyName(), "bank", "send", user.KeyName(), toUser.FormattedAddress(), token, "--fees=0token")

	// none
	res := bankSendWithFees(t, ctx, chain, user, toUser.FormattedAddress(), token, "0token", 200000)
	require.Contains(t, res, "no fees were specified", res)

	// not enough
	res = bankSendWithFees(t, ctx, chain, user, toUser.FormattedAddress(), token, "1token", 200000)
	require.Contains(t, res, "insufficient fees", res)

	// wrong fee
	res = bankSendWithFees(t, ctx, chain, user, toUser.FormattedAddress(), token, "1NOTTOKEN", 200000)
	require.Contains(t, res, "this fee denom is not accepted", res)

	// success
	res = bankSendWithFees(t, ctx, chain, user, toUser.FormattedAddress(), token, "500token", 200000)
	require.Contains(t, res, "code: 0", res)
}

// We ignore some of the safeguards interchaintest puts in place (such as gas prices and adjustment, since we are testing fees)
func bankSendWithFees(t *testing.T, ctx context.Context, chain *cosmos.CosmosChain, from ibc.Wallet, toAddr, coins, feeCoin string, gasAmt int64) string {
	cmd := []string{chain.Config().Bin, "tx", "bank", "send", from.KeyName(), toAddr, coins,
		"--node", chain.GetRPCAddress(),
		"--home", chain.HomeDir(),
		"--chain-id", chain.Config().ChainID,
		"--gas", fmt.Sprintf("%d", gasAmt),
		"--fees", feeCoin,
		"--keyring-dir", chain.HomeDir(),
		"--keyring-backend", keyring.BackendTest,
		"-y",
	}
	stdout, _, err := chain.Exec(ctx, cmd, nil)
	require.NoError(t, err)

	if err := testutil.WaitForBlocks(ctx, 2, chain); err != nil {
		t.Fatal(err)
	}

	return string(stdout)
}
