package db

import (
	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func SetupDB(sqliteFile string) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(sqliteFile), &gorm.Config{})
	if err != nil {
		logrus.Fatal("sqlite: ", err)
	}
	return db
}
