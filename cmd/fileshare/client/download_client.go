package client

import (
	"github.com/chanmaoganda/fileshare/internal/config"
	"github.com/chanmaoganda/fileshare/internal/db"
	"github.com/chanmaoganda/fileshare/internal/fileshare/download"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var DownloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download file, either with sharelink code or file checksum256 hash",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			logrus.Error("Too few arguments, size is", len(args))
			return
		}
		key := args[0]

		settings, err := config.ReadSettings("settings.yml")
		if err != nil {
			logrus.Error(err)
			return
		}

		logrus.Debug("Connecting to ", settings.Address)

		conn, err := grpc.NewClient(settings.Address, grpc.WithTransportCredentials(insecure.NewCredentials()))

		if err != nil {
			logrus.Error(err)
		}

		defer conn.Close()

		DB := db.SetupDB(settings.Database)

		client := download.NewDownloadClient(cmd.Context(), settings, conn, DB)

		if err := client.DownloadFile(cmd.Context(), key); err != nil {
			logrus.Error(err)
		}
	},
}
