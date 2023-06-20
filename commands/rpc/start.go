package rpc

import (
	"github.com/cometbft/rpc-companion/rpc"
	"github.com/spf13/cobra"
	"log"
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
		err := rpcSvc.Serve()
		if err != nil {
			log.Fatalln("There's an error starting the RPC service:", err)
		} else {
			log.Println("Started RPC service...")
		}

	},
}

func init() {
}
