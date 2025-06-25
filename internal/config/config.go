package config

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/chanmaoganda/fileshare/internal/pkg/util"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

var cfg Config

func Cfg() *Config {
	return &cfg
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

const CONFIG_FILE = "config.yml"

var configDir string
var configPath string
var configFileMode fs.FileMode = os.FileMode(0644)

func ReadConfig() error {
	var err error
	if err = setupConfigPath(); err != nil {
		return err
	}
	logrus.Debug("config path is ", configPath)

	bytes, err := os.ReadFile(configPath)
	if err != nil {
		logrus.WithError(err).Warn("cannot open configuration file, use default config")
	} else {
		if err := yaml.Unmarshal(bytes, &cfg); err != nil {
			return err
		}
	}

	fillMissingWithDefault(&cfg)
	if err = setupDirectories(); err != nil {
		return err
	}

	print(&cfg)
	saveConfig()
	return nil
}

func fillMissingWithDefault(s *Config) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Print("Cannot Get Home Directory: ", err)
		return
	}

	if s.CacheDirectory == "" {
		s.CacheDirectory = fmt.Sprintf("%s/%s", homeDir, ".fileshare")
	}
	if s.DownloadDirectory == "" {
		s.DownloadDirectory = fmt.Sprintf("%s/%s", homeDir, "Downloads")
	}

	if s.GrpcAddress == "" {
		s.GrpcAddress = ":60011"
	}
	if s.WebAddress == "" {
		s.WebAddress = ":8080"
	}
	if s.Database == "" {
		s.Database = filepath.Join(configDir, "default.db")
	}
	if s.ShareCodeLength == 0 {
		s.ShareCodeLength = 8
	}

	if s.ValidDays == 0 {
		s.ValidDays = 7
	}
}

func print(s *Config) {
	logrus.Debugf("[Settings] Grpc Address: %s", util.Render(s.GrpcAddress))
	logrus.Debugf("[Settings] Web Address: %s", util.Render(s.WebAddress))
	logrus.Debugf("[Settings] Database: %s", util.Render(s.Database))
	logrus.Debugf("[Settings] ShareCodeLength: %s", util.Render(s.ShareCodeLength))
	logrus.Debugf("[Settings] CacheDirectory %s", util.Render(s.CacheDirectory))
	logrus.Debugf("[Settings] DownloadDirectory %s", util.Render(s.DownloadDirectory))
	logrus.Debugf("[Settings] CertPath %s", util.Render(s.CertsPath))
	logrus.Debugf("[Settings] Valid Days %s", util.Render(s.ValidDays))
	logrus.Debugf("[Settings] Blocked Ips %s", util.Render(s.BlockedIps))
}

func saveConfig() {
	bytes, err := yaml.Marshal(&cfg)
	if err != nil {
		logrus.Error(err)
		return
	}
	if err = os.WriteFile(configPath, bytes, configFileMode); err != nil {
		logrus.WithError(err).Error("cannot write back to ", configPath)
	} else {
		logrus.Debug("write back to ", configPath)
	}
}

func setupConfigPath() error {
	if util.FileExists(CONFIG_FILE) {
		configDir = "."
		configPath = CONFIG_FILE
		return nil
	}

	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		return err
	}
	configDir = filepath.Join(userConfigDir, "fileshare")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	configPath = filepath.Join(configDir, CONFIG_FILE)
	return nil
}

func setupDirectories() error {
	var err error
	logrus.Debugf("Setting up Directories, %s, %s", cfg.CacheDirectory, cfg.DownloadDirectory)

	if util.FileExists(cfg.CacheDirectory) {
		return nil
	}
	if err = os.Mkdir(cfg.CacheDirectory, 0755); err != nil {
		return err
	}

	if util.FileExists(cfg.DownloadDirectory) {
		return nil
	}
	if err = os.Mkdir(cfg.DownloadDirectory, 0755); err != nil {
		return err
	}
	return nil
}
