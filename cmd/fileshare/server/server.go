package server

import (
	"fmt"
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
		PrintBanner()

		settings, err := config.ReadSettings("settings.yml")
		if err != nil {
			logrus.Error(err)
			return
		}

		settings.PrintSettings()

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

func PrintBanner() {
	banner := []byte{32, 32, 32, 32, 95, 95, 95, 95, 95, 32, 95, 95, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 95, 95, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 10, 32, 32, 32, 47, 32, 95, 95, 40, 95, 41, 32, 47, 95, 95, 32, 32, 32, 32, 32, 95, 95, 95, 95, 95, 47, 32, 47, 95, 32, 32, 95, 95, 95, 95, 32, 95, 95, 95, 95, 95, 95, 95, 95, 95, 32, 10, 32, 32, 47, 32, 47, 95, 47, 32, 47, 32, 47, 32, 95, 32, 92, 32, 32, 32, 47, 32, 95, 95, 95, 47, 32, 95, 95, 32, 92, 47, 32, 95, 95, 32, 96, 47, 32, 95, 95, 95, 47, 32, 95, 32, 92, 10, 32, 47, 32, 95, 95, 47, 32, 47, 32, 47, 32, 32, 95, 95, 47, 32, 32, 40, 95, 95, 32, 32, 41, 32, 47, 32, 47, 32, 47, 32, 47, 95, 47, 32, 47, 32, 47, 32, 32, 47, 32, 32, 95, 95, 47, 10, 47, 95, 47, 32, 47, 95, 47, 95, 47, 92, 95, 95, 95, 47, 32, 32, 47, 95, 95, 95, 95, 47, 95, 47, 32, 47, 95, 47, 92, 95, 95, 44, 95, 47, 95, 47, 32, 32, 32, 92, 95, 95, 95, 47, 32, 10}
	fmt.Printf("\n%s\n\n\n", string(banner))
}
