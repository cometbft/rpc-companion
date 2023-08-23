package commands

import (
	"fmt"
	ingest "github.com/cometbft/rpc-companion/commands/ingest"
	rpc "github.com/cometbft/rpc-companion/commands/rpc"
	cfg "github.com/cometbft/rpc-companion/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log/slog"
	"os"
)

var (
	config    = cfg.DefaultConfig()
	log       = slog.New(slog.NewTextHandler(os.Stdout, nil))
	configVar = "RPCDATACOMP"
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(RootCmd.Execute())
}

func init() {

	registerFlagsRootCmd(RootCmd)

	options := cobra.CompletionOptions{
		DisableDefaultCmd:   true,
		DisableNoDescFlag:   true,
		DisableDescriptions: false,
	}
	RootCmd.CompletionOptions = options

	cobra.EnableCommandSorting = true

	RootCmd.AddCommand(ingest.IngestCmd)
	RootCmd.AddCommand(rpc.RpcCmd)
}

func ParseConfig(cmd *cobra.Command) (*cfg.Config, error) {
	conf := cfg.DefaultConfig()
	err := viper.Unmarshal(conf)
	if err != nil {
		return nil, err
	}

	var home string
	if os.Getenv(configVar) != "" {
		home = os.Getenv(configVar)
	} else {
		home, err = cmd.Flags().GetString("config")
		if err != nil {
			return nil, err
		}
	}

	conf.RootDir = home

	//conf.SetRoot(conf.RootDir)
	//cfg.EnsureRoot(conf.RootDir)
	if err := conf.ValidateBasic(); err != nil {
		return nil, fmt.Errorf("error in config file: %v", err)
	}
	return conf, nil
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("toml")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}

func registerFlagsRootCmd(cmd *cobra.Command) {
	cmd.PersistentFlags().String("config", config.RootDir, "")
}

// RootCmd is the root command for CometBFT core.
var RootCmd = &cobra.Command{
	Use:   "rpc-companion",
	Short: "A RPC Data Companion for CometBFT",
	Long:  `A RPC Data Companion is an implementation of a Data Companion for CometBFT based chains`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
		config, err = ParseConfig(cmd)
		if err != nil {
			return err
		}

		return nil
	},
}
