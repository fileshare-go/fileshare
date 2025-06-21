package db

import (
	"github.com/chanmaoganda/fileshare/internal/model"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func OpenClientDB(sqliteFile string) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(sqliteFile), &gorm.Config{})
	if err != nil {
		logrus.Fatalf("sqlite %s: %v", sqliteFile, err)
	}
	if err := db.AutoMigrate(&model.FileInfo{}); err != nil {
		logrus.Panic(err)
	}
	return db
}

func OpenServerDB(sqliteFile string) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(sqliteFile), &gorm.Config{})
	if err != nil {
		logrus.Fatalf("sqlite %s: %v", sqliteFile, err)
	}
	if err := db.AutoMigrate(&model.FileInfo{}, &model.ShareLink{}, &model.Record{}); err != nil {
		logrus.Panic(err)
	}
	return db
}
