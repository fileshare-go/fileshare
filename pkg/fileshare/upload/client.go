package upload

import (
	"context"

	pb "github.com/chanmaoganda/fileshare/proto/upload"
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

func (c *UploadClient) UploadFile(filename string) (*pb.UploadSummary, error) {
	task := &pb.UploadTask{
		Filename: filename,
		FileSize: 128,
	}

	summary, err := c.Client.PreUpload(context.Background(), task)
	if err != nil {
		return nil, err
	}

	return summary, nil
}
