package config

import (
	"fmt"
	"io/fs"
	"os"

	"github.com/chanmaoganda/fileshare/internal/pkg/util"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

var config Config

func Cfg() *Config {
	return &config
}

type Config struct {
	GrpcAddress       string   `yaml:"grpc_address"`
	WebAddress        string   `yaml:"web_address"`
	Database          string   `yaml:"database"`
	ShareCodeLength   int      `yaml:"share_code_length"`
	CacheDirectory    string   `yaml:"cache_directory"`
	DownloadDirectory string   `yaml:"download_directory"`
	CertsPath         string   `yaml:"certs_path"`
	ValidDays         int      `yaml:"valid_days"`
	BlockedIps        []string `yaml:"blocked_ips"`
}

const CONFIG_FOLDER = "fileshare"
const CONFIG_FILE = "config.yml"

var configPath string
var configFileMode fs.FileMode = os.FileMode(0644)

// setup config path according to os
func init() {
	setupConfigPath()
}

func ReadConfig() error {
	logrus.Debug("config path is ", configPath)

	bytes, err := os.ReadFile(configPath)
	if err != nil {
		logrus.Warn("cannot open configuration file, use default config")
		config.FillMissingWithDefault()
		// if configFile not found, save with default configurations
		saveConfig()
		return err
	}

	if err := yaml.Unmarshal(bytes, &config); err != nil {
		return err
	}

	config.FillMissingWithDefault()
	return nil
}

func (s *Config) FillMissingWithDefault() {
	if s.GrpcAddress == "" {
		s.GrpcAddress = ":60011"
	}
	if s.WebAddress == "" {
		s.WebAddress = ":8080"
	}
	if s.Database == "" {
		s.Database = "default.db"
	}
	if s.ShareCodeLength == 0 {
		s.ShareCodeLength = 8
	}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Cannot Get Home Directory: %v\n", err)
		return
	}
	if s.CacheDirectory == "" {
		s.CacheDirectory = fmt.Sprintf("%s/%s", homeDir, ".fileshare")
	}
	if s.DownloadDirectory == "" {
		s.DownloadDirectory = fmt.Sprintf("%s/%s", homeDir, "Downloads")
	}
	if s.ValidDays == 0 {
		s.ValidDays = 7
	}
}

func (s *Config) PrintConfig() {
	logrus.Debugf("[Settings] Grpc Address: %s", util.Render(s.GrpcAddress))
	logrus.Debugf("[Settings] Web Address: %s", util.Render(s.WebAddress))
	logrus.Debugf("[Settings] Database: %s", util.Render(s.Database))
	logrus.Debugf("[Settings] ShareCodeLength: %s", util.Render(s.ShareCodeLength))
	logrus.Debugf("[Settings] CacheDirectory %s", util.Render(s.CacheDirectory))
	logrus.Debugf("[Settings] DownloadDirectory %s", util.Render(s.CacheDirectory))
	logrus.Debugf("[Settings] CertPath %s", util.Render(s.CertsPath))
	logrus.Debugf("[Settings] Valid Days %s", util.Render(s.ValidDays))
	logrus.Debugf("[Settings] Blocked Ips %s", util.Render(s.BlockedIps))
}

func saveConfig() {
	bytes, err := yaml.Marshal(&config)
	if err != nil {
		logrus.Error(err)
		return
	}
	if err = os.WriteFile(configPath, bytes, configFileMode); err != nil {
		logrus.Error(err)
	}
}
