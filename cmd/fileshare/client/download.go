package client

import (
	"github.com/chanmaoganda/fileshare/cmd/fileshare"
	"github.com/chanmaoganda/fileshare/internal/core/download"
	"github.com/chanmaoganda/fileshare/internal/config"
	"github.com/chanmaoganda/fileshare/pkg/db"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var DownloadCmd = &cobra.Command{
	Use:   "download <checksum256 | linkcode>",
	Short: "Download file, either with sharelink code or file checksum256 hash",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]

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

		DB := db.SetupClientDB(settings.Database)

		client := download.NewDownloadClient(cmd.Context(), settings, conn, DB)

		if err := client.DownloadFile(cmd.Context(), key); err != nil {
			logrus.Error(err)
		}

		if err := conn.Close(); err != nil {
			logrus.Error(err)
		}
	},
}
