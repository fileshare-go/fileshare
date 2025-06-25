package client

import (
	"github.com/chanmaoganda/fileshare/internal/fileshare"
	"github.com/spf13/cobra"
)

var DownloadCmd = &cobra.Command{
	Use:   "download <checksum256 | linkcode>",
	Short: "Download file, either with sharelink code or file checksum256 hash",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fileshare.Download(cmd.Context(), args)
	},
}
