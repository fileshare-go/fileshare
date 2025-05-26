package dbmanager

import (
	"github.com/chanmaoganda/fileshare/internal/model"
	"gorm.io/gorm"
)

type DBManager struct {
	DB *gorm.DB
}

func NewDBManager(DB *gorm.DB) *DBManager {
	return &DBManager{
		DB: DB,
	}
}

func (m *DBManager) SelectFileInfo(fileInfo *model.FileInfo) bool {
	result := m.DB.Where(fileInfo).First(fileInfo)
	return result.Error == nil && result.RowsAffected == 1
}

func (m *DBManager) CreateFileInfo(fileInfo *model.FileInfo) bool {
	result := m.DB.Create(fileInfo)
	return result.Error == nil && result.RowsAffected == 1
}

func (m *DBManager) UpdateFileInfo(fileInfo *model.FileInfo) bool {
	result := m.DB.Save(fileInfo)
	return result.Error != nil && result.RowsAffected == 1
}

func (m *DBManager) SelectShareLink(shareLink *model.ShareLink) bool {
	result := m.DB.Where(shareLink).First(shareLink)
	return result.Error == nil && result.RowsAffected == 1
}

func (m *DBManager) CreateShareLink(shareLink *model.ShareLink) bool {
	result := m.DB.Create(shareLink)
	return result.Error == nil && result.RowsAffected == 1
}

func (m *DBManager) UpdateShareLink(shareLink *model.ShareLink) bool {
	result := m.DB.Save(shareLink).First(shareLink)
	return result.Error == nil && result.RowsAffected == 1
}
