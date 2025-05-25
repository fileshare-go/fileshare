package model

type Link struct {
	LinkCode string `gorm:"primaryKey"`
	Sha256   string `gorm:"primaryKey"`
}
