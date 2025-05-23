package client

import (
	"context"

	"github.com/chanmaoganda/fileshare/internal/fileshare/upload"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var UploadCmd = &cobra.Command{
	Use: "upload",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			logrus.Error("Too few arguments, size is", len(args))
			return
		}

		address := "127.0.0.1:60011"
		logrus.Debug("Uploading file to ", address)

		conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))

		if err != nil {
			logrus.Error(err)
		}

		defer conn.Close()

		client := upload.NewUploadClient(conn)

		if err := client.UploadFile(context.Background(), args[0]); err != nil {
			logrus.Error(err)
		}
	},
}
