package pkg

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

// Config represents a list of watches.
type Config struct {
	Watches []Watch `yaml:"watches"`
	Emby    Emby    `yaml:"emby,omitempty"`
}

// Watch represents a watch.
type Watch struct {
	Src      string `yaml:"src"`
	Dst      string `yaml:"dst"`
	Type     string `yaml:"type"`
	URL      string `yaml:"url,omitempty"`
	Username string `yaml:"username,omitempty"`
	Password string `yaml:"password,omitempty"`
}

// Emby represents an Emby server.
type Emby struct {
	URL    string `yaml:"url"`
	APIKey string `yaml:"apiKey"`
}

// ReadConfig reads a YAML file and returns a Config.
func ReadConfig(configPath string) (Config, error) {
	var ret Config
	yamlFile, err := os.ReadFile(configPath)
	if err != nil {
		return ret, fmt.Errorf("Failed to read YAML file: %v", err)
	}
	err = yaml.Unmarshal(yamlFile, &ret)
	if err != nil {
		return ret, fmt.Errorf("Failed to unmarshal YAML: %v", err)
	}
	return ret, nil
}
