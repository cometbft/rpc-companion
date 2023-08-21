package ingest

import (
	"github.com/cometbft/rpc-companion/service/ingest"
	"github.com/spf13/cobra"
	"log/slog"
	"os"
)

// TODO: make this configurable via a config or parameter
const connString = "postgres://postgres:postgres@0.0.0.0:15432/postgres?sslmode=disable"

// StartCmd start ingest service
var StartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start Ingest Service",
	Long:  `The start command runs an instance of the Ingest Service`,
	Run: func(cmd *cobra.Command, args []string) {
		textHandler := slog.NewTextHandler(os.Stdout, nil)
		logger := slog.New(textHandler)
		service := ingest.NewIngestService(*logger)
		err := service.OnStart()
		if err != nil {
			service.Logger.Error("Failed to start the Ingest Service: %v", err.Error())
		}
		select {}
	},
}

func init() {
}
