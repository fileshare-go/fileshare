package download

import (
	"context"

	"github.com/chanmaoganda/fileshare/internal/config"
	pb "github.com/chanmaoganda/fileshare/proto/gen"
	"github.com/sirupsen/logrus"
)

type DownloadServer struct {
	pb.UnimplementedDownloadServiceServer
	Settings *config.Settings
}

func (s *DownloadServer) PreDownload(_ context.Context, meta *pb.FileMeta) (*pb.DownloadSummary, error) {
	logrus.Debugf("File meta [filename: %s, sha256: %s]", meta.Filename, meta.Sha256)

	return nil, nil
}
