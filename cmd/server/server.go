package server

import (
	"net"

	"github.com/chanmaoganda/fileshare/internal/config"
	"github.com/chanmaoganda/fileshare/internal/db"
	"github.com/chanmaoganda/fileshare/internal/fileshare/download"
	"github.com/chanmaoganda/fileshare/internal/fileshare/sharelink"
	"github.com/chanmaoganda/fileshare/internal/fileshare/upload"
	pb "github.com/chanmaoganda/fileshare/proto/gen"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

var ServerCmd = &cobra.Command{
	Use: "server",
	Run: func(cmd *cobra.Command, args []string) {
		settings, err := config.ReadSettings("settings.yml")
		if err != nil {
			logrus.Error(err)
			return
		}

		listen, err := net.Listen("tcp", settings.Address)

		logrus.Debug("Server listening on ", settings.Address)

		if err != nil {
			logrus.Fatalln("cannot bind address")
		}

		grpcServer := grpc.NewServer()

		DB := db.SetupDB(settings.Database)

		pb.RegisterUploadServiceServer(grpcServer, upload.NewUploadServer(settings, DB))
		pb.RegisterDownloadServiceServer(grpcServer, download.NewDownloadServer(settings, DB))
		pb.RegisterShareLinkServiceServer(grpcServer, sharelink.NewShareLinkServer(settings, DB))

		if err := grpcServer.Serve(listen); err != nil {
			logrus.Error(err)
		}
	},
}
