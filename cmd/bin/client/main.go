package main

import (
	"github.com/chanmaoganda/fileshare/cmd/cache"
	"github.com/chanmaoganda/fileshare/cmd/client"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "fileshare <subcommands>",
	Short: "Fileshare is designed for lightweight file server.",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Usage()
	},
}

func init() {
	// client commands
	rootCmd.AddCommand(client.UploadCmd)
	rootCmd.AddCommand(client.DownloadCmd)
	rootCmd.AddCommand(client.ShareLinkGenCmd)

	// cache commands
	rootCmd.AddCommand(cache.CacheCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		logrus.Error(err)
	}
}
