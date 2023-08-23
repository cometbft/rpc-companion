package commands

import (
	"github.com/cometbft/rpc-companion/rpc"
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
	RpcCmd.AddCommand(rpcStartCmd)
}

// TODO: make this configurable via a config or parameter
const connString = "postgres://postgres:postgres@0.0.0.0:15432/postgres?sslmode=disable"

// StartCmd start RPC service
var rpcStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start RPC Service",
	Long:  `The start command runs an instance of the RPC Service`,
	Run: func(cmd *cobra.Command, args []string) {

		//Instantiate a new Ingest Service
		rpcSvc := rpc.NewService(connString)

		//Start the service
		rpcSvc.Serve()
	},
}
