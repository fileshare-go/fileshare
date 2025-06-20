package service

import (
	"github.com/chanmaoganda/fileshare/internal/model"
	"github.com/chanmaoganda/fileshare/internal/pkg/service"
)

type FileInfo struct {
	service.Service
}

func (f *FileInfo) Get() ([]model.FileInfo, error) {
	var fileinfoList []model.FileInfo
	db := f.Orm.Find(&fileinfoList)
	return fileinfoList, db.Error
}

func (f *FileInfo) Insert(fileinfo model.FileInfo) error {
	db := f.Orm.Create(&fileinfo)
	if db.RowsAffected != 1 {
		return notAffectedError
	}
	return db.Error
}
