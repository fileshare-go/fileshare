package util

import (
	"encoding/base64"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func OsInfo() string {
	info := strings.Join([]string{
		runtime.GOOS, runtime.GOARCH, getHostname(),
	}, ",")
	return base64.StdEncoding.EncodeToString([]byte(info))
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
