package upload

import (
	"context"
	"os"

	"github.com/chanmaoganda/fileshare/internal/core/chunkstream/send"
	"github.com/chanmaoganda/fileshare/internal/pkg/util"
	pb "github.com/chanmaoganda/fileshare/internal/proto/gen"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type UploadClient struct {
	Client pb.UploadServiceClient
	Stream pb.UploadService_UploadClient
}

func NewUploadClient(ctx context.Context, conn *grpc.ClientConn) *UploadClient {
	client := pb.NewUploadServiceClient(conn)
	stream, err := client.Upload(ctx)
	if err != nil {
		logrus.Panic(err)
	}

	return &UploadClient{
		Client: client,
		Stream: stream,
	}
}

func (c *UploadClient) UploadFile(ctx context.Context, filePath string) error {
	builder := TaskBuilder{Client: c.Client}

	task, err := builder.GetTask(ctx, filePath)
	if err != nil {
		return err
	}
	util.DebugUploadTask(task)

	sendStream := send.NewClientSendStream(task, filePath, c.Stream)

	if err := sendStream.SendStreamChunks(); err != nil {
		return err
	}

	logrus.Debug("[Upload] Upload done")
	return sendStream.CloseStream()
}

type TaskBuilder struct {
	Client pb.UploadServiceClient
}

// build upload request
func (b *TaskBuilder) BuildRequest(filePath string) (*pb.UploadRequest, error) {
	stat, err := os.Stat(filePath)
	if err != nil {
		return nil, err
	}

	sha256, err := util.CalculateFileSHA256(filePath)
	if err != nil {
		return nil, err
	}

	request := &pb.UploadRequest{
		Meta: &pb.FileMeta{
			Filename: util.GetFileName(filePath),
			Sha256:   sha256,
		},
		FileSize: stat.Size(),
	}
	return request, nil
}

// recv
func (b *TaskBuilder) GetTask(ctx context.Context, filePath string) (*pb.UploadTask, error) {
	request, err := b.BuildRequest(filePath)
	if err != nil {
		return nil, err
	}
	logrus.Debugf("request [filename: %s, sha256: %s, file size: %d]", request.Meta.Filename, request.Meta.Sha256, request.FileSize)

	task, err := b.Client.PreUpload(ctx, request)
	if err != nil {
		return nil, err
	}
	return task, nil
}
