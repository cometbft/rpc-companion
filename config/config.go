package config

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	// Top level options use an anonymous struct
	BaseConfig `mapstructure:",squash"`

	// Options for services
	Storage    *StorageConfig    `mapstructure:"storage"`
	GRPCClient *GRPCClientConfig `mapstructure:"grpc_client"`
}

// ValidateBasic performs basic validation and
// returns an error if any check fails.
func (cfg *Config) ValidateBasic() error {
	if err := cfg.GRPCClient.ValidateBasic(); err != nil {
		return fmt.Errorf("error in [grpc_client] section: %w", err)
	}
	return nil
}

type BaseConfig struct { //nolint: maligned
	// The root directory for all data.
	// This should be set in viper, so it can unmarshal into this struct
	RootDir string `mapstructure:"home"`
}

//-----------------------------------------------------------------------------
// StorageConfig

// StorageConfig defines the configuration options for the storage layer
type StorageConfig struct { //nolint: maligned
	// Connection credentials
	Connection string `mapstructure:"connection"`
}

// DefaultStorageConfig returns a default configuration for the Storage layer
func DefaultStorageConfig() *StorageConfig {
	return &StorageConfig{
		Connection: "",
	}
}

//-----------------------------------------------------------------------------
// StorageConfig

// GRPCClientConfig defines the configuration options for the gRPC fetcher layer
type GRPCClientConfig struct { //nolint: maligned
	// GRPC service address
	ListenAddress           string `mapstructure:"address"`
	ListenAddressPrivileged string `mapstructure:"privileged_address"`
}

// ValidateBasic performs basic validation for the
// [grpc_client] config section
func (cfg *GRPCClientConfig) ValidateBasic() error {
	if len(cfg.ListenAddress) > 0 {
		addrParts := strings.SplitN(cfg.ListenAddress, "://", 2)
		if len(addrParts) != 2 {
			return fmt.Errorf(
				"invalid listening address %s (use fully formed addresses, including the tcp:// or unix:// prefix)",
				cfg.ListenAddress,
			)
		}
	} else {
		return fmt.Errorf("invalid gRPC fetcher listening address, cannot be blank, please ensure a value is set in the config")
	}

	if len(cfg.ListenAddressPrivileged) > 0 {
		addrParts := strings.SplitN(cfg.ListenAddressPrivileged, "://", 2)
		if len(addrParts) != 2 {
			return fmt.Errorf(
				"invalid priviledged listening address %s (use fully formed addresses, including the tcp:// or unix:// prefix)",
				cfg.ListenAddressPrivileged,
			)
		}
	} else {
		return fmt.Errorf("invalid priviledged listening address, cannot be blank, please ensure a value is set in the config")
	}

	return nil
}

func LoadConfig(configPath string) (Config, error) {
	config := Config{}
	if configPath != "" {
		p := filepath.Join(configPath)
		filename := filepath.Base(p)
		ext := filepath.Ext(filename)
		configName := strings.TrimSuffix(filename, ext)
		path := filepath.Dir(p)
		viper.SetConfigName(configName)
		viper.SetConfigType(strings.Replace(ext, ".", "", 1))
		viper.AddConfigPath(path)
	} else {
		viper.SetConfigName("config") // name of config file (without extension)
		viper.SetConfigType("toml")
		viper.AddConfigPath("$HOME/.rpc-companion") // call multiple times to add many search paths
		viper.AddConfigPath(".")
	}

	err := viper.ReadInConfig() // Find and read the config file

	if err != nil { // Handle errors reading the config file
		return config, fmt.Errorf("error reading configuration file")
	} else {
		err := viper.Unmarshal(&config)
		if err != nil {
			return config, fmt.Errorf("cannot unmarshall configuration file")
		}
		if err := config.GRPCClient.ValidateBasic(); err != nil {
			return config, fmt.Errorf("error validating gRPC client configuration: %v", err)
		}
		return config, nil
	}
}
