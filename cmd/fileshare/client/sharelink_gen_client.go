package client

import (
	"github.com/chanmaoganda/fileshare/cmd/fileshare"
	"github.com/chanmaoganda/fileshare/internal/config"
	"github.com/chanmaoganda/fileshare/internal/fileshare/sharelink"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var ShareLinkGenCmd = &cobra.Command{
	Use:   "linkgen",
	Short: "Generates sharelink code for friends to easily download",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			logrus.Error("Too few arguments, size is", len(args))
			return
		}
		transferFile := args[0]
		sha256 := args[1]

		settings, err := config.ReadSettings("settings.yml")
		if err != nil {
			logrus.Error(err)
			return
		}

		logrus.Debug("Connecting to ", settings.GrpcAddress)

		conn, err := fileshare.NewClientConn(settings)
		if err != nil {
			logrus.Panic(err)
		}

		client := sharelink.NewShareLinkClient(cmd.Context(), conn)

		code := client.GenerateLink(transferFile, sha256)
		logrus.Infof("Generated Code is: [%s]", code)

		if err := conn.Close(); err != nil {
			logrus.Error(err)
		}
	},
}
