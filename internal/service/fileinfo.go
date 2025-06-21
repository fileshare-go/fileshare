package service

import (
	"github.com/chanmaoganda/fileshare/internal/model"
	"github.com/chanmaoganda/fileshare/internal/pkg/service"
)

type FileInfo struct {
	service.Service
}

func (f *FileInfo) GetFileInfo() ([]model.FileInfo, error) {
	var fileinfoList []model.FileInfo
	db := f.Orm.Find(&fileinfoList)
	return fileinfoList, db.Error
}

func (f *FileInfo) SelectFileInfo(fileinfo *model.FileInfo) error {
	db := f.Orm.Where(fileinfo).Find(fileinfo)
	return oneRowAffected(db)
}

func (f *FileInfo) InsertFileInfo(fileinfo *model.FileInfo) error {
	db := f.Orm.Create(&fileinfo)
	return oneRowAffected(db)
}

func (f *FileInfo) UpdateFileInfo(fileinfo *model.FileInfo) error {
	db := f.Orm.Save(fileinfo)
	return oneRowAffected(db)
}
