package cache

import (
	"os"

	"github.com/chanmaoganda/fileshare/internal/config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Cleans db file and cache folder by settings.yml, if not set then clean the default ones(default.db, $HOME/.fileshare)",
	Run: func(cmd *cobra.Command, args []string) {
		settings, err := config.ReadSettings("settings.yml")
		if err != nil {
			logrus.Error(err)
			return
		}

		cleanCache(settings)
	},
}

func cleanCache(settings *config.Settings) {
	if err := os.RemoveAll(settings.CacheDirectory); err != nil {
		logrus.Errorf("Error removing directory %s, err: %s", settings.CacheDirectory, err.Error())
	}

	if err := os.Remove(settings.Database); err != nil {
		if !os.IsNotExist(err) {
			logrus.Errorf("Error removing database %s, err: %s", settings.Database, err.Error())
			return
		}
	}

	logrus.Info("Cache cleaned successfully")
}
