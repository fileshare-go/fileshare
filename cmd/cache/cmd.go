package cache

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var CacheCmd = &cobra.Command{
	Use: "cache",
	Run: func(cmd *cobra.Command, args []string) {
		logrus.Info("Cache operations for fileshare")
	},
}

func init() {
	CacheCmd.AddCommand(cleanCmd)
}
