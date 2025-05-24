package upload

import (
	"context"
	"io"
	"os"

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

func (c *UploadClient) getSummary(ctx context.Context, filePath string) (*pb.UploadSummary, error) {
	task, err := CreateTask(filePath)
	if err != nil {
		return nil, err
	}
	logrus.Debugf("task [filename: %s, sha256: %s, file size: %d]", task.Meta.Filename, task.Meta.Sha256, task.FileSize)
	summary, err := c.Client.PreUpload(ctx, task)
	if err != nil {
		return nil, err
	}
	return summary, nil
}

func (c *UploadClient) UploadFile(ctx context.Context, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}

	summary, err := c.getSummary(ctx, filePath)
	if err != nil {
		return err
	}

	fileName := fileutil.GetFileName(filePath)
	logrus.Debugf("File summary: [filename: %s, sha256: %s, chunk number: %d, chunk size: %d, uploadList: %v]", summary.Meta.Filename, summary.Meta.Sha256, summary.GetChunkNumber(), summary.GetChunkSize(), summary.GetChunkList())

	for _, chunkIndex := range summary.ChunkList {
		chunk := MakeChunk(file, fileName, summary.Meta.Sha256, summary.ChunkSize, summary.ChunkNumber, chunkIndex)

		logrus.Debugf("File Chunk:[filename: %s, sha256: %s, chunk index: %d, chunk size: %d]", summary.Meta.Filename, summary.Meta.Sha256, chunk.GetIndex(), len(chunk.Data))
		err = c.Stream.Send(chunk)

		if err != nil {
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

func MakeChunk(file *os.File, fileName, sha256 string, chunkSize int64, totalChunkNumber, chunkIndex int32) *pb.FileChunk {
	data := make([]byte, chunkSize)
	file.Seek(chunkSize*int64(chunkIndex), 0)
	n, err := file.Read(data)
	if err != nil {
		logrus.Error(err)
	}

	return &pb.FileChunk{
		Meta: &pb.FileMeta{
			Filename: fileName,
			Sha256:   sha256,
		},
		Total: totalChunkNumber,
		Index: chunkIndex,
		Data:  data[:n],
	}
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
		Meta: &pb.FileMeta{
			Filename: fileutil.GetFileName(filePath),
			Sha256:   sha256,
		},
		FileSize: stat.Size(),
	}
	return task, nil
}
