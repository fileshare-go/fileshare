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
	Long: `Fileshare is designed for lightweight file server.
Example Usages:
- fileshare server
- fileshare upload llvm-2.2.tar.gz
- fileshare download fzHghSyr
- fileshare download 788d871aec139e0c61d49533d0252b21c4cd030e91405491ee8cb9b2d0311072
- fileshare linkgen llvm-2.2.tar.gz 788d871aec139e0c61d49533d0252b21c4cd030e91405491ee8cb9b2d0311072
- fileshare cache clean
	`,
	Run: func(cmd *cobra.Command, args []string) {
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
