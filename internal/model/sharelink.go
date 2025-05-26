package model

type ShareLink struct {
	LinkCode string
	Sha256   string `gorm:"primaryKey"`
}
