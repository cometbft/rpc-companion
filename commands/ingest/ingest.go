package ingest

import (
	"github.com/spf13/cobra"
)

// IngestCmd ingest service commands
var IngestCmd = &cobra.Command{
	Use:   "ingest",
	Short: "Ingest Service commands",
	Long:  `The Ingest Service pulls data from a CometBFT full node and store the information in a database`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
}
