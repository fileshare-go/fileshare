package sharelink

import (
	"context"
	"errors"

	"github.com/chanmaoganda/fileshare/internal/config"
	"github.com/chanmaoganda/fileshare/internal/fileshare/dbmanager"
	"github.com/chanmaoganda/fileshare/internal/model"
	"github.com/chanmaoganda/fileshare/internal/pkg/debugprint"
	pb "github.com/chanmaoganda/fileshare/internal/proto/gen"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ShareLinkServer struct {
	pb.UnimplementedShareLinkServiceServer
	Settings *config.Settings
	Manager  *dbmanager.DBManager
	RangGen  *RandStringGen
}

func NewShareLinkServer(settings *config.Settings, DB *gorm.DB) *ShareLinkServer {
	return &ShareLinkServer{
		Settings: settings,
		RangGen:  NewRandStringGen(),
		Manager:  dbmanager.NewDBManager(DB),
	}
}

func (s *ShareLinkServer) GenerateLink(_ context.Context, meta *pb.FileMeta) (*pb.ShareLink, error) {
	logrus.Debugf("Generating sharelink for %s", debugprint.Render(meta.Sha256[:8]))

	link := &model.ShareLink{
		Sha256: meta.Sha256,
	}
	if s.Manager.SelectShareLink(link) {
		logrus.Debugf("Existing sharelink for %s is %s", debugprint.Render(meta.Sha256[:8]), debugprint.Render(link.LinkCode))

		return &pb.ShareLink{
			LinkCode: link.LinkCode,
		}, nil
	}

	// if db doesn't have records then construct this ShareLink
	linkCode := s.RangGen.generateCode(s.Settings.ShareCodeLength)
	logrus.Debugf("Generated code is %s", debugprint.Render(linkCode))

	fileInfo := &model.FileInfo{
		Sha256: meta.Sha256,
	}

	if !s.Manager.SelectFileInfo(fileInfo) {
		return nil, errors.New("file meta not exist")
	}

	link.LinkCode = linkCode
	link.Sha256 = fileInfo.Sha256

	if !s.Manager.CreateShareLink(link) {
		return nil, errors.New("link code gen error")
	}

	logrus.Debugf("Generated sharelink for %s is %s", debugprint.Render(meta.Sha256[:8]), debugprint.Render(link.LinkCode))

	return &pb.ShareLink{
		LinkCode: linkCode,
	}, nil
}
