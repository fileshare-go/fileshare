package osinfo

import (
	"os"
	"runtime"
	"strings"
)

func OsInfo() string {
	return strings.Join([]string{
		runtime.GOOS, runtime.GOARCH, getHostname(),
	}, ",")
}

func getHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		return "unknown"
	}
	return hostname
}
