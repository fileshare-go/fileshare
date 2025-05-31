package sharelink

import (
	"context"
	"strconv"

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

func (c *ShareLinkClient) GenerateLink(args []string) string {
	req := ParseArgs2Request(args)

	link, err := c.Client.GenerateLink(context.Background(), req)
	if err != nil {
		logrus.Error(err)
		return ""
	}

	return link.LinkCode
}

func ParseArgs2Request(args []string) *pb.ShareLinkRequest {
	var validDays int
	var err error
	if len(args) < 3 {
		validDays = 0
	} else {
		validDays, err = strconv.Atoi(args[2])
		if err != nil {
			validDays = 0
		}
	}

	return &pb.ShareLinkRequest{
		Meta: &pb.FileMeta{
			Filename: args[0],
			Sha256:   args[1],
		},
		ValidDays: int32(validDays),
	}
}
