package util

import (
	"encoding/base64"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

var once sync.Once
var info string

func OsInfo() string {
	once.Do(func() {
		infoStr := strings.Join([]string{
			runtime.GOOS, runtime.GOARCH, getHostname(),
		}, ",")
		info = base64.StdEncoding.EncodeToString([]byte(infoStr))
	})
	return info
}

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

func getHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		return "unknown"
	}
	return hostname
}
