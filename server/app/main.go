package main

import (
	"crypto/rand"
	"encoding/binary"
	insecureRand "math/rand"
	"time"

	"Dur4nC2/cli"
)

// Attempt to seed insecure rand with secure rand, but we really
// don't care that much if it fails since it's insecure anyways
func init() {
	buf := make([]byte, 8)
	_, err := rand.Read(buf)
	if err != nil {
		insecureRand.Seed(int64(time.Now().Unix()))
	} else {
		insecureRand.Seed(int64(binary.LittleEndian.Uint64(buf)))
	}
}

func main() {
	cli.Execute()

}
