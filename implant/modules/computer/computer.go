package computer

import (
	"os"
	"runtime"
)

func GetOS() string {
	return runtime.GOOS
}

func GetArch() string {
	return runtime.GOARCH
}
func GetHostname() string {
	hostname, _ := os.Hostname()
	return hostname
}
