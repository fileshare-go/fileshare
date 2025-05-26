package cmd

import (
	"github.com/chanmaoganda/fileshare/cmd/cache"
	"github.com/chanmaoganda/fileshare/cmd/fileshare/client"
	"github.com/chanmaoganda/fileshare/cmd/fileshare/server"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use: "fileshare",
	Run: func(cmd *cobra.Command, args []string) {
		logrus.Info("File share is a good tool for handling your own cloud disk")
	},
}

func init() {
	// server commands
	RootCmd.AddCommand(server.ServerCmd)

	// client commands
	RootCmd.AddCommand(client.UploadCmd)
	RootCmd.AddCommand(client.DownloadCmd)
	RootCmd.AddCommand(client.ShareLinkGenCmd)

	// cache commands
	RootCmd.AddCommand(cache.CacheCmd)

	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors: true,
	})
}
