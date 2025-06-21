package cache

import (
	"github.com/spf13/cobra"
)

var CacheCmd = &cobra.Command{
	Use:   "cache",
	Short: "Operating cache set in config.yml, if not set then operate the default cache folder ($HOME/.fileshare)",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Usage()
	},
}

func init() {
	CacheCmd.AddCommand(cleanCmd)
}
