package client

import (
	"github.com/chanmaoganda/fileshare/internal/core/setup"
	"github.com/chanmaoganda/fileshare/internal/fileshare"
	"github.com/spf13/cobra"
)

var ShareLinkGenCmd = &cobra.Command{
	Use:     "linkgen <filename> <checksum256> <expire days(optional)>",
	Short:   "Generates sharelink code for friends to easily download",
	Args:    cobra.MinimumNArgs(2),
	PreRunE: setup.Setup,
	Run: func(cmd *cobra.Command, args []string) {
		fileshare.LinkGen(cmd.Context(), args)
	},
}
