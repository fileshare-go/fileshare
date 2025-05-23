package fileutil

import (
	"os"
	"strings"

	"github.com/sirupsen/logrus"
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
	return list[len(list) - 1]
}


func CreateLockDir(lockDirectory string) {
	if FileExists(lockDirectory) {
		return
	}

	if err := os.Mkdir(lockDirectory, 0755); err != nil {
		logrus.Fatalf("While creating lockDirectory %s, %s", lockDirectory, err.Error())
	}
}
