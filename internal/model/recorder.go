package model

import "time"

type Record struct {
	Sha256         string    `gorm:"primaryKey"`
	InteractAction string    `gorm:"primaryKey,size:8" comment:"this field records which action is done, either upload, linkgen, or download"`
	ClientIp       string    `gorm:"primaryKey,size:64"`
	Os             string    `gorm:"primaryKey,size:64"`
	Time           time.Time `gorm:"primaryKey"`
}
