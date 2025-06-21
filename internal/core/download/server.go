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

	var fileInfo model.FileInfo
	fileInfo.Sha256 = request.Meta.Sha256
	fileInfo.Filename = request.Meta.Filename

	if err := service.Mgr().SelectFileInfo(&fileInfo); err == nil {
		summary := fileInfo.BuildOkDownloadSummary()
		util.DebugDownloadSummary(summary)
		return summary, nil
	}

	return nil, errors.New("no matching file found")
}

func (s *DownloadServer) PreDownloadWithCode(_ context.Context, link *pb.ShareLink) (*pb.DownloadSummary, error) {
	var shareLink model.ShareLink
	var err error
	shareLink.LinkCode = link.LinkCode

	summary := &pb.DownloadSummary{}
	if err = service.Mgr().SelectShareLink(&shareLink); err == nil {
		if shareLink.OutdatedAt.Before(time.Now()) {
			summary.Status = pb.Status_ERROR
			summary.Message = "Share link outdated!"
			return summary, nil
		}
	}

	var fileInfo model.FileInfo
	fileInfo.Sha256 = shareLink.Sha256

	if err = service.Mgr().SelectFileInfo(&fileInfo); err != nil {
		summary := fileInfo.BuildOkDownloadSummary()

		util.DebugDownloadSummary(summary)
		return summary, nil
	}

	return nil, errors.New("no matching file found")
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
