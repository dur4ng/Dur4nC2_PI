package compressors

import (
	"bytes"
	"compress/flate"
)

func DeflateBuf(data []byte) []byte {
	var buf bytes.Buffer
	flateWriter, _ := flate.NewWriter(&buf, flate.BestCompression)
	flateWriter.Write(data)
	flateWriter.Close()
	return buf.Bytes()
}
