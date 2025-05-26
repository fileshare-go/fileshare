package download

import (
	"context"

	"github.com/chanmaoganda/fileshare/internal/fileshare/dbmanager"
	"github.com/chanmaoganda/fileshare/internal/model"
	pb "github.com/chanmaoganda/fileshare/proto/gen"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

type DownloadClient struct {
	Client  pb.DownloadServiceClient
	Manager *dbmanager.DBManager
}

func NewDownloadClient(ctx context.Context, conn *grpc.ClientConn, DB *gorm.DB) *DownloadClient {
	client := pb.NewDownloadServiceClient(conn)

	return &DownloadClient{
		Client:  client,
		Manager: dbmanager.NewDBManager(DB),
	}
}

func (c *DownloadClient) getSummary(ctx context.Context, key string) (*pb.DownloadSummary, error) {
	// if the key is not the fixed size of sha256, then recognize this as link code
	if len(key) != 64 {
		return c.Client.PreDownloadWithCode(ctx, &pb.ShareLink{LinkCode: key})
	} else {
		return c.Client.PreDownload(ctx, &pb.DownloadRequest{Meta: &pb.FileMeta{Sha256: key}})
	}
}

func (c *DownloadClient) getTask(ctx context.Context, key string) (*pb.DownloadTask, error) {
	summary, err := c.getSummary(ctx, key)
	if err != nil {
		return nil, err
	}

	fileInfo := &model.FileInfo{
		Sha256: summary.Meta.Sha256,
	}

	if c.Manager.SelectFileInfo(fileInfo) {
		return fileInfo.BuildDownloadTask(), nil
	}

	fileInfo = model.NewFileInfoFromDownload(summary)

	c.Manager.CreateFileInfo(fileInfo)

	task := &pb.DownloadTask{
		Meta:        summary.Meta,
		ChunkNumber: summary.ChunkNumber,
		ChunkList:   summary.ChunkList,
	}
	return task, nil
}

func (c *DownloadClient) downloadStream(ctx context.Context, key string) (pb.DownloadService_DownloadClient, error) {
	logrus.Debugf("Download request [key: %s]", key)

	task, err := c.getTask(ctx, key)
	if err != nil {
		return nil, err
	}

	return c.Client.Download(ctx, task)
}

func (c *DownloadClient) DownloadFile(ctx context.Context, key string) error {
	stream, err := c.downloadStream(ctx, key)

	if err != nil {
		return err
	}

	handler := NewHandler(stream, c.Manager)

	// if recv or saving has any err, just close and return err
	if err := handler.Recv(); err != nil {
		return handler.CloseWithErr(err)
	}

	// if recv and saving do not has any error, validate and close
	handler.ValidateAndClose()

	logrus.Debug("[Download] Ending Download Process!")
	return nil
}
