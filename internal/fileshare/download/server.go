package download

import (
	"context"

	"github.com/chanmaoganda/fileshare/internal/config"
	"github.com/chanmaoganda/fileshare/internal/lockfile"
	pb "github.com/chanmaoganda/fileshare/proto/gen"
	"github.com/sirupsen/logrus"
)

type DownloadServer struct {
	pb.UnimplementedDownloadServiceServer
	Settings *config.Settings
}

func (s *DownloadServer) PreDownload(_ context.Context, meta *pb.FileMeta) (*pb.DownloadSummary, error) {
	logrus.Debugf("File meta [filename: %s, sha256: %s]", meta.Filename, meta.Sha256)

	summary := pb.DownloadSummary{
		Meta: &pb.FileMeta{
			Filename: meta.Filename,
			Sha256:   meta.Sha256,
		},
	}

	lockfile, err := lockfile.ReadLockFile(meta.Sha256)
	if err != nil {
		return &summary, nil
	}

	summary.ChunkList = lockfile.ChunkList
	summary.ChunkSize = lockfile.ChunkSize
	summary.ChunkNumber = lockfile.TotalChunkNumber
	// summary.FileSize = lockfile.FileSize
	return &summary, nil
}
