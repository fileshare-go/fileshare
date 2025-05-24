package upload

import (
	"context"
	"io"
	"os"

	"github.com/chanmaoganda/fileshare/internal/fileshare/chunkio"
	"github.com/chanmaoganda/fileshare/internal/fileutil"
	"github.com/chanmaoganda/fileshare/internal/sha256"
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
	request, err := CreateRequest(filePath)
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

// TODO: refactor
func (c *UploadClient) UploadFile(ctx context.Context, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}

	task, err := c.getTask(ctx, filePath)
	if err != nil {
		return err
	}

	// fileName := fileutil.GetFileName(filePath)
	logrus.Debugf("File task: [filename: %s, sha256: %s, chunk number: %d, chunk size: %d, uploadList: %v]", task.Meta.Filename, task.Meta.Sha256, task.GetChunkNumber(), task.GetChunkSize(), task.GetChunkList())

	if len(task.ChunkList) == 0 {
		// if no chunk is needed, just send the first chunk for messaging
		// at least one chunk is sent cause server side needs meta for recording information
		task.ChunkList = append(task.ChunkList, 0)
	}

	for _, chunkIndex := range task.ChunkList {
		chunk := chunkio.MakeChunk(file, task.Meta.Sha256, task.ChunkSize, chunkIndex)

		logrus.Debugf("File Chunk:[filename: %s, sha256: %s, chunk index: %d, chunk size: %d]", task.Meta.Filename, task.Meta.Sha256, chunk.ChunkIndex, len(chunk.Data))

		if err := c.Stream.Send(chunk); err != nil {
			logrus.Error(err)
			break
		}
	}

	status, err := c.Stream.CloseAndRecv()
	if err != nil && err != io.EOF {
		return err
	}

	logrus.Debugf("Status Info [status: %d]", status.Status)
	return nil
}

func CreateRequest(filePath string) (*pb.UploadRequest, error) {
	stat, err := os.Stat(filePath)
	if err != nil {
		return nil, err
	}

	sha256, err := sha256.CalculateSHA256(filePath)
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
