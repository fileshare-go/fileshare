package dbmanager

import (
	"time"

	"github.com/chanmaoganda/fileshare/internal/model"
	"gorm.io/gorm"
)

// handle fileinfo, sharelink and records in database
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
	if invalidFileInfo(fileInfo) {
		return false
	}
	result := m.DB.Create(fileInfo)
	return result.Error == nil && result.RowsAffected == 1
}

func (m *DBManager) UpdateFileInfo(fileInfo *model.FileInfo) bool {
	if invalidFileInfo(fileInfo) {
		return false
	}
	result := m.DB.Save(fileInfo)
	return result.Error != nil && result.RowsAffected == 1
}

func (m *DBManager) SelectShareLink(shareLink *model.ShareLink) bool {
	result := m.DB.Where(shareLink).First(shareLink)
	return result.Error == nil && result.RowsAffected == 1
}

func (m *DBManager) SelectValidShareLink(shareLink *model.ShareLink) bool {
	result := m.DB.Scopes(NotOutdated()).Find(&shareLink)
	return result.Error == nil && result.RowsAffected == 1
}

func (m *DBManager) CreateShareLink(shareLink *model.ShareLink) bool {
	if invalidShareLink(shareLink) {
		return false
	}
	result := m.DB.Create(shareLink)
	return result.Error == nil && result.RowsAffected == 1
}

func (m *DBManager) UpdateShareLink(shareLink *model.ShareLink) bool {
	if invalidShareLink(shareLink) {
		return false
	}
	result := m.DB.Where("sha256 = ?", shareLink.Sha256).Save(shareLink)
	return result.Error == nil && result.RowsAffected == 1
}

func (m *DBManager) CreateRecord(record *model.Record) bool {
	if invalidRecord(record) {
		return false
	}
	result := m.DB.Create(record)
	return result.Error == nil && result.RowsAffected == 1
}

func NotOutdated() func(DB *gorm.DB) *gorm.DB {
	return func(DB *gorm.DB) *gorm.DB {
		return DB.Where("outdated_at > ?", time.Now())
	}
}

func invalidFileInfo(fileInfo *model.FileInfo) bool {
	return fileInfo.Filename == "" || fileInfo.Sha256 == ""
}

func invalidShareLink(sharelink *model.ShareLink) bool {
	return sharelink.Sha256 == "" || sharelink.LinkCode == ""
}

func invalidRecord(record *model.Record) bool {
	return record.Sha256 == "" || record.InteractAction == ""
}
