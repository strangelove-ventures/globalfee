package main

import (
	"fmt"
	"os"

	"github.com/strangelove-ventures/globalfee/app"
	"github.com/strangelove-ventures/globalfee/cmd/globald/cmd"

	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
)

func main() {
	rootCmd := cmd.NewRootCmd()

	if err := svrcmd.Execute(rootCmd, "globalfee", app.DefaultNodeHome); err != nil {
		fmt.Fprintln(rootCmd.OutOrStderr(), err)
		os.Exit(1)
	}
}
