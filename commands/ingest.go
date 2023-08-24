package commands

import (
	"log/slog"
	"os"

	"github.com/cometbft/rpc-companion/config"
	"github.com/cometbft/rpc-companion/ingest"
	rpcos "github.com/cometbft/rpc-companion/libs/os"
	"github.com/spf13/cobra"
)

// IngestCmd ingest service commands
var IngestCmd = &cobra.Command{
	Use:   "ingest",
	Short: "Ingest Service commands",
	Long: `The Ingest Service pulls data from a CometBFT full node and store the information in a database. 

There should be just one running instance of the Ingest Service targeting a full node.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	IngestCmd.AddCommand(ingestStartCmd)
}

// StartCmd start ingest service
var ingestStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start Ingest Service",
	Long:  `The start command runs an instance of the Ingest Service`,
	Run: func(cmd *cobra.Command, args []string) {
		textHandler := slog.NewTextHandler(os.Stdout, nil)
		logger := slog.New(textHandler)

		// Load configuration file
		config, err := config.LoadConfig(FlagConfigPath)
		if err != nil {
			logger.Error("error reading configuration file:", err)
			os.Exit(1)
		}

		service, err := ingest.NewIngestService(*logger, config)
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
