package cmd

import (
	"fmt"

	"github.com/chanmaoganda/fileshare/cmd/client"
	"github.com/chanmaoganda/fileshare/cmd/server"
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
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors: true,
	})
	// logrus.SetReportCaller(true)

	RootCmd.AddCommand(client.UploadCmd)
	RootCmd.AddCommand(client.DownloadCmd)
	RootCmd.AddCommand(server.ServerCmd)

	PrintBanner()
}

func PrintBanner() {
	banner := []byte{32, 32, 32, 32, 95, 95, 95, 95, 95, 32, 95, 95, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 95, 95, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 10, 32, 32, 32, 47, 32, 95, 95, 40, 95, 41, 32, 47, 95, 95, 32, 32, 32, 32, 32, 95, 95, 95, 95, 95, 47, 32, 47, 95, 32, 32, 95, 95, 95, 95, 32, 95, 95, 95, 95, 95, 95, 95, 95, 95, 32, 10, 32, 32, 47, 32, 47, 95, 47, 32, 47, 32, 47, 32, 95, 32, 92, 32, 32, 32, 47, 32, 95, 95, 95, 47, 32, 95, 95, 32, 92, 47, 32, 95, 95, 32, 96, 47, 32, 95, 95, 95, 47, 32, 95, 32, 92, 10, 32, 47, 32, 95, 95, 47, 32, 47, 32, 47, 32, 32, 95, 95, 47, 32, 32, 40, 95, 95, 32, 32, 41, 32, 47, 32, 47, 32, 47, 32, 47, 95, 47, 32, 47, 32, 47, 32, 32, 47, 32, 32, 95, 95, 47, 10, 47, 95, 47, 32, 47, 95, 47, 95, 47, 92, 95, 95, 95, 47, 32, 32, 47, 95, 95, 95, 95, 47, 95, 47, 32, 47, 95, 47, 92, 95, 95, 44, 95, 47, 95, 47, 32, 32, 32, 92, 95, 95, 95, 47, 32, 10}
	fmt.Printf("\n%s\n\n\n", string(banner))
}
