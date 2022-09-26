// Copyright 2020 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xor_test

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/fzipp/gg/crypt/xor"
)

func TestWriterReaderRoundTrip(t *testing.T) {
	key := xor.DefaultKey

	original := []byte("This is a test.")
	encodedBuf := &bytes.Buffer{}
	_, err := key.EncodingWriter(encodedBuf, int64(len(original))).Write(original)
	if err != nil {
		t.Errorf("encoding writer returned an error: %s", err)
	}
	encoded := encodedBuf.Bytes()

	decoded := make([]byte, len(encoded))
	_, err = key.DecodingReader(bytes.NewBuffer(encoded), int64(len(encoded))).Read(decoded)
	if err != nil {
		t.Errorf("decoding reader returned an error: %s", err)
	}

	if !reflect.DeepEqual(decoded, original) {
		t.Errorf("decoded data is not equal to original data! Original: %q vs. decoded: %q", string(original), string(decoded))
	}
}
