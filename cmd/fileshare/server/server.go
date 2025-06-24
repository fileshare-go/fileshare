package server

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/chanmaoganda/fileshare/cmd/fileshare"
	"github.com/chanmaoganda/fileshare/internal/config"
	"github.com/chanmaoganda/fileshare/internal/core/download"
	"github.com/chanmaoganda/fileshare/internal/core/setup"
	"github.com/chanmaoganda/fileshare/internal/core/sharelink"
	"github.com/chanmaoganda/fileshare/internal/core/upload"
	pb "github.com/chanmaoganda/fileshare/internal/proto/gen"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var ServerCmd = &cobra.Command{
	Use:     "server",
	Short:   "Starts fileshare server",
	PreRunE: setup.SetupServer,
	Run: func(cmd *cobra.Command, args []string) {
		PrintBanner()

		cfg := config.Cfg()

		cfg.PrintConfig()

		logrus.Debug("Server listening on ", cfg.GrpcAddress)

		listen, err := net.Listen("tcp", cfg.GrpcAddress)
		if err != nil {
			logrus.Fatalln("cannot bind address")
		}

		grpcServer, err := fileshare.NewServerConn(cfg)
		if err != nil {
			logrus.Fatal(err)
		}

		pb.RegisterUploadServiceServer(grpcServer, upload.NewUploadServer())
		pb.RegisterDownloadServiceServer(grpcServer, download.NewDownloadServer())
		pb.RegisterShareLinkServiceServer(grpcServer, sharelink.NewShareLinkServer())

		go func() {
			if err := grpcServer.Serve(listen); err != nil {
				logrus.Error(err)
			}
		}()

		// web := web.NewWebService(DB)
		// go func() {
		// 	if err := web.Run(settings.WebAddress); err != nil {
		// 		logrus.Error(err)
		// 	}
		// }()

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
