package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config holds the CLI configuration
type Config struct {
	AWSRegion     string `mapstructure:"aws_region"`
	AWSProfile    string `mapstructure:"aws_profile"`
	AWSAccountID  string `mapstructure:"aws_account_id"`
	DefaultRole   string `mapstructure:"default_role"`
	CatalogURL    string `mapstructure:"catalog_url"`
	S3DataLake    string `mapstructure:"s3_data_lake"`
	SchemaRegistry string `mapstructure:"schema_registry_url"`
}

// LoadConfig loads the configuration from the config file and environment variables
func LoadConfig() (*Config, error) {
	// Find home directory
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("could not find home directory: %w", err)
	}

	// Create config directory if it doesn't exist
	configDir := filepath.Join(home, ".fedramp-data-mesh")
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		if err := os.MkdirAll(configDir, 0755); err != nil {
			return nil, fmt.Errorf("could not create config directory: %w", err)
		}
	}

	configName := "config"
	configType := "yaml"
	configPath := configDir

	// Configure viper
	viper.SetConfigName(configName)
	viper.SetConfigType(configType)
	viper.AddConfigPath(configPath)

	// Set default values
	viper.SetDefault("aws_region", "us-east-1")
	viper.SetDefault("aws_profile", "")
	viper.SetDefault("aws_account_id", "")
	viper.SetDefault("default_role", "")
	viper.SetDefault("catalog_url", "")
	viper.SetDefault("s3_data_lake", "")
	viper.SetDefault("schema_registry_url", "")

	// Read from environment variables
	viper.AutomaticEnv()
	viper.SetEnvPrefix("DATAMESH")

	// Create default config file if it doesn't exist
	configFile := filepath.Join(configPath, configName+"."+configType)
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		if err := viper.SafeWriteConfigAs(configFile); err != nil {
			return nil, fmt.Errorf("could not write default config file: %w", err)
		}
	}

	// Read the config file
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("could not read config file: %w", err)
	}

	// Unmarshal into config struct
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("could not unmarshal config: %w", err)
	}

	return &config, nil
}

// SaveConfig saves the configuration to the config file
func SaveConfig(config *Config) error {
	// Set all config values
	viper.Set("aws_region", config.AWSRegion)
	viper.Set("aws_profile", config.AWSProfile)
	viper.Set("aws_account_id", config.AWSAccountID)
	viper.Set("default_role", config.DefaultRole)
	viper.Set("catalog_url", config.CatalogURL)
	viper.Set("s3_data_lake", config.S3DataLake)
	viper.Set("schema_registry_url", config.SchemaRegistry)

	// Write the config file
	if err := viper.WriteConfig(); err != nil {
		return fmt.Errorf("could not write config file: %w", err)
	}

	return nil
}
