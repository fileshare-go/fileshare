package download

import (
	"context"

	"github.com/chanmaoganda/fileshare/internal/fileshare/model"
	pb "github.com/chanmaoganda/fileshare/proto/gen"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

type DownloadClient struct {
	Client pb.DownloadServiceClient
	DB     *gorm.DB
}

func NewDownloadClient(ctx context.Context, conn *grpc.ClientConn, DB *gorm.DB) *DownloadClient {
	client := pb.NewDownloadServiceClient(conn)

	return &DownloadClient{
		Client: client,
		DB:     DB,
	}
}

func (c *DownloadClient) getTask(ctx context.Context, sha256 string) (*pb.DownloadTask, error) {
	summary, err := c.Client.PreDownload(ctx, &pb.DownloadRequest{Meta: &pb.FileMeta{Sha256: sha256}})
	if err != nil {
		return nil, err
	}

	fileInfo, ok := model.GetFileInfo(sha256, c.DB)
	if ok {
		return fileInfo.BuildDownloadTask(), nil
	}

	fileInfo = model.NewFileInfoFromDownload(summary)

	c.DB.Create(fileInfo)

	task := &pb.DownloadTask{
		Meta:        summary.Meta,
		ChunkNumber: summary.ChunkNumber,
		ChunkList:   summary.ChunkList,
	}
	return task, err
}

func (c *DownloadClient) PreDownload(ctx context.Context, sha256 string) (pb.DownloadService_DownloadClient, error) {
	logrus.Debugf("Download request [sha256: %s]", sha256)

	task, err := c.getTask(ctx, sha256)
	if err != nil {
		return nil, err
	}

	return c.Client.Download(ctx, task)
}

func (c *DownloadClient) DownloadFile(ctx context.Context, sha256 string) error {
	stream, err := c.PreDownload(ctx, sha256)

	if err != nil {
		return err
	}

	handler := NewHandler(stream, c.DB)

	// if recv or saving has any err, just close and return err
	if err := handler.Recv(); err != nil {
		return handler.CloseWithErr(err)
	}

	// if recv and saving do not has any error, validate and close
	handler.ValidateAndClose()

	logrus.Debug("[Download] Ending Download Process!")
	return nil
}
