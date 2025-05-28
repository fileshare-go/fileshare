package model

import "time"

type ShareLink struct {
	Sha256     string `gorm:"primaryKey"`
	LinkCode   string
	CreatedBy  string
	CreatedAt  time.Time
	OutdatedAt time.Time
}
