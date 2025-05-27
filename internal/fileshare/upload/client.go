package upload

import (
	"context"
	"os"

	"github.com/chanmaoganda/fileshare/internal/fileshare/chunkstream/send"
	"github.com/chanmaoganda/fileshare/internal/pkg/debugprint"
	"github.com/chanmaoganda/fileshare/internal/pkg/fileutil"
	"github.com/chanmaoganda/fileshare/internal/pkg/sha256"
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
	task, err := c.getTask(ctx, filePath)
	if err != nil {
		return err
	}

	if err := c.uploadWithTask(task, filePath); err != nil {
		return err
	}

	logrus.Debug("[Upload] Upload done")
	return nil
}

func (c *UploadClient) getTask(ctx context.Context, filePath string) (*pb.UploadTask, error) {
	request, err := c.createRequest(filePath)
	if err != nil {
		return nil, err
	}
	logrus.Debugf("request [filename: %s, sha256: %s, file size: %d]", request.Meta.Filename, request.Meta.Sha256, request.FileSize)
	task, err := c.Client.PreUpload(ctx, request)
	if err != nil {
		return nil, err
	}
	return task, nil
}

// use task to upload chunks to server
func (c *UploadClient) uploadWithTask(task *pb.UploadTask, filePath string) error {
	debugprint.DebugUploadTask(task)

	sendStream := send.NewClientSendStream(task, filePath, c.Stream)

	if err := sendStream.SendStreamChunks(); err != nil {
		return err
	}

	return sendStream.CloseStream()
}

// create upload request
func (c *UploadClient) createRequest(filePath string) (*pb.UploadRequest, error) {
	stat, err := os.Stat(filePath)
	if err != nil {
		return nil, err
	}

	sha256, err := sha256.CalculateFileSHA256(filePath)
	if err != nil {
		return nil, err
	}

	request := &pb.UploadRequest{
		Meta: &pb.FileMeta{
			Filename: fileutil.GetFileName(filePath),
			Sha256:   sha256,
		},
		FileSize: stat.Size(),
	}
	return request, nil
}
