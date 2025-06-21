package download

import (
	"context"
	"errors"

	"github.com/chanmaoganda/fileshare/internal/core/chunkstream/recv"
	"github.com/chanmaoganda/fileshare/internal/model"
	pb "github.com/chanmaoganda/fileshare/internal/proto/gen"
	"github.com/chanmaoganda/fileshare/internal/service"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type DownloadClient struct {
	Client pb.DownloadServiceClient
}

func NewDownloadClient(ctx context.Context, conn *grpc.ClientConn) *DownloadClient {
	client := pb.NewDownloadServiceClient(conn)

	return &DownloadClient{
		Client: client,
	}
}

// get the download stream from grpc, aimed to handle download tasks
func (c *DownloadClient) getDownloadStream(ctx context.Context, key string) (pb.DownloadService_DownloadClient, error) {
	logrus.Debugf("Download request [key: %s]", key)

	builder := taskBuilder{Client: c.Client}

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

	recvStream := recv.NewClientRecvStream(stream)
	if err := recvStream.RecvStreamChunks(); err != nil {
		return recvStream.CloseStream(false)
	}

	validate := recvStream.ValidateRecvChunks()
	return recvStream.CloseStream(validate)
}

type taskBuilder struct {
	Client pb.DownloadServiceClient
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
	var err error
	fileInfo := &model.FileInfo{
		Sha256: summary.Meta.Sha256,
	}

	if err = service.Mgr().SelectFileInfo(fileInfo); err == nil {
		return fileInfo.BuildDownloadTask(), nil
	}

	fileInfo = model.NewFileInfoFromDownload(summary)

	if err = service.Mgr().InsertFileInfo(fileInfo); err != nil {
		return nil, err
	}

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
