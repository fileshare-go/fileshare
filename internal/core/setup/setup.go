package setup

import (
	"os"

	"github.com/chanmaoganda/fileshare/internal/config"
	"github.com/chanmaoganda/fileshare/internal/pkg/util"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func Setup(cmd *cobra.Command, args []string) error {
	var err error
	if err = config.ReadConfig(); err != nil {
		logrus.Error(err)
		return err
	}

	// if directory cannot be set correctly, following actions will panic
	return setupDirectory()
}

func setupDirectory() error {
	c := config.Cfg()
	logrus.Debugf("Setting up Directories, %s, %s", c.CacheDirectory, c.DownloadDirectory)
	if util.FileExists(c.CacheDirectory) {
		return nil
	}
	if err := os.Mkdir(c.CacheDirectory, 0755); err != nil {
		return err
	}

	if util.FileExists(c.DownloadDirectory) {
		return nil
	}
	if err := os.Mkdir(c.DownloadDirectory, 0755); err != nil {
		return err
	}

	return nil
}
