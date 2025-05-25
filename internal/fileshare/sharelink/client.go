package sharelink

import (
	"context"

	pb "github.com/chanmaoganda/fileshare/proto/gen"
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

func (c *ShareLinkClient) GenerateLink(filename, sha256 string) string {
	meta := &pb.FileMeta{
		Filename: filename,
		Sha256: sha256,
	}

	link, err := c.Client.GenerateLink(context.Background(), meta)
	if err != nil {
		logrus.Error(err)
		return ""
	}

	return link.LinkCode
}
