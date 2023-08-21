package ingest

import (
	"fmt"
	"github.com/cometbft/rpc-companion/ingest"
	"github.com/spf13/cobra"
	"log"
	"os"
	"strconv"
)

// TODO: make this configurable via a config or parameter
const connString = "postgres://postgres:postgres@0.0.0.0:15432/postgres?sslmode=disable"

// StartCmd start ingest service
var StartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start Ingest Service",
	Long:  `The start command runs an instance of the Ingest Service`,
	Run: func(cmd *cobra.Command, args []string) {

		//Instantiate a new Ingest Service
		ingestSvc := ingest.NewService(connString)

		//Insert some blocks
		numberHeights := int64(10)
		initialHeightParameter := os.Getenv("COMPANION_INITIAL_HEIGHT")
		initialHeight, err := strconv.ParseInt(initialHeightParameter, 10, 64)
		if err != nil {
			fmt.Printf("Invalid initial height %s: %s\n", initialHeightParameter, err)
		}

		for height := initialHeight; height <= initialHeight+numberHeights; height++ {

			header, err := ingestSvc.Client.Header(int64(height))
			if err != nil {
				log.Fatalf("Error fetching block at height %d: %s\n", height, err)
			}

			err = ingestSvc.Storage.InsertHeader(int64(height), *header)
			if err != nil {
				fmt.Printf("error inserting header at height %d: %s\n", height, err)
			}
		}
	},
}

func init() {
}
