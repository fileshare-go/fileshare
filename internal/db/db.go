package db

import (
	"github.com/chanmaoganda/fileshare/internal/fileshare/model"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func SetupDB(sqliteFile string) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(sqliteFile), &gorm.Config{})
	if err != nil {
		logrus.Fatal("sqlite: ", err)
	}
	db.AutoMigrate(&model.FileInfo{})
	return db
}
