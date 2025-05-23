package upload

import (
	"context"
	"io"

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

func (c *UploadClient) UploadFile(ctx context.Context, filename string) error {
	task := &pb.UploadTask{
		Filename: filename,
		FileSize: 128,
	}

	summary, err := c.Client.PreUpload(ctx, task)
	if err != nil {
		return err
	}

	stream, err := c.Client.Upload(ctx)
	if err != nil {
		return err
	}

	logrus.Debugf("File summary: [number: %d, size: %d]", summary.GetChunkNumber(), summary.GetChunkSize())

	for i := range 3 {
		err = stream.Send(&pb.FileChunk{
			Filename: filename,
			Index: int32(i),
			Data: []byte{},
		})

		if err != nil {
			return err
		}
	}

	status, err := stream.CloseAndRecv()
	if err != nil && err != io.EOF {
		return err
	}

	logrus.Debugf("Status Info [status: %d]", status.Status)
	return nil
}
