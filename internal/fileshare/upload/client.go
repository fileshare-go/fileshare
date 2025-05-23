package upload

import (
	"context"
	"io"
	"os"

	"github.com/chanmaoganda/fileshare/internal/sha256"
	pb "github.com/chanmaoganda/fileshare/proto/upload"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type UploadClient struct {
	Client pb.UploadServiceClient
}

func NewUploadClient(conn *grpc.ClientConn) *UploadClient {
	client := pb.NewUploadServiceClient(conn)
	return &UploadClient{
		Client: client,
	}
}

func (c *UploadClient) UploadFile(ctx context.Context, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}

	task, err := CreateTask(filePath)
	if err != nil {
		return err
	}

	summary, err := c.Client.PreUpload(ctx, task)
	if err != nil {
		return err
	}

	stream, err := c.Client.Upload(ctx)
	if err != nil {
		return err
	}

	logrus.Debugf("File summary: [filename: %s, sha256: %s, chunk number: %d, chunk size: %d, uploadList: %v]", summary.GetFilename(), summary.GetSha256(), summary.GetChunkNumber(), summary.GetChunkSize(), summary.GetChunkList())

	for _, chunkIndex := range summary.ChunkList {
		data := make([]byte, summary.ChunkSize)
		file.Seek(summary.ChunkSize * int64(chunkIndex), 0)
		n, err := file.Read(data)
		if err != nil {
			break
		}

		logrus.Debugf("File Chunk:[filename: %s, chunk index: %d, chunk size: %d]", summary.GetFilename(), chunkIndex, n)
		err = stream.Send(&pb.FileChunk{
			Filename: filePath,
			Index:    chunkIndex,
			Data:     data,
		})

		if err != nil {
			break
		}
	}

	status, err := stream.CloseAndRecv()
	if err != nil && err != io.EOF {
		return err
	}

	logrus.Debugf("Status Info [status: %d]", status.Status)
	return nil
}

func CreateTask(filePath string) (*pb.UploadTask, error) {
	stat, err := os.Stat(filePath)
	if err != nil {
		return nil, err
	}

	sha256, err := sha256.CalculateSHA256(filePath)
	if err != nil {
		return nil, err
	}

	task := &pb.UploadTask{
		Filename: filePath,
		FileSize: stat.Size(),
		Sha256: sha256,
	}
	return task, nil
}
