package internal

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type User struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type ServerConfig struct {
	Port  int    `yaml:"port"`
	Users []User `yaml:"users"`
}

type Config struct {
	Server ServerConfig `yaml:"server"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		return nil, fmt.Errorf("failed to read config.yaml: %v", err)
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config.yaml: %v", err)
	}

	return &config, nil
}
