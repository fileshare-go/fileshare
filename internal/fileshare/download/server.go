package download

import (
	"context"
	"errors"

	"github.com/chanmaoganda/fileshare/internal/config"
	"github.com/chanmaoganda/fileshare/internal/fileshare/chunkio"
	"github.com/chanmaoganda/fileshare/internal/fileshare/model"
	pb "github.com/chanmaoganda/fileshare/proto/gen"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type DownloadServer struct {
	pb.UnimplementedDownloadServiceServer
	Settings *config.Settings
	DB       *gorm.DB
}

func NewDownloadServer(settings *config.Settings, DB *gorm.DB) *DownloadServer {
	return &DownloadServer{
		Settings: settings,
		DB: DB,
	}
}

func (s *DownloadServer) PreDownload(_ context.Context, request *pb.DownloadRequest) (*pb.DownloadSummary, error) {
	logrus.Debugf("File meta [filename: %s, sha256: %s]", request.Meta.Filename, request.Meta.Sha256)
	fileInfo, ok := model.GetFileInfo(request.Meta.Sha256, s.DB)
	if ok {
		return fileInfo.BuildDownloadSummary(), nil
	}

	return nil, nil
}

func (s *DownloadServer) PreDownloadWithCode(_ context.Context, link *pb.ShareLink) (*pb.DownloadSummary, error) {
	var fileLink model.Link
	if s.DB.Where("link_code = ?", link.LinkCode).First(&fileLink).RowsAffected == 0 {
		return nil, errors.New("no file associated is found!")
	}

	fileInfo, ok := model.GetFileInfo(fileLink.Sha256, s.DB)
	if ok {
		return fileInfo.BuildDownloadSummary(), nil
	}

	return nil, nil
}

func (s *DownloadServer) Download(task *pb.DownloadTask, stream pb.DownloadService_DownloadServer) error {
	logrus.Debugf("Download Task: %s", task.Meta.Sha256)

	// if chunklist is empty, at least send one chunk
	if len(task.ChunkList) == 0 {
		task.ChunkList = append(task.ChunkList, 0)
	}

	for _, chunkIndex := range task.ChunkList {
		bytes := chunkio.UploadChunk(task.Meta.Sha256, chunkIndex)

		logrus.Debugf("File Chunk:[filename: %s, sha256: %s, chunk index: %d]", task.Meta.Filename, task.Meta.Sha256, chunkIndex)

		chunk := &pb.FileChunk{
			Sha256:     task.Meta.Sha256,
			ChunkIndex: chunkIndex,
			Data:       bytes,
		}

		if err := stream.Send(chunk); err != nil {
			logrus.Error(err)
			break
		}
	}

	logrus.Debugf("File Sent! %s", task.Meta.Filename)
	return nil
}
