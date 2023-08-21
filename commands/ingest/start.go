package ingest

import (
	rpcos "github.com/cometbft/rpc-companion/libs/os"
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

		// Stop upon receiving SIGTERM or CTRL-C.
		rpcos.TrapSignal(*logger, func() {
			// Cleanup
			if err := service.Stop(); err != nil {
				logger.Error("Error while stopping server", "err", err)
			}
		})

		select {}
	},
}

func init() {
}
