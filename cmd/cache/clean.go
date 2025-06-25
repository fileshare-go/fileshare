package cache

import (
	"os"

	"github.com/chanmaoganda/fileshare/internal/config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Cleans db file and cache folder by config.yml, if not set then clean the default ones(default.db, $HOME/.fileshare)",
	Run: func(cmd *cobra.Command, args []string) {
		cleanCache()
	},
}

func cleanCache() {
	cfg := config.Cfg()

	if err := os.RemoveAll(cfg.CacheDirectory); err != nil {
		logrus.Errorf("Error removing directory %s, err: %s", cfg.CacheDirectory, err.Error())
	}

	if err := os.Remove(cfg.Database); err != nil {
		if !os.IsNotExist(err) {
			logrus.Errorf("Error removing database %s, err: %s", cfg.Database, err.Error())
			return
		}
	}

	logrus.Info("Cache cleaned successfully")
}
