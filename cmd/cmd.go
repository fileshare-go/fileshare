package cmd

import (
	"github.com/chanmaoganda/fileshare/cmd/client"
	"github.com/chanmaoganda/fileshare/cmd/server"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use: "fileshare",
	Run: func(cmd *cobra.Command, args []string) {
		logrus.Info("File share is a good tool for handling your own disk")
	},
}

func init() {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors: true,
	})
	// logrus.SetReportCaller(true)

	RootCmd.AddCommand(client.UploadCmd)
	RootCmd.AddCommand(server.ServerCmd)
}
