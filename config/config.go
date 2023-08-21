package config

type Config struct {
	// Options for services
	Storage *StorageConfig `mapstructure:"storage"`
}

// DefaultConfig returns a default configuration for a node
func DefaultConfig() *Config {
	return &Config{}
}

// DefaultStorageConfig returns a default configuration for the Storage layer
func DefaultStorageConfig() *StorageConfig {
	return &StorageConfig{
		Connection: "", //TODO: Add postgres connection
	}
}

//-----------------------------------------------------------------------------
// StorageConfig

// StorageConfig defines the configuration options for the storage layer
type StorageConfig struct { //nolint: maligned
	// Connection credentials
	Connection string `mapstructure:"connection"`
}
