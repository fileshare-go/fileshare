package db

import (
	"time"

	"gorm.io/gorm"
)

func NotOutdated() func(DB *gorm.DB) *gorm.DB {
	return func(DB *gorm.DB) *gorm.DB {
		return DB.Where("outdated_at > ?", time.Now())
	}
}
