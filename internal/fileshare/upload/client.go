package upload

import (
	"context"
	"io"
	"os"

	"github.com/chanmaoganda/fileshare/internal/fileshare/chunkio"
	"github.com/chanmaoganda/fileshare/internal/pkg/debugprint"
	"github.com/chanmaoganda/fileshare/internal/pkg/fileutil"
	"github.com/chanmaoganda/fileshare/internal/pkg/sha256"
	pb "github.com/chanmaoganda/fileshare/proto/gen"
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

func (c *UploadClient) UploadFile(ctx context.Context, filePath string) error {
	task, err := c.getTask(ctx, filePath)
	if err != nil {
		return err
	}

	c.uploadWithTask(task, filePath)

	status, err := c.Stream.CloseAndRecv()
	if err != nil && err != io.EOF {
		return err
	}

	logrus.Debugf("[Upload] Status Info [status: %d]", status.Status)
	return nil
}

// use task to upload chunks to server
func (c *UploadClient) uploadWithTask(task *pb.UploadTask, filePath string) {
	debugprint.DebugUploadTask(task)

	if len(task.ChunkList) == 0 {
		// if no chunk is needed, just send the first chunk for messaging
		// at least one chunk is sent cause server side needs meta for recording information
		c.uploadEmptyTask(task)
		return
	}

	file, err := os.Open(filePath)
	if err != nil {
		logrus.Error("File cannot open, upload abort! ERR: ", err)
		return
	}

	for _, chunkIndex := range task.ChunkList {
		chunk := chunkio.MakeChunk(file, task.Meta.Sha256, task.ChunkSize, chunkIndex)

		debugprint.DebugChunk(chunk)

		if err := c.Stream.Send(chunk); err != nil {
			logrus.Error(err)
			break
		}
	}
	logrus.Debug("[Upload] Upload done")
}

func (c *UploadClient) uploadEmptyTask(task *pb.UploadTask) {
	logrus.Debug("Upload Task is empty, just send empty data instead")
	chunk := &pb.FileChunk{
		Sha256:     task.Meta.Sha256,
		ChunkIndex: 0,
		Data:       []byte{},
	}
	if err := c.Stream.Send(chunk); err != nil {
		logrus.Error(err)
	}
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
