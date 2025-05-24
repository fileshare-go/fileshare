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

	fileInfo, ok := model.GetFileInfo(request.Meta.Sha256, s.DB)
	if ok {
		return fileInfo.BuildUploadTask(), nil
	}

	fileInfo = model.NewFileInfoFromUpload(request)

	s.DB.Create(fileInfo)

	return fileInfo.BuildUploadTask(), nil
}

// upload receives chunks from client, save lockfile
func (s *UploadServer) Upload(stream pb.UploadService_UploadServer) error {
	logrus.Debug("[Upload] Starting Upload Process!")

	handler := NewHandler(stream, s.DB)

	// if recv or saving has any err, just close and return err
	if err := handler.Recv(); err != nil {
		return handler.CloseWithErr(err)
	}

	// if recv and saving do not has any error, validate and close
	handler.ValidateAndClose()

	logrus.Debug("[Upload] Ending Upload Process!")
	return nil
}
