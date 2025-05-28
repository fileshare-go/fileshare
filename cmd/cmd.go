package cmd

import (
	"github.com/chanmaoganda/fileshare/cmd/cache"
	"github.com/chanmaoganda/fileshare/cmd/fileshare/client"
	"github.com/chanmaoganda/fileshare/cmd/fileshare/server"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "fileshare <subcommands>",
	Short: "Fileshare is designed for lightweight file server.",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Usage()
	},
}

func init() {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors: true,
	})

	// server commands
	RootCmd.AddCommand(server.ServerCmd)

	// client commands
	RootCmd.AddCommand(client.UploadCmd)
	RootCmd.AddCommand(client.DownloadCmd)
	RootCmd.AddCommand(client.ShareLinkGenCmd)

	// cache commands
	RootCmd.AddCommand(cache.CacheCmd)
}
