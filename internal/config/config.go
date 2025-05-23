package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Settings struct {
	Address  string `yaml:"address"`
	Database string `yaml:"database"`
	ShareCodeLength int `yaml:"share_code_length"`
}

func ReadSettings(filename string) (*Settings, error) {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var settings Settings

	if err := yaml.Unmarshal(bytes, &settings); err != nil {
		return nil, err
	}

	return &settings, nil
}
