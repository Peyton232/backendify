package config

import (
	"backendify/pkg/models"

	"github.com/spf13/viper"
)

// LoadConfig reads the configuration file and returns a Config struct
func LoadConfig() (*models.Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config models.Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}
	return &config, nil
}
