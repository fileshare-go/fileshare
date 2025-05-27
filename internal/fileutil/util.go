package fileutil

import (
	"path/filepath"
	"os"
)

func FileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	if err == nil || os.IsExist(err) {
		return true
	}
	return false
}

func GetFileName(path string) string {
	return filepath.Base(path)
}
