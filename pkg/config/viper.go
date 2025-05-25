package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// LoadConfig loads the configuration from the specified path
func LoadConfig(path string) *Config {
	v := viper.New()
	v.SetConfigName("config") // name of config file
	v.SetConfigType("yaml")
	v.AddConfigPath(path) // path to look for the config file in

	// Read the config file
	if err := v.ReadInConfig(); err != nil {
		panic(fmt.Errorf("failed to read config file: %w", err))
	}

	// Unmarshal the config into our Database struct
	var config Config
	if err := v.UnmarshalKey("database", &config); err != nil {
		panic(fmt.Errorf("failed to unmarshal database config: %w", err))
	}

	return &config
}

// LoadConfigFromDefaultPath loads the configuration from the default path
func LoadConfigFromDefaultPath() *Config {
	return LoadConfig(".")
}
