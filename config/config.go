package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	ProjectName string
	Modules     map[string]ModuleSpec
}

type ModuleSpec struct {
	Type  string
	Group string
}

func Load(file string) (Config, error) {
	if _, err := os.Stat(file); err != nil {
		// config file doesn't exist, so just return an empty config
		return Config{}, nil
	}

	viper.SetConfigFile(file)
	viper.SetConfigType("yaml")

	err := viper.ReadInConfig()
	if err != nil {
		return Config{}, fmt.Errorf("could not load config file: %w", err)
	}

	var config Config
	err = viper.Unmarshal(&config)
	if err != nil {
		return Config{}, fmt.Errorf("could not unmarshal config: %w", err)
	}

	return config, nil
}
