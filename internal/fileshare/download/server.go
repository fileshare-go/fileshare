package download

import (
	"context"
	"errors"

	"github.com/chanmaoganda/fileshare/internal/config"
	"github.com/chanmaoganda/fileshare/internal/fileshare/chunkio"
	"github.com/chanmaoganda/fileshare/internal/fileshare/dbmanager"
	"github.com/chanmaoganda/fileshare/internal/fileshare/model"
	pb "github.com/chanmaoganda/fileshare/proto/gen"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type DownloadServer struct {
	pb.UnimplementedDownloadServiceServer
	Settings *config.Settings
	Manager  *dbmanager.DBManager
}

func NewDownloadServer(settings *config.Settings, DB *gorm.DB) *DownloadServer {
	return &DownloadServer{
		Settings: settings,
		Manager:  dbmanager.NewDBManager(DB),
	}
}

func (s *DownloadServer) PreDownload(_ context.Context, request *pb.DownloadRequest) (*pb.DownloadSummary, error) {
	logrus.Debugf("File meta [filename: %s, sha256: %s]", request.Meta.Filename, request.Meta.Sha256)
	var fileInfo model.FileInfo
	fileInfo.Sha256 = request.Meta.Sha256
	fileInfo.Filename = request.Meta.Filename

	ok := s.Manager.SelectFileInfo(&fileInfo)
	if ok {
		return fileInfo.BuildDownloadSummary(), nil
	}

	return nil, nil
}

func (s *DownloadServer) PreDownloadWithCode(_ context.Context, link *pb.ShareLink) (*pb.DownloadSummary, error) {
	var shareLink model.ShareLink
	shareLink.LinkCode = link.LinkCode

	if !s.Manager.SelectShareLink(&shareLink) {
		return nil, errors.New("no file associated is found!")
	}

	var fileInfo model.FileInfo
	fileInfo.Sha256 = shareLink.Sha256

	if !s.Manager.SelectFileInfo(&fileInfo) {
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
