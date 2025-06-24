package client

import (
	"github.com/chanmaoganda/fileshare/cmd/fileshare"
	"github.com/chanmaoganda/fileshare/internal/config"
	"github.com/chanmaoganda/fileshare/internal/core/download"
	"github.com/chanmaoganda/fileshare/internal/core/setup"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var DownloadCmd = &cobra.Command{
	Use:     "download <checksum256 | linkcode>",
	Short:   "Download file, either with sharelink code or file checksum256 hash",
	Args:    cobra.MinimumNArgs(1),
	PreRunE: setup.SetupClient,
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]

		cfg := config.Cfg()

		logrus.Debug("Connecting to ", cfg.GrpcAddress)

		conn, err := fileshare.NewClientConn(cfg)
		if err != nil {
			logrus.Fatal(err)
		}

		client := download.NewDownloadClient(cmd.Context(), conn)

		if err := client.DownloadFile(cmd.Context(), key); err != nil {
			logrus.Error(err)
		}

		if err := conn.Close(); err != nil {
			logrus.Error(err)
		}
	},
}
