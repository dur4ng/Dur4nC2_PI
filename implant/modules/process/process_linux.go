package process

import "os"

func GetPID() int {
	return os.Getpid()
}
func GetFilename() string {
	return "golang"
}
