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
		service, err := ingest.NewIngestService(*logger)
		if err != nil {
			logger.Error("failed to instantiate a new ingest service", "error", err)
		}
		err = service.Start()
		if err != nil {
			logger.Error("failed to start the ingest service", "error", err)
			if service.IsRunning() {
				service.Stop()
			}
			logger.Info("ingest service start aborted")
			os.Exit(1)
		}

		// Stop upon receiving SIGTERM or CTRL-C.
		rpcos.TrapSignal(*logger, func() {
			// Cleanup
			if err := service.Stop(); err != nil {
				logger.Error("error while stopping server", "err", err)
			}
		})

		select {}
	},
}

func init() {
}
