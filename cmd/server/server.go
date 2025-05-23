package server

import (
	"net"

	"github.com/chanmaoganda/fileshare/pkg/fileshare/upload"
	pb "github.com/chanmaoganda/fileshare/proto/upload"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

var ServerCmd = &cobra.Command{
	Use: "server",
	Run: func(cmd *cobra.Command, args []string) {
		address := "127.0.0.1:60011"

		listen, err := net.Listen("tcp", address)

		logrus.Debug("Server listening on ", address)

		if err != nil {
			logrus.Fatalln("cannot bind address")
		}

		grpcServer := grpc.NewServer()

		pb.RegisterUploadServiceServer(grpcServer, &upload.UploadServer{})

		err = grpcServer.Serve(listen)
	},
}
