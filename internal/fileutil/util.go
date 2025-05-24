package fileutil

import (
	"os"
	"strings"
)

func FileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	if err == nil || os.IsExist(err) {
		return true
	}
	return false
}

func GetFileName(path string) string {
	list := strings.Split(path, "/")
	return list[len(list)-1]
}
