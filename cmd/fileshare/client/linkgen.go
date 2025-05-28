package client

import (
	"strconv"

	"github.com/chanmaoganda/fileshare/cmd/fileshare"
	"github.com/chanmaoganda/fileshare/internal/config"
	"github.com/chanmaoganda/fileshare/internal/fileshare/sharelink"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var ShareLinkGenCmd = &cobra.Command{
	Use:   "linkgen <filename> <checksum256> <expire days>",
	Short: "Generates sharelink code for friends to easily download",
	Args:  cobra.MinimumNArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		settings, err := config.ReadSettings("settings.yml")
		if err != nil {
			logrus.Error(err)
			return
		}

		transferFile := args[0]
		sha256 := args[1]
		validDays, err := strconv.Atoi(args[2])
		if err != nil {
			validDays = settings.ValidDays
		}

		logrus.Debug("Connecting to ", settings.GrpcAddress)

		conn, err := fileshare.NewClientConn(settings)
		if err != nil {
			logrus.Panic(err)
		}

		client := sharelink.NewShareLinkClient(cmd.Context(), conn)

		code := client.GenerateLink(transferFile, sha256, validDays)
		logrus.Infof("Generated Code is: [%s]", code)

		if err := conn.Close(); err != nil {
			logrus.Error(err)
		}
	},
}
