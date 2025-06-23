package model

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/chanmaoganda/fileshare/internal/pkg/util"
	"github.com/sirupsen/logrus"
)

type FileInfo struct {
	Filename       string `gorm:"primaryKey;size:64"`
	Sha256         string `gorm:"primaryKey;size:64"`
	ChunkSize      int64
	ChunkNumber    int32
	FileSize       int64
	UploadedChunks string
	Link           []ShareLink `gorm:"foreignKey:Sha256;references:Sha256"`
	Record         []Record    `gorm:"foreignKey:Sha256;references:Sha256"`
}

func (f *FileInfo) GetUploadedChunks() []int32 {
	var chunks []int32
	if err := json.Unmarshal([]byte(f.UploadedChunks), &chunks); err != nil {
		logrus.Error(err)
		return []int32{}
	}

	return chunks
}

func (f *FileInfo) GetMissingChunks() []int32 {
	loaded := f.GetUploadedChunks()
	all := f.GetAllChunks()

	return util.MissingElementsInSortedList(all, loaded)
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
	merged := util.MergeList(loaded, newChunks)

	bytes, err := json.Marshal(merged)
	if err != nil {
		logrus.Error(err)
		return
	}

	f.UploadedChunks = string(bytes)
}

func (f *FileInfo) ValidateChunks(cache_directory, download_directory string) bool {
	filePath := fmt.Sprintf("%s/%s", download_directory, f.Filename)
	logrus.Debug("[Validate] File: ", util.Render(filePath))

	if util.FileExists(filePath) {
		checkSum, err := util.CalculateFileSHA256(filePath)
		if err != nil {
			logrus.Warn("[Validate]", err)
		}

		if checkSum == f.Sha256 {
			logrus.Debugf("Existing file [%s] matches checksum!", filePath)
			return true
		}
		logrus.Debugf("Existing file [%s] does not match checksum, remaking new file", filePath)
	}

	if err := f.RemakeFile(cache_directory, filePath); err != nil {
		logrus.Error("[Validate] ", err)
	}

	checkSum, err := util.CalculateFileSHA256(filePath)
	if err != nil {
		logrus.Error("[Validate] ", err)
		return false
	}

	return checkSum == f.Sha256
}

func (f *FileInfo) RemakeFile(cache_directory, filePath string) error {
	out, err := os.Create(filePath)
	if err != nil {
		return err
	}

	for _, index := range f.GetUploadedChunks() {
		in, err := os.Open(fmt.Sprintf("%s/%s/%d", cache_directory, f.Sha256, index))
		if err != nil {
			return err
		}

		_, err = io.Copy(out, in)
		if err != nil {
			return err
		}
		if err := in.Close(); err != nil {
			return err
		}
	}
	return out.Close()
}
