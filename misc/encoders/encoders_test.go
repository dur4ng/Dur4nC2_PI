package encoders

import (
	"bytes"
	"crypto/rand"
	insecureRand "math/rand"
	"testing"

	implantEncoders "github.com/bishopfox/sliver/implant/sliver/encoders"
)

func randomData() []byte {
	buf := make([]byte, insecureRand.Intn(1024))
	rand.Read(buf)
	return buf
}

func TestNopNonce(t *testing.T) {
	nop := NopNonce()
	_, _, err := EncoderFromNonce(nop)
	if err != nil {
		t.Errorf("Nop nonce returned error %s", err)
	}

	nop2 := implantEncoders.NopNonce()
	_, _, err = EncoderFromNonce(nop2)
	if err != nil {
		t.Errorf("Nop nonce returned error %s", err)
	}
}

func TestRandomEncoder(t *testing.T) {
	for index := 0; index < 20; index++ {
		sample := randomData()

		nonce, encoder := RandomEncoder()
		_, encoder2, err := EncoderFromNonce(nonce)
		if err != nil {
			t.Errorf("RandomEncoder() nonce returned error %s", err)
		}
		output := encoder.Encode(sample)
		data, err := encoder2.Decode(output)
		if err != nil {
			t.Errorf("RandomEncoder() encoder2 returned error %s", err)
		}
		if !bytes.Equal(sample, data) {
			t.Errorf("RandomEncoder() encoder2 failed to decode encoder data %s", err)
		}

		nonce, encoder = implantEncoders.RandomEncoder()
		_, encoder2, err = implantEncoders.EncoderFromNonce(nonce)
		if err != nil {
			t.Errorf("RandomEncoder() nonce returned error %s", err)
		}
		output = encoder.Encode(sample)
		data, err = encoder2.Decode(output)
		if err != nil {
			t.Errorf("RandomEncoder() encoder2 returned error %s", err)
		}
		if !bytes.Equal(sample, data) {
			t.Errorf("RandomEncoder() encoder2 failed to decode encoder data %s", err)
		}

		nonce, encoder = RandomEncoder()
		_, encoder2, err = implantEncoders.EncoderFromNonce(nonce)
		if err != nil {
			t.Errorf("RandomEncoder() nonce returned error %s", err)
		}
		output = encoder.Encode(sample)
		data, err = encoder2.Decode(output)
		if err != nil {
			t.Errorf("RandomEncoder() encoder2 returned error %s", err)
		}
		if !bytes.Equal(sample, data) {
			t.Errorf("RandomEncoder() encoder2 failed to decode encoder data %s", err)
		}

		nonce, encoder = implantEncoders.RandomEncoder()
		_, encoder2, err = EncoderFromNonce(nonce)
		if err != nil {
			t.Errorf("RandomEncoder() nonce returned error %s", err)
		}
		output = encoder.Encode(sample)
		data, err = encoder2.Decode(output)
		if err != nil {
			t.Errorf("RandomEncoder() encoder2 returned error %s", err)
		}
		if !bytes.Equal(sample, data) {
			t.Errorf("RandomEncoder() encoder2 failed to decode encoder data %s", err)
		}

	}
}
