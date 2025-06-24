package config

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/sirupsen/logrus"
)

var configDir string = "."

func setupConfigPath() {
	goOs := runtime.GOOS
	switch goOs {
	case "windows":
		setupWindows()
	case "linux":
		setupLinux()
	default:
		logrus.Warn("unsupported platform ", goOs)
	}

	if err := os.MkdirAll(configDir, 0755); err != nil {
		logrus.WithError(err).Error("Failed to create config directory ", configDir)
	}

	configPath = filepath.Join(configDir, CONFIG_FILE)
}

func setupWindows() {
	appData := os.Getenv("APPDATA")
	if appData == "" {
		logrus.Warn("APPDATA environment variable not set, falling back to current directory for configPath on Windows.")
		return
	}
	configDir = filepath.Join(appData, CONFIG_FOLDER)
}

func setupLinux() {
	home, err := os.UserHomeDir()
	if err != nil {
		logrus.WithError(err).Warn("Failed to get user home directory, falling back to current directory for configPath on Linux.")
		return
	}

	configDir = filepath.Join(home, CONFIG_FOLDER)
}
