package encoders

import "encoding/base64"

// Base64EncoderID - EncoderID
const Base64EncoderID = 13

// Base64 Encoder
type Base64 struct{}

var base64Alphabet = "a0b2c5def6hijklmnopqr_st-uvwxyzA1B3C4DEFGHIJKLM7NO9PQR8ST+UVWXYZ"
var base64encoder = base64.NewEncoding(base64Alphabet).WithPadding(base64.NoPadding)

// Encode - Base64 Encode
func (e Base64) Encode(data []byte) []byte {
	return []byte(base64encoder.EncodeToString(data))
}

// Decode - Base64 Decode
func (e Base64) Decode(data []byte) ([]byte, error) {
	return base64encoder.DecodeString(string(data))
}
