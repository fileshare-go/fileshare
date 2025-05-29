package config

import (
	"fmt"
	"os"

	"github.com/chanmaoganda/fileshare/internal/pkg/fileutil"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

type Settings struct {
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

func ReadSettings(filename string) (*Settings, error) {
	var settings Settings

	bytes, err := os.ReadFile(filename)
	if err != nil {
		logrus.Warn("cannot open configuration file, use default config")
		if err := settings.SetupEssentials(); err != nil {
			return nil, err
		}
		return &settings, nil
	}

	if err := yaml.Unmarshal(bytes, &settings); err != nil {
		return nil, err
	}

	if err := settings.SetupEssentials(); err != nil {
		return nil, err
	}
	return &settings, nil
}

func (s *Settings) SetupEssentials() error {
	s.FillMissingWithDefault()
	return s.SetupDirectory()
}

func (s *Settings) FillMissingWithDefault() {
	if s.GrpcAddress == "" {
		s.GrpcAddress = ":8080"
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

func (s *Settings) PrintSettings() {
	logrus.Debugf("[Settings] Grpc Address: %s", s.GrpcAddress)
	logrus.Debugf("[Settings] Web Address: %s", s.WebAddress)
	logrus.Debugf("[Settings] Database: %s", s.Database)
	logrus.Debugf("[Settings] ShareCodeLength: %d", s.ShareCodeLength)
	logrus.Debugf("[Settings] CacheDirectory %s", s.CacheDirectory)
	logrus.Debugf("[Settings] DownloadDirectory %s", s.CacheDirectory)
	logrus.Debugf("[Settings] CertPath %s", s.CertsPath)
	logrus.Debugf("[Settings] Valid Days %d", s.ValidDays)
	logrus.Debugf("[Settings] Blocked Ips %v", s.BlockedIps)
}

func (s *Settings) SetupDirectory() error {
	logrus.Debugf("Setting up Directories, %s, %s", s.CacheDirectory, s.DownloadDirectory)
	if fileutil.FileExists(s.CacheDirectory) {
		return nil
	}
	if err := os.Mkdir(s.CacheDirectory, 0755); err != nil {
		return err
	}

	if fileutil.FileExists(s.DownloadDirectory) {
		return nil
	}
	if err := os.Mkdir(s.DownloadDirectory, 0755); err != nil {
		return err
	}

	return nil
}
