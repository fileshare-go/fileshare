package upload

import (
	"context"

	"github.com/chanmaoganda/fileshare/internal/config"
	"github.com/chanmaoganda/fileshare/internal/fileshare/chunkstream/recv"
	"github.com/chanmaoganda/fileshare/internal/fileshare/dbmanager"
	"github.com/chanmaoganda/fileshare/internal/model"
	"github.com/chanmaoganda/fileshare/internal/pkg/debugprint"
	pb "github.com/chanmaoganda/fileshare/proto/gen"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UploadServer struct {
	pb.UnimplementedUploadServiceServer
	Settings *config.Settings
	Manager  *dbmanager.DBManager
}

func NewUploadServer(settings *config.Settings, DB *gorm.DB) *UploadServer {
	return &UploadServer{
		Settings: settings,
		Manager:  dbmanager.NewDBManager(DB),
	}
}

// pre upload receives a task from client, calculate missing chunks and send the task back
func (s *UploadServer) PreUpload(ctx context.Context, request *pb.UploadRequest) (*pb.UploadTask, error) {
	logrus.Debugf("PreUpload request [filename: %s, file size: %d, sha256: %s]", debugprint.Render(request.Meta.Filename), request.FileSize, debugprint.Render(request.Meta.Sha256[:8]))

	fileInfo := &model.FileInfo{
		Sha256: request.Meta.Sha256,
	}
	// , ok := model.GetFileInfo(request.Meta.Sha256, s.DB)
	if s.Manager.SelectFileInfo(fileInfo) {
		logrus.Debug("Existing file info ", fileInfo.Filename)
		return fileInfo.BuildUploadTask(), nil
	}

	fileInfo = model.NewFileInfoFromUpload(request)

	logrus.Debug("Creating file info ", fileInfo.Filename)
	s.Manager.CreateFileInfo(fileInfo)

	return fileInfo.BuildUploadTask(), nil
}

func (s *UploadServer) Upload(stream pb.UploadService_UploadServer) error {
	logrus.Debug("[Upload] Starting Upload Process!")

	chunkStream := recv.NewServerRecvStream(s.Settings, s.Manager, stream)
	if err := chunkStream.RecvStreamChunks(); err != nil {
		return chunkStream.CloseStream(false)
	}

	validate := chunkStream.Validate()
	return chunkStream.CloseStream(validate)
}
