package sharelink

import (
	"context"
	"errors"

	"github.com/chanmaoganda/fileshare/internal/config"
	"github.com/chanmaoganda/fileshare/internal/fileshare/model"
	pb "github.com/chanmaoganda/fileshare/proto/gen"
	"gorm.io/gorm"
)

type ShareLinkServer struct {
	pb.UnimplementedShareLinkServiceServer
	Settings *config.Settings
	DB *gorm.DB
	RangGen *RandStringGen
}

func NewShareLinkServer(settings *config.Settings, DB *gorm.DB) *ShareLinkServer {
	return &ShareLinkServer{
		Settings: settings,
		RangGen: NewRandStringGen(),
		DB: DB,
	}
}

func (s *ShareLinkServer) GenerateLink(_ context.Context, meta *pb.FileMeta) (*pb.ShareLink, error) {
	var link model.Link
	if s.DB.Where("sha256 = ?", meta.Sha256).First(&link).RowsAffected != 0 {
		return &pb.ShareLink{
			LinkCode: link.LinkCode,
		}, nil
	}

	linkCode := s.RangGen.generateCode(s.Settings.ShareCodeLength)
	var fileInfo model.FileInfo

	if s.DB.Where("sha256 = ?", meta.Sha256).First(&fileInfo).RowsAffected == 0 {
		return nil, errors.New("file meta not exist")
	}

	link.LinkCode = linkCode
	link.Sha256 = fileInfo.Sha256

	if s.DB.Create(&link).RowsAffected == 0 {
		return nil, errors.New("link code gen error")
	}

	return &pb.ShareLink{
		LinkCode: linkCode,
	}, nil
}
