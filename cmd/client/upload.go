package client

import (
	"github.com/chanmaoganda/fileshare/internal/core/setup"
	"github.com/chanmaoganda/fileshare/internal/fileshare"
	"github.com/spf13/cobra"
)

var UploadCmd = &cobra.Command{
	Use:     "upload <filepath>",
	Short:   "Uploads the file, requires the filepath as argument",
	Args:    cobra.MinimumNArgs(1),
	PreRunE: setup.Setup,
	Run: func(cmd *cobra.Command, args []string) {
		fileshare.Upload(cmd.Context(), args)
	},
}
