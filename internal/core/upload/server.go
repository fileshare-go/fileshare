package upload

import (
	"context"

	"github.com/chanmaoganda/fileshare/internal/core/chunkstream/recv"
	"github.com/chanmaoganda/fileshare/internal/model"
	"github.com/chanmaoganda/fileshare/internal/pkg/util"
	pb "github.com/chanmaoganda/fileshare/internal/proto/gen"
	"github.com/chanmaoganda/fileshare/internal/service"
	"github.com/sirupsen/logrus"
)

type UploadServer struct {
	pb.UnimplementedUploadServiceServer
}

func NewUploadServer() *UploadServer {
	return &UploadServer{}
}

// pre upload receives a task from client, calculate missing chunks and send the task back
func (s *UploadServer) PreUpload(ctx context.Context, request *pb.UploadRequest) (*pb.UploadTask, error) {
	var err error
	logrus.Debugf("PreUpload request [filename: %s, file size: %d, sha256: %s]", util.Render(request.Meta.Filename), request.FileSize, util.Render(request.Meta.Sha256[:8]))

	fileInfo := &model.FileInfo{
		Sha256: request.Meta.Sha256,
	}

	db := service.Orm().Find(fileInfo)

	if db.RowsAffected == 1 {
		// if fileinfo exists, then use stored info
		logrus.Debug("Existing file info ", fileInfo.Filename)
		return assembleUploadTask(fileInfo), nil
	}

	fileInfo = assembleFileInfo(request)

	logrus.Debug("Creating file info ", fileInfo.Filename)
	if err = service.Orm().Save(fileInfo).Error; err != nil {
		return nil, err
	}

	return assembleUploadTask(fileInfo), nil
}

func (s *UploadServer) Upload(stream pb.UploadService_UploadServer) error {
	logrus.Debug("[Upload] Starting Upload Process!")

	recvStream := recv.NewServerRecvStream(stream)
	if err := recvStream.RecvStreamChunks(); err != nil {
		return recvStream.CloseStream(false)
	}

	validate := recvStream.ValidateRecvChunks()
	return recvStream.CloseStream(validate)
}

func assembleFileInfo(req *pb.UploadRequest) *model.FileInfo {
	fileInfo := model.FileInfo{}

	chunkSummary := dealChunkSize(req.FileSize)

	// avoid filename injection
	fileInfo.Filename = util.GetFileName(req.Meta.Filename)
	fileInfo.Sha256 = req.Meta.Sha256
	fileInfo.FileSize = req.FileSize
	fileInfo.ChunkNumber = chunkSummary.Number
	fileInfo.ChunkSize = chunkSummary.Size
	fileInfo.UploadedChunks = "[]"

	return &fileInfo
}

func assembleUploadTask(f *model.FileInfo) *pb.UploadTask {
	return &pb.UploadTask{
		Meta: &pb.FileMeta{
			Filename: f.Filename,
			Sha256:   f.Sha256,
			FileSize: f.FileSize,
		},
		ChunkNumber: f.ChunkNumber,
		ChunkSize:   f.ChunkSize,
		ChunkList:   f.GetMissingChunks(),
	}
}
