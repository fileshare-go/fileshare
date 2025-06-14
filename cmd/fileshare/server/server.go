package server

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/chanmaoganda/fileshare/cmd/fileshare"
	"github.com/chanmaoganda/fileshare/internal/config"
	"github.com/chanmaoganda/fileshare/internal/fileshare/download"
	"github.com/chanmaoganda/fileshare/internal/fileshare/sharelink"
	"github.com/chanmaoganda/fileshare/internal/fileshare/upload"
	"github.com/chanmaoganda/fileshare/internal/pkg/db"
	pb "github.com/chanmaoganda/fileshare/internal/proto/gen"
	"github.com/chanmaoganda/fileshare/internal/web"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var ServerCmd = &cobra.Command{
	Use:   "server",
	Short: "Starts fileshare server",
	Run: func(cmd *cobra.Command, args []string) {
		PrintBanner()

		settings, err := config.ReadSettings("settings.yml")
		if err != nil {
			logrus.Error(err)
			return
		}

		settings.PrintSettings()

		logrus.Debug("Server listening on ", settings.GrpcAddress)

		listen, err := net.Listen("tcp", settings.GrpcAddress)
		if err != nil {
			logrus.Fatalln("cannot bind address")
		}

		grpcServer, err := fileshare.NewServerConn(settings)
		if err != nil {
			logrus.Panic(err)
		}

		DB := db.SetupServerDB(settings.Database)

		pb.RegisterUploadServiceServer(grpcServer, upload.NewUploadServer(settings, DB))
		pb.RegisterDownloadServiceServer(grpcServer, download.NewDownloadServer(settings, DB))
		pb.RegisterShareLinkServiceServer(grpcServer, sharelink.NewShareLinkServer(settings, DB))

		go func() {
			if err := grpcServer.Serve(listen); err != nil {
				logrus.Error(err)
			}
		}()

		web := web.NewWebService(DB)
		go func() {
			if err := web.Run(settings.WebAddress); err != nil {
				logrus.Error(err)
			}
		}()

		shutdown := make(chan os.Signal, 1)
		signal.Notify(shutdown, os.Interrupt, syscall.SIGINT)
		<-shutdown
		logrus.Info("Shutting down servers...")
		grpcServer.GracefulStop()
	},
}

func PrintBanner() {
	banner := []byte{32, 32, 32, 32, 95, 95, 95, 95, 95, 32, 95, 95, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 95, 95, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 10, 32, 32, 32, 47, 32, 95, 95, 40, 95, 41, 32, 47, 95, 95, 32, 32, 32, 32, 32, 95, 95, 95, 95, 95, 47, 32, 47, 95, 32, 32, 95, 95, 95, 95, 32, 95, 95, 95, 95, 95, 95, 95, 95, 95, 32, 10, 32, 32, 47, 32, 47, 95, 47, 32, 47, 32, 47, 32, 95, 32, 92, 32, 32, 32, 47, 32, 95, 95, 95, 47, 32, 95, 95, 32, 92, 47, 32, 95, 95, 32, 96, 47, 32, 95, 95, 95, 47, 32, 95, 32, 92, 10, 32, 47, 32, 95, 95, 47, 32, 47, 32, 47, 32, 32, 95, 95, 47, 32, 32, 40, 95, 95, 32, 32, 41, 32, 47, 32, 47, 32, 47, 32, 47, 95, 47, 32, 47, 32, 47, 32, 32, 47, 32, 32, 95, 95, 47, 10, 47, 95, 47, 32, 47, 95, 47, 95, 47, 92, 95, 95, 95, 47, 32, 32, 47, 95, 95, 95, 95, 47, 95, 47, 32, 47, 95, 47, 92, 95, 95, 44, 95, 47, 95, 47, 32, 32, 32, 92, 95, 95, 95, 47, 32, 10}
	fmt.Printf("\n%s\n\n\n", string(banner))
}
