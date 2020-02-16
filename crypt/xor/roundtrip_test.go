package xor

import (
	"bytes"
	"reflect"
	"testing"
)

var testKey = DefaultKey

func TestTransformerRoundTrip(t *testing.T) {
	original := []byte("secret")
	size := len(original)
	enc := &encoder{xorSum: byte(size), key: testKey}
	encoded := make([]byte, size)
	enc.Transform(encoded, original)

	decoded := make([]byte, size)
	dec := &decoder{xorSum: byte(size), key: testKey}
	dec.Transform(decoded, encoded)

	if !reflect.DeepEqual(decoded, original) {
		t.Errorf("decoded data is not equal to original data! Original: %q vs. decoded: %q", string(original), string(decoded))
	}
}

func TestWriterReaderRoundTrip(t *testing.T) {
	original := []byte("secret")
	encodedBuf := &bytes.Buffer{}
	_, err := EncodingWriter(encodedBuf, testKey, int64(len(original))).Write(original)
	if err != nil {
		t.Errorf("encoding writer returned an error: %s", err)
	}
	encoded := encodedBuf.Bytes()

	decoded := make([]byte, len(encoded))
	_, err = DecodingReader(bytes.NewBuffer(encoded), testKey, int64(len(encoded))).Read(decoded)
	if err != nil {
		t.Errorf("decoding reader returned an error: %s", err)
	}

	if !reflect.DeepEqual(decoded, original) {
		t.Errorf("decoded data is not equal to original data! Original: %q vs. decoded: %q", string(original), string(decoded))
	}
}
