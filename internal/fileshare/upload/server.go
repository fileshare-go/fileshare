package upload

import (
	"context"

	"github.com/chanmaoganda/fileshare/internal/config"
	"github.com/chanmaoganda/fileshare/internal/fileshare/model"
	pb "github.com/chanmaoganda/fileshare/proto/gen"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UploadServer struct {
	pb.UnimplementedUploadServiceServer
	Settings *config.Settings
	DB       *gorm.DB
}

// pre upload receives a task from client, calculate missing chunks and send the task back
func (s *UploadServer) PreUpload(_ context.Context, request *pb.UploadRequest) (*pb.UploadTask, error) {
	logrus.Debugf("Upload request [filename: %s, file size: %d, sha256: %s]", request.Meta.Filename, request.FileSize, request.Meta.Sha256)

	fileInfo, ok := s.getFileInfo(request.Meta.Sha256)
	if ok {
		return fileInfo.BuildUploadTask(), nil
	}

	fileInfo = model.NewFileInfo(request)

	s.DB.Create(fileInfo)

	return fileInfo.BuildUploadTask(), nil
}

// upload receives chunks from client, save lockfile
func (s *UploadServer) Upload(stream pb.UploadService_UploadServer) error {
	logrus.Debug("Starting Upload Process!")

	handler := NewHandler(stream)

	// if recv or saving has any err, just close and return err
	if err := handler.Recv(); err != nil {
		return handler.CloseWithErr(err)
	}

	// if recv and saving do not has any error, validate and close
	handler.ValidateAndClose()

	logrus.Debug("Ending Upload Process!")
	return nil
}

func (s *UploadServer) getFileInfo(sha256 string) (*model.FileInfo, bool) {
	var fileInfo model.FileInfo

	if s.DB.First(&fileInfo, sha256).RowsAffected != 0 {
		return &fileInfo, true
	}

	return &fileInfo, false
}
