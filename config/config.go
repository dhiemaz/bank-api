package config

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"sync"
)

type Config struct {
	SymmetricKey string `mapstructure:"symmetric_key"`
	Database     struct {
		MigrationPath string `mapstructure:"migration_path"`
		URL           string `mapstructure:"url"`
		Driver        string `mapstructure:"driver"`
		User          string `mapstructure:"user"`
		Password      string `mapstructure:"password"`
		Name          string `mapstructure:"name"`
	} `mapstructure:"database"`
	Env string `mapstructure:"env"`
}

var (
	cfg  *Config
	once sync.Once
)

func LoadConfig(path string) (*Config, error) {
	viper.SetConfigName("config") // Name of config file (without extension)
	viper.SetConfigType("yaml")   // Config file type
	viper.AddConfigPath(path)     // Path to look for the config file

	// Read config file
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	// Unmarshal into struct
	var newCfg Config
	if err := viper.Unmarshal(&newCfg); err != nil {
		return nil, fmt.Errorf("error unmarshalling config: %w", err)
	}

	return &newCfg, nil
}

// GetConfig ensures the config is loaded only once
func GetConfig() *Config {
	once.Do(func() {
		var err error
		cfg, err = LoadConfig(".") // Load config from current directory
		if err != nil {
			log.Fatalf("Failed to load config: %v", err)
		}
	})
	return cfg
}
