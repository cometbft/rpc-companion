package commands

import (
	"log/slog"
	"os"

	"github.com/spf13/cobra"
)

var (
	log = slog.New(slog.NewTextHandler(os.Stdout, nil))
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(RootCmd.Execute())
}

func init() {

	addGlobalFlags(RootCmd)

	options := cobra.CompletionOptions{
		DisableDefaultCmd:   true,
		DisableNoDescFlag:   true,
		DisableDescriptions: false,
	}
	RootCmd.CompletionOptions = options

	cobra.EnableCommandSorting = true

	RootCmd.AddCommand(IngestCmd)
}

// RootCmd is the root command for CometBFT core.
var RootCmd = &cobra.Command{
	Use:   "rpc-companion",
	Short: "A RPC Data Companion for CometBFT",
	Long:  `A RPC Data Companion is an implementation of a Data Companion for CometBFT based chains`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
		if err != nil {
			return err
		}

		return nil
	},
}
