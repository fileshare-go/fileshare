package download

import (
	"context"
	"errors"

	"github.com/chanmaoganda/fileshare/internal/config"
	"github.com/chanmaoganda/fileshare/internal/fileshare/chunkstream/recv"
	"github.com/chanmaoganda/fileshare/internal/model"
	"github.com/chanmaoganda/fileshare/internal/pkg/dbmanager"
	pb "github.com/chanmaoganda/fileshare/internal/proto/gen"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

type DownloadClient struct {
	Settings *config.Settings
	Client   pb.DownloadServiceClient
	Manager  *dbmanager.DBManager
}

func NewDownloadClient(ctx context.Context, settings *config.Settings, conn *grpc.ClientConn, DB *gorm.DB) *DownloadClient {
	client := pb.NewDownloadServiceClient(conn)

	return &DownloadClient{
		Settings: settings,
		Client:   client,
		Manager:  dbmanager.NewDBManager(DB),
	}
}

// get the download stream from grpc, aimed to handle download tasks
func (c *DownloadClient) getDownloadStream(ctx context.Context, key string) (pb.DownloadService_DownloadClient, error) {
	logrus.Debugf("Download request [key: %s]", key)

	builder := taskBuilder{Client: c.Client, Manager: c.Manager}

	task, err := builder.BuildTask(ctx, key)
	if err != nil {
		return nil, err
	}

	return c.Client.Download(ctx, task)
}

// download file from the key, both sha256 and sharelink are accepted
func (c *DownloadClient) DownloadFile(ctx context.Context, key string) error {
	stream, err := c.getDownloadStream(ctx, key)
	if err != nil {
		return err
	}

	recvStream := recv.NewClientRecvStream(c.Settings, c.Manager, stream)
	if err := recvStream.RecvStreamChunks(); err != nil {
		return recvStream.CloseStream(false)
	}

	validate := recvStream.ValidateRecvChunks()
	return recvStream.CloseStream(validate)
}

type taskBuilder struct {
	Client  pb.DownloadServiceClient
	Manager *dbmanager.DBManager
}

// get summary from grpc
func (b *taskBuilder) getSummary(ctx context.Context, key string) (*pb.DownloadSummary, error) {
	// if the key is not the fixed size of sha256, then recognize this as link code
	if len(key) != 64 {
		return b.Client.PreDownloadWithCode(ctx, &pb.ShareLink{LinkCode: key})
	} else {
		return b.Client.PreDownload(ctx, &pb.DownloadRequest{Meta: &pb.FileMeta{Sha256: key}})
	}
}

// build download task from summary
func (b *taskBuilder) buildTaskFromSummary(summary *pb.DownloadSummary) (*pb.DownloadTask, error) {
	fileInfo := &model.FileInfo{
		Sha256: summary.Meta.Sha256,
	}

	if b.Manager.SelectFileInfo(fileInfo) {
		return fileInfo.BuildDownloadTask(), nil
	}

	fileInfo = model.NewFileInfoFromDownload(summary)

	b.Manager.CreateFileInfo(fileInfo)

	task := &pb.DownloadTask{
		Meta:        summary.Meta,
		ChunkNumber: summary.ChunkNumber,
		ChunkList:   summary.ChunkList,
	}
	return task, nil
}

func (b *taskBuilder) BuildTask(ctx context.Context, key string) (*pb.DownloadTask, error) {
	summary, err := b.getSummary(ctx, key)
	if err != nil {
		return nil, err
	}
	if summary.Status != pb.Status_OK {
		return nil, errors.New(summary.Message)
	}

	return b.buildTaskFromSummary(summary)
}
