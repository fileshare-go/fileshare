package fileshare

import (
	"github.com/chanmaoganda/fileshare/internal/config"
	"github.com/chanmaoganda/fileshare/internal/core/download"
	"github.com/chanmaoganda/fileshare/internal/core/sharelink"
	"github.com/chanmaoganda/fileshare/internal/core/upload"
	"github.com/chanmaoganda/fileshare/internal/pkg/grpc"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

func Upload(ctx context.Context, args []string) {
	transferFile := args[0]

	cfg := config.Cfg()

	logrus.Debug("Uploading file to ", cfg.GrpcAddress)

	conn, err := grpc.NewClientConn(cfg)
	if err != nil {
		logrus.Fatal(err)
	}

	client := upload.NewUploadClient(ctx, conn)

	if err := client.UploadFile(ctx, transferFile); err != nil {
		logrus.Error(err)
	}

	if err := conn.Close(); err != nil {
		logrus.Error(err)
	}
}

func LinkGen(ctx context.Context, args []string) {
	cfg := config.Cfg()

	logrus.Debug("Connecting to ", cfg.GrpcAddress)

	conn, err := grpc.NewClientConn(cfg)
	if err != nil {
		logrus.Fatal(err)
	}

	client := sharelink.NewShareLinkClient(ctx, conn)

	code := client.GenerateLink(args)
	logrus.Infof("Generated Code is: [%s]", code)

	if err := conn.Close(); err != nil {
		logrus.Error(err)
	}
}

func Download(ctx context.Context, args []string) {
	key := args[0]

	cfg := config.Cfg()

	logrus.Debug("Connecting to ", cfg.GrpcAddress)

	conn, err := grpc.NewClientConn(cfg)
	if err != nil {
		logrus.Fatal(err)
	}

	client := download.NewDownloadClient(ctx, conn)

	if err := client.DownloadFile(ctx, key); err != nil {
		logrus.Error(err)
	}

	if err := conn.Close(); err != nil {
		logrus.Error(err)
	}
}
