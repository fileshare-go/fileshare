package model

type File struct {
	Filename    string `gorm:"primaryKey"`
	Sha256      string `gorm:"primaryKey"`
	ChunkSize   int64
	ChunkNumber int32
	FileSize    int64
	UploadedChunks string
}
