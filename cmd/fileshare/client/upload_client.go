package client

import (
	"github.com/chanmaoganda/fileshare/cmd/fileshare"
	"github.com/chanmaoganda/fileshare/internal/config"
	"github.com/chanmaoganda/fileshare/internal/fileshare/upload"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var UploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Uploads the file, requires the filename as argument",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			logrus.Error("Too few arguments, size is", len(args))
			return
		}
		transferFile := args[0]

		settings, err := config.ReadSettings("settings.yml")
		if err != nil {
			logrus.Error(err)
			return
		}

		logrus.Debug("Uploading file to ", settings.GrpcAddress)

		conn, err := fileshare.NewClientConn(settings)
		if err != nil {
			logrus.Panic(err)
		}

		client := upload.NewUploadClient(cmd.Context(), conn)

		if err := client.UploadFile(cmd.Context(), transferFile); err != nil {
			logrus.Error(err)
		}

		if err := conn.Close(); err != nil {
			logrus.Error(err)
		}
	},
}
