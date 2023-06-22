package rpc

import (
	"github.com/cometbft/rpc-companion/rpc"
	"github.com/spf13/cobra"
)

// TODO: make this configurable via a config or parameter
const connString = "postgres://postgres:postgres@0.0.0.0:15432/postgres?sslmode=disable"

// StartCmd start RPC service
var StartCmd = &cobra.Command{
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

func init() {
}
