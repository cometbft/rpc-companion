package rpc

import (
	"github.com/spf13/cobra"
)

// RpcCmd RPC service commands
var RpcCmd = &cobra.Command{
	Use:   "rpc",
	Short: "RPC Service commands",
	Long:  `The RPC Service expose a RPC endpoint compatible with the CometBFT RPC endpoint.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	RpcCmd.AddCommand(StartCmd)
}
