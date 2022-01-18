package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Settings stores the app configuration values
type Settings struct {
	Conn struct {
		Schema         string `yaml:"schema"`
		Username       string `yaml:"username"`
		Password       string `yaml:"password"`
		Host           string `yaml:"host"`
		Port           string `yaml:"port"`
		VHost          string `yaml:"vhost"`
		ConnectionName string `yaml:"connection_name"`
	}
}

// Read loads the config from `onfigPath` yml file
func Read(configPath string) (Settings, error) {
	var config Settings

	file, err := os.Open(configPath)
	if err != nil {
		return config, err
	}

	defer file.Close()

	if file == nil {
		return config, fmt.Errorf("Could not parse file %s", configPath)
	}

	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return config, err
	}

	return config, nil
}
