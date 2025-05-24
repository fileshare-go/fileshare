package upload

import (
	"context"

	"github.com/chanmaoganda/fileshare/internal/config"
	"github.com/chanmaoganda/fileshare/internal/fileshare/chunker"
	"github.com/chanmaoganda/fileshare/internal/fileshare/model"
	"github.com/chanmaoganda/fileshare/internal/fileutil"
	"github.com/chanmaoganda/fileshare/internal/lockfile"
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

	chunkSummary := chunker.DealChunkSize(request.FileSize)
	logrus.Debugf("Chunk Summary [chunk number: %d, chunk size: %d]", chunkSummary.Number, chunkSummary.Size)

	chunkList := make([]int32, 0)
	for index := range chunkSummary.Number {
		chunkList = append(chunkList, index)
	}

	required := getMissingChunks(request.Meta.Sha256, chunkSummary.Number)

	s.storeFileInfo(request, chunkSummary)

	return &pb.UploadTask{
		Meta:        request.Meta,
		ChunkNumber: chunkSummary.Number,
		ChunkSize:   chunkSummary.Size,
		ChunkList:   required,
	}, nil
}

// upload receives chunks from client, save lockfile
func (s *UploadServer) Upload(stream pb.UploadService_UploadServer) error {
	logrus.Debug("Starting Upload Process!")

	handler := NewHandler(stream)

	// if recv or saving has any err, just close and return err
	if err := handler.RecvAndSaveLock(); err != nil {
		return handler.CloseWithErr(err)
	}

	// if recv and saving do not has any error, validate and close
	handler.ValidateAndClose()

	logrus.Debug("Ending Upload Process!")
	return nil
}

func (s *UploadServer) storeFileInfo(request *pb.UploadRequest, summary chunker.ChunkSummary) {


	file := model.File{
		Filename:    request.Meta.Filename,
		Sha256:      request.Meta.Sha256,
		ChunkSize:   summary.Size,
		ChunkNumber: summary.Number,
		FileSize:    request.FileSize,
	}

	s.DB.Create(&file)
}

// get missing chunks from lockfile if exists, either return the enum of `total`
func getMissingChunks(lockFolder string, total int32) []int32 {
	lockPath := lockfile.GetLockPath(lockFolder)
	if fileutil.FileExists(lockPath) {
		lock, err := lockfile.ReadLockFile(lockFolder)
		if err != nil {
			logrus.Error(err)
			return []int32{}
		}

		return lock.RemainingChunks()
	}
	result := []int32{}
	for i := range total {
		result = append(result, i)
	}
	return result
}

