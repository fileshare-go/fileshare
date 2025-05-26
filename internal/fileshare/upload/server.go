package upload

import (
	"context"

	"github.com/chanmaoganda/fileshare/internal/config"
	"github.com/chanmaoganda/fileshare/internal/debugprint"
	"github.com/chanmaoganda/fileshare/internal/fileshare/dbmanager"
	"github.com/chanmaoganda/fileshare/internal/model"
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
		return fileInfo.BuildUploadTask(), nil
	}

	fileInfo = model.NewFileInfoFromUpload(request)

	s.Manager.CreateFileInfo(fileInfo)

	return fileInfo.BuildUploadTask(), nil
}

// upload receives chunks from client, save lockfile
func (s *UploadServer) Upload(stream pb.UploadService_UploadServer) error {
	logrus.Debug("[Upload] Starting Upload Process!")

	handler := NewHandler(s.Settings, stream, s.Manager.DB)

	// if recv or saving has any err, just close and return err
	if err := handler.Recv(); err != nil {
		return handler.CloseWithErr(err)
	}

	// if recv and saving do not has any error, validate and close
	handler.ValidateAndClose()

	logrus.Debug("[Upload] Ending Upload Process!")
	return nil
}
