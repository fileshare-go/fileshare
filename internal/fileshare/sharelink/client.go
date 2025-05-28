package sharelink

import (
	"context"

	pb "github.com/chanmaoganda/fileshare/internal/proto/gen"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type ShareLinkClient struct {
	Client pb.ShareLinkServiceClient
}

func NewShareLinkClient(ctx context.Context, conn *grpc.ClientConn) *ShareLinkClient {
	client := pb.NewShareLinkServiceClient(conn)

	return &ShareLinkClient{
		Client: client,
	}
}

func (c *ShareLinkClient) GenerateLink(filename, sha256 string, validDays int) string {
	req := &pb.ShareLinkRequest{
		Meta: &pb.FileMeta{
			Filename: filename,
			Sha256:   sha256,
		},
		ValidDays: int32(validDays),
	}

	link, err := c.Client.GenerateLink(context.Background(), req)
	if err != nil {
		logrus.Error(err)
		return ""
	}

	return link.LinkCode
}
