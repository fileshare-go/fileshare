package cache

import (
	"github.com/spf13/cobra"
)

var CacheCmd = &cobra.Command{
	Use:   "cache",
	Short: "Operating cache set in settings.yml, if not set then operate the default cache folder ($HOME/.fileshare)",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}

func init() {
	CacheCmd.AddCommand(cleanCmd)
}
