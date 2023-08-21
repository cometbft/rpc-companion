package commands

import (
	ingest "github.com/cometbft/rpc-companion/commands/ingest"
	rpc "github.com/cometbft/rpc-companion/commands/rpc"
	"github.com/spf13/cobra"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "http-companion",
	Short: "RPC Companion - CometBFT",
	Long:  `RPC Companion is an implementation of a Data Companion for CometBFT based chains`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(RootCmd.Execute())
}

func init() {

	options := cobra.CompletionOptions{
		DisableDefaultCmd:   true,
		DisableNoDescFlag:   true,
		DisableDescriptions: false,
	}
	RootCmd.CompletionOptions = options

	cobra.EnableCommandSorting = true

	RootCmd.AddCommand(ingest.IngestCmd)
	RootCmd.AddCommand(rpc.RpcCmd)
}

func initConfig() {

}
