package setup

import (
	"github.com/chanmaoganda/fileshare/internal/config"
	"github.com/chanmaoganda/fileshare/internal/pkg/logger"
	"github.com/chanmaoganda/fileshare/internal/service"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func Setup(cmd *cobra.Command, args []string) error {
	logger.SetupLogger()
	var err error
	if err = config.ReadConfig(); err != nil {
		logrus.Error(err)
		return err
	}
	service.InitService()
	// if directory cannot be set correctly, following actions will panic
	return nil
}
