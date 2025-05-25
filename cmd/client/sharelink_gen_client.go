package client

import (
	"github.com/chanmaoganda/fileshare/internal/config"
	"github.com/chanmaoganda/fileshare/internal/fileshare/sharelink"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var ShareLinkGenCmd = &cobra.Command{
	Use: "linkgen",
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

		logrus.Debug("Uploading file to ", settings.Address)

		conn, err := grpc.NewClient(settings.Address, grpc.WithTransportCredentials(insecure.NewCredentials()))

		if err != nil {
			logrus.Error(err)
		}

		defer conn.Close()

		client := sharelink.NewShareLinkClient(cmd.Context(), conn)

		code := client.GenerateLink(transferFile, sha256)
		logrus.Infof("Generated Code is: [%s]", code)
	},
}
