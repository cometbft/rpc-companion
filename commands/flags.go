package commands

import "github.com/spf13/cobra"

var (
	FlagConfigPath string
)

// addGlobalFlags defines flags to be used regardless of the command used
func addGlobalFlags(cmd *cobra.Command) {
	RootCmd.PersistentFlags().StringVarP(&FlagConfigPath, "config", "f", "", "configuration file")
}
