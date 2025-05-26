package config

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

type Settings struct {
	Address         string `yaml:"address"`
	Database        string `yaml:"database"`
	ShareCodeLength int    `yaml:"share_code_length"`
	CacheDirectory string `yaml:"cache_directory"`
}

func ReadSettings(filename string) (*Settings, error) {
	var settings Settings

	bytes, err := os.ReadFile(filename)
	if err != nil {
		logrus.Warn("cannot open configuration file, use default config")
		settings.FillMissingWithDefault()
		return &settings, nil
	}

	if err := yaml.Unmarshal(bytes, &settings); err != nil {
		return nil, err
	}

	settings.FillMissingWithDefault()
	return &settings, nil
}

func (s *Settings) FillMissingWithDefault() {
	if s.Address == "" {
		s.Address = ":8080"
	}
	if s.Database == "" {
		s.Database = "default.db"
	}
	if s.ShareCodeLength == 0 {
		s.ShareCodeLength = 8
	}
	if s.CacheDirectory == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Printf("Cannot Get Home Directory: %v\n", err)
			return
		}
		s.CacheDirectory = fmt.Sprintf("%s/%s", homeDir, ".fileshare")
	}
}
