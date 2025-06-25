package main

import (
	"github.com/chanmaoganda/fileshare/cmd/cache"
	"github.com/chanmaoganda/fileshare/cmd/server"
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
	// server commands
	rootCmd.AddCommand(server.ServerCmd)

	// cache commands
	rootCmd.AddCommand(cache.CacheCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		logrus.Error(err)
	}
}
