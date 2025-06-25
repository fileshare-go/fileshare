package server

import (
	"github.com/chanmaoganda/fileshare/internal/core/setup"
	"github.com/chanmaoganda/fileshare/internal/fileshare"
	"github.com/spf13/cobra"
)

var ServerCmd = &cobra.Command{
	Use:     "server",
	Short:   "Starts fileshare server",
	PreRunE: setup.Setup,
	Run: func(cmd *cobra.Command, args []string) {
		fileshare.Server()
	},
}
