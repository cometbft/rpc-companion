package config

import (
	"fmt"
	"strings"
)

type Config struct {
	// Options for services
	Storage    *StorageConfig    `mapstructure:"storage"`
	GRPCClient *GRPCClientConfig `mapstructure:"grpc_client"`
}

// DefaultConfig returns a default configuration for a node
func DefaultConfig() *Config {
	return &Config{
		Storage:    DefaultStorageConfig(),
		GRPCClient: DefaultGRPCClientConfig(),
	}
}

// ValidateBasic performs basic validation and
// returns an error if any check fails.
func (cfg *Config) ValidateBasic() error {
	if err := cfg.GRPCClient.ValidateBasic(); err != nil {
		return fmt.Errorf("error in [grpc_client] section: %w", err)
	}
	return nil
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
	ListenAddress string `mapstructure:"address"`
}

// DefaultGRPCClientConfig returns a default configuration for the gRPC fetcher layer
func DefaultGRPCClientConfig() *GRPCClientConfig {
	return &GRPCClientConfig{
		ListenAddress: "",
	}
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
	return nil
}
