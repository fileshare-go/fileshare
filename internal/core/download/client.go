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

// get summary from grpc
func (c *DownloadClient) getSummary(ctx context.Context, key string) (*pb.DownloadSummary, error) {
	// if the key is not the fixed size of sha256, then recognize this as link code
	if len(key) != 64 {
		return c.Client.PreDownloadWithCode(ctx, &pb.ShareLink{LinkCode: key})
	} else {
		return c.Client.PreDownload(ctx, &pb.DownloadRequest{Meta: &pb.FileMeta{Sha256: key}})
	}
}

// get the download stream from grpc, aimed to handle download tasks
func (c *DownloadClient) getDownloadStream(ctx context.Context, key string) (pb.DownloadService_DownloadClient, error) {
	logrus.Debugf("Download request [key: %s]", key)
	var err error
	summary, err := c.getSummary(ctx, key)
	if err != nil {
		return nil, err
	}

	// if summary is not OK, then return with failure message
	if summary.Status != pb.Status_OK {
		return nil, errors.New(summary.Message)
	}

	task, err := buildTaskFromSummary(summary)
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

// build download task from summary
func buildTaskFromSummary(summary *pb.DownloadSummary) (*pb.DownloadTask, error) {
	var err error
	fileInfo := &model.FileInfo{
		Sha256: summary.Meta.Sha256,
	}

	if service.Orm().Find(fileInfo).RowsAffected == 1 {
		return assembleDownloadTask(fileInfo), nil
	}

	fileInfo = assembleFileInfo(summary)

	if service.Orm().Save(fileInfo).Error != nil {
		return nil, err
	}

	return &pb.DownloadTask{
		Meta:        summary.Meta,
		ChunkNumber: summary.ChunkNumber,
		ChunkList:   summary.ChunkList,
	}, nil
}

func assembleDownloadTask(f *model.FileInfo) *pb.DownloadTask {
	return &pb.DownloadTask{
		Meta: &pb.FileMeta{
			Filename: f.Filename,
			Sha256:   f.Sha256,
			FileSize: f.FileSize,
		},
		ChunkNumber: f.ChunkNumber,
		ChunkList:   f.GetMissingChunks(),
	}
}

func assembleFileInfo(summary *pb.DownloadSummary) *model.FileInfo {
	fileInfo := model.FileInfo{}

	fileInfo.Filename = summary.Meta.Filename
	fileInfo.Sha256 = summary.Meta.Sha256
	fileInfo.FileSize = summary.FileSize
	fileInfo.ChunkNumber = summary.ChunkNumber
	fileInfo.ChunkSize = summary.ChunkSize
	fileInfo.UploadedChunks = "[]"

	return &fileInfo
}

