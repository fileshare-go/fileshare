package service

import (
	"github.com/chanmaoganda/fileshare/internal/model"
	"github.com/chanmaoganda/fileshare/internal/pkg/service"
)

type ShareLink struct {
	service.Service
}

func (r ShareLink) Get() ([]model.ShareLink, error) {
	var sharelinkList []model.ShareLink
	db := r.Orm.Find(&sharelinkList)
	return sharelinkList, db.Error
}

func (r ShareLink) Insert(sharelink model.ShareLink) error {
	db := r.Orm.Create(&sharelink)
	if db.RowsAffected != 1 {
		return notAffectedError
	}
	return db.Error
}
