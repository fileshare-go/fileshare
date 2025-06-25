package client

import (
	"github.com/chanmaoganda/fileshare/cmd/fileshare"
	"github.com/chanmaoganda/fileshare/internal/config"
	"github.com/chanmaoganda/fileshare/internal/core/setup"
	"github.com/chanmaoganda/fileshare/internal/core/sharelink"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var ShareLinkGenCmd = &cobra.Command{
	Use:     "linkgen <filename> <checksum256> <expire days(optional)>",
	Short:   "Generates sharelink code for friends to easily download",
	Args:    cobra.MinimumNArgs(2),
	PreRunE: setup.Setup,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.Cfg()

		logrus.Debug("Connecting to ", cfg.GrpcAddress)

		conn, err := fileshare.NewClientConn(cfg)
		if err != nil {
			logrus.Fatal(err)
		}

		client := sharelink.NewShareLinkClient(cmd.Context(), conn)

		code := client.GenerateLink(args)
		logrus.Infof("Generated Code is: [%s]", code)

		if err := conn.Close(); err != nil {
			logrus.Error(err)
		}
	},
}
