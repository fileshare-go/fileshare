package model

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/chanmaoganda/fileshare/internal/algorithms"
	"github.com/chanmaoganda/fileshare/internal/fileshare/debugprint"
	"github.com/chanmaoganda/fileshare/internal/sha256"
	pb "github.com/chanmaoganda/fileshare/proto/gen"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type FileInfo struct {
	Filename       string `gorm:"primaryKey;size:64"`
	Sha256         string `gorm:"primaryKey;size:128"`
	ChunkSize      int64
	ChunkNumber    int32
	FileSize       int64
	UploadedChunks string
	Link           ShareLink `gorm:"foreignKey:Sha256;references:Sha256"`
}

func (f *FileInfo) GetUploadedChunks() []int32 {
	var chunks []int32
	if err := json.Unmarshal([]byte(f.UploadedChunks), &chunks); err != nil {
		logrus.Error(err)
		return f.GetAllChunks()
	}

	return chunks
}

func (f *FileInfo) GetMissingChunks() []int32 {
	loaded := f.GetUploadedChunks()
	all := f.GetAllChunks()

	return algorithms.MissingElementsInSortedList(all, loaded)
}

func (f *FileInfo) GetAllChunks() []int32 {
	result := []int32{}
	for i := range f.ChunkNumber {
		result = append(result, i)
	}
	return result
}

func (f *FileInfo) UpdateChunks(newChunks []int32) {
	loaded := f.GetUploadedChunks()
	merged := algorithms.MergeList(loaded, newChunks)

	bytes, err := json.Marshal(merged)
	if err != nil {
		logrus.Error(err)
		return
	}

	f.UploadedChunks = string(bytes)
}

func (f *FileInfo) BuildUploadTask() *pb.UploadTask {
	return &pb.UploadTask{
		Meta: &pb.FileMeta{
			Filename: f.Filename,
			Sha256:   f.Sha256,
			FileSize: f.FileSize,
		},
		ChunkNumber: f.ChunkNumber,
		ChunkSize:   f.ChunkSize,
		ChunkList:   f.GetMissingChunks(),
	}
}

func (f *FileInfo) BuildDownloadTask() *pb.DownloadTask {
	return &pb.DownloadTask{
		Meta: &pb.FileMeta{
			Filename: f.Filename,
			Sha256:   f.Sha256,
			FileSize: f.FileSize,
		},
		ChunkNumber: f.ChunkNumber,
		ChunkList:   f.GetMissingChunks(),
	}
}

func (f *FileInfo) BuildDownloadSummary() *pb.DownloadSummary {
	return &pb.DownloadSummary{
		Meta: &pb.FileMeta{
			Filename: f.Filename,
			Sha256:   f.Sha256,
			FileSize: f.FileSize,
		},
		FileSize:    f.FileSize,
		ChunkNumber: f.ChunkNumber,
		ChunkSize:   f.ChunkSize,
		ChunkList:   f.GetUploadedChunks(),
	}
}

func (f *FileInfo) ValidateChunks() bool {
	filePath := fmt.Sprintf("%s/%s", f.Sha256, f.Filename)
	logrus.Debug("[Validate] File: ", debugprint.Render(filePath))
	out, err := os.Create(filePath)
	if err != nil {
		logrus.Error("[Validate]", err)
		return false
	}

	for _, index := range f.GetUploadedChunks() {
		in, err := os.Open(fmt.Sprintf("%s/%d", f.Sha256, index))
		if err != nil {
			logrus.Error("[Validate]", err)
			return false
		}

		_, err = io.Copy(out, in)
		if err != nil {
			logrus.Error("[Validate]", err)
			return false
		}
		in.Close()
	}
	out.Close()

	checkSum, err := sha256.CalculateSHA256(filePath)
	if err != nil {
		logrus.Error("[Validate]", err)
		return false
	}

	return checkSum == f.Sha256
}

func NewFileInfoFromUpload(req *pb.UploadRequest) *FileInfo {
	fileInfo := FileInfo{}

	chunkSummary := dealChunkSize(req.FileSize)

	fileInfo.Filename = req.Meta.Filename
	fileInfo.Sha256 = req.Meta.Sha256
	fileInfo.FileSize = req.FileSize
	fileInfo.ChunkNumber = chunkSummary.Number
	fileInfo.ChunkSize = chunkSummary.Size
	fileInfo.UploadedChunks = "[]"

	return &fileInfo
}

func NewFileInfoFromDownload(summary *pb.DownloadSummary) *FileInfo {
	fileInfo := FileInfo{}

	fileInfo.Filename = summary.Meta.Filename
	fileInfo.Sha256 = summary.Meta.Sha256
	fileInfo.FileSize = summary.FileSize
	fileInfo.ChunkNumber = summary.ChunkNumber
	fileInfo.ChunkSize = summary.ChunkSize
	fileInfo.UploadedChunks = "[]"

	return &fileInfo
}

func GetFileInfo(sha256 string, DB *gorm.DB) (*FileInfo, bool) {
	var fileInfo FileInfo

	if DB.First(&fileInfo, "sha256 = ?", sha256).RowsAffected != 0 {
		return &fileInfo, true
	}

	return &fileInfo, false
}

const (
	SMALL  = 1024 * 1024 // 1MB
	MEDIUM = 2 * SMALL   // 2MB
	LARGE  = 4 * SMALL   // 4MB
)

type ChunkSummary struct {
	Size   int64
	Number int32
}

func dealChunkSize(fileSize int64) ChunkSummary {
	var chunkSize int
	if fileSize < 64 * SMALL {
		chunkSize = SMALL
	} else if fileSize < 1024 * SMALL {
		chunkSize = MEDIUM
	} else {
		chunkSize = LARGE
	}

	chunkNumber := fileSize / int64(chunkSize)

	if fileSize%int64(chunkSize) != 0 {
		chunkNumber += 1
	}

	return ChunkSummary{
		Size:   int64(chunkSize),
		Number: int32(chunkNumber),
	}
}
