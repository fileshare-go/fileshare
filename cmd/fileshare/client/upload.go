package client

import (
	"github.com/chanmaoganda/fileshare/cmd/fileshare"
	"github.com/chanmaoganda/fileshare/internal/config"
	"github.com/chanmaoganda/fileshare/internal/core/setup"
	"github.com/chanmaoganda/fileshare/internal/core/upload"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var UploadCmd = &cobra.Command{
	Use:     "upload <filepath>",
	Short:   "Uploads the file, requires the filepath as argument",
	Args:    cobra.MinimumNArgs(1),
	PreRunE: setup.Setup,
	Run: func(cmd *cobra.Command, args []string) {
		transferFile := args[0]

		cfg := config.Cfg()

		logrus.Debug("Uploading file to ", cfg.GrpcAddress)

		conn, err := fileshare.NewClientConn(cfg)
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
