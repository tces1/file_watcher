package pkg

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

// Example YAML file:
// watches:
//   - src: /root/test
//     dst: /root/tmptest
//     type: local

// WatchList represents a list of watches.
type WatchList struct {
	Watches []Watch `yaml:"watches"`
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

// ReadConfig reads a YAML file and returns a WatchList.
func ReadConfig(configPath string) (WatchList, error) {
	var ret WatchList
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
