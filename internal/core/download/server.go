package download

import (
	"context"
	"errors"
	"time"

	"github.com/chanmaoganda/fileshare/internal/core/chunkstream/send"
	"github.com/chanmaoganda/fileshare/internal/model"
	"github.com/chanmaoganda/fileshare/internal/pkg/util"
	pb "github.com/chanmaoganda/fileshare/internal/proto/gen"
	"github.com/chanmaoganda/fileshare/internal/service"
	"github.com/sirupsen/logrus"
)

type DownloadServer struct {
	pb.UnimplementedDownloadServiceServer
}

func NewDownloadServer() *DownloadServer {
	return &DownloadServer{}
}

func (s *DownloadServer) PreDownload(_ context.Context, request *pb.DownloadRequest) (*pb.DownloadSummary, error) {
	util.DebugMeta(request.Meta)

	fileInfo := &model.FileInfo{
		Sha256: request.Meta.Sha256,
	}

	if service.Orm().Find(fileInfo).RowsAffected == 0 {
		return nil, errors.New("no matching file found")
	}

	summary := buildOkDownloadSummary(fileInfo)
	// util.DebugDownloadSummary(summary)
	return summary, nil
}

func (s *DownloadServer) PreDownloadWithCode(_ context.Context, link *pb.ShareLink) (*pb.DownloadSummary, error) {
	shareLink := &model.ShareLink{
		LinkCode: link.LinkCode,
	}

	if service.Orm().Find(shareLink).RowsAffected == 1 {
		// if link is outdated, just return with error
		if shareLink.OutdatedAt.Before(time.Now()) {
			return &pb.DownloadSummary{
				Status:  pb.Status_ERROR,
				Message: "Share link outdated!",
			}, nil
		}
	}

	// if not outdated, query for existing file info
	fileInfo := &model.FileInfo{
		Sha256: shareLink.Sha256,
	}

	if service.Orm().Find(fileInfo).RowsAffected == 0 {
		return nil, errors.New("no matching file found")
	}

	return buildOkDownloadSummary(fileInfo), nil
}

func (s *DownloadServer) Download(task *pb.DownloadTask, stream pb.DownloadService_DownloadServer) error {
	util.DebugDownloadTask(task)
	sendStream := send.NewServerSendStream(task, stream)

	if err := sendStream.SendStreamChunks(); err != nil {
		return err
	}

	logrus.Debugf("File Sent! %s", task.Meta.Filename)
	return sendStream.CloseStream()
}

func buildOkDownloadSummary(f *model.FileInfo) *pb.DownloadSummary {
	return &pb.DownloadSummary{
		Meta: &pb.FileMeta{
			Filename: f.Filename,
			Sha256:   f.Sha256,
			FileSize: f.FileSize,
		},
		FileSize:    f.FileSize,
		ChunkNumber: f.ChunkNumber,
		ChunkSize:   f.ChunkSize,
		ChunkList:   f.GetUploadedChunks(),
		Status:      pb.Status_OK,
		Message:     "Share link found!",
	}
}
