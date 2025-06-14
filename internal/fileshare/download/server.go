package download

import (
	"context"
	"errors"
	"time"

	"github.com/chanmaoganda/fileshare/internal/config"
	"github.com/chanmaoganda/fileshare/internal/fileshare/chunkstream/send"
	"github.com/chanmaoganda/fileshare/internal/model"
	"github.com/chanmaoganda/fileshare/internal/pkg/dbmanager"
	"github.com/chanmaoganda/fileshare/internal/pkg/debugprint"
	pb "github.com/chanmaoganda/fileshare/internal/proto/gen"
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
	debugprint.DebugMeta(request.Meta)

	var fileInfo model.FileInfo
	fileInfo.Sha256 = request.Meta.Sha256
	fileInfo.Filename = request.Meta.Filename

	if s.Manager.SelectFileInfo(&fileInfo) {
		summary := fileInfo.BuildDownloadSummary()
		debugprint.DebugDownloadSummary(summary)
		return summary, nil
	}

	return nil, errors.New("no matching file found")
}

func (s *DownloadServer) PreDownloadWithCode(_ context.Context, link *pb.ShareLink) (*pb.DownloadSummary, error) {
	var shareLink model.ShareLink
	shareLink.LinkCode = link.LinkCode

	summary := &pb.DownloadSummary{}
	if s.Manager.SelectShareLink(&shareLink) {
		if shareLink.OutdatedAt.Before(time.Now()) {
			summary.Status = pb.Status_ERROR
			summary.Message = "Share link outdated!"
			return summary, nil
		}
	}

	var fileInfo model.FileInfo
	fileInfo.Sha256 = shareLink.Sha256

	if s.Manager.SelectFileInfo(&fileInfo) {
		summary := fileInfo.BuildDownloadSummary()
		summary.Status = pb.Status_OK
		summary.Message = "Share link found!"

		debugprint.DebugDownloadSummary(summary)
		return summary, nil
	}

	return nil, errors.New("no matching file found")
}

func (s *DownloadServer) Download(task *pb.DownloadTask, stream pb.DownloadService_DownloadServer) error {
	debugprint.DebugDownloadTask(task)
	sendStream := send.NewServerSendStream(s.Settings, s.Manager, task, stream)

	if err := sendStream.SendStreamChunks(); err != nil {
		return err
	}

	logrus.Debugf("File Sent! %s", task.Meta.Filename)
	return sendStream.CloseStream()
}
