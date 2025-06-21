package service

import (
	"github.com/chanmaoganda/fileshare/internal/model"
	"github.com/chanmaoganda/fileshare/internal/pkg/db"
	"github.com/chanmaoganda/fileshare/internal/pkg/service"
)

type ShareLink struct {
	service.Service
}

func (s ShareLink) GetShareLink() ([]model.ShareLink, error) {
	var sharelinkList []model.ShareLink
	db := s.Orm.Find(&sharelinkList)
	return sharelinkList, db.Error
}

func (s *ShareLink) SelectValidShareLink(sharelink *model.ShareLink) error {
	db := s.Orm.Scopes(db.NotOutdated()).Where(sharelink).Find(sharelink)
	return oneRowAffected(db)
}

func (s *ShareLink) SelectShareLink(sharelink *model.ShareLink) error {
	db := s.Orm.Where(sharelink).Find(sharelink)
	return oneRowAffected(db)
}

func (s *ShareLink) InsertShareLink(sharelink *model.ShareLink) error {
	db := s.Orm.Create(&sharelink)
	return oneRowAffected(db)
}

func (s *ShareLink) UpdateShareLink(sharelink *model.ShareLink) error {
	db := s.Orm.Save(sharelink)
	return oneRowAffected(db)
}
