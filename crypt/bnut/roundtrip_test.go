// Copyright 2020 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bnut_test

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/fzipp/gg/crypt/bnut"
)

var testBnutScript = `__<-"This is a test input."

TestRoom <-
{
 background = "TestRoom"

 enter = function()
 {
 }
}
`

func TestWriterReaderRoundTrip(t *testing.T) {
	original := []byte(testBnutScript)
	encodedBuf := &bytes.Buffer{}
	_, err := bnut.EncodingWriter(encodedBuf, int64(len(original))).Write(original)
	if err != nil {
		t.Errorf("encoding writer returned an error: %s", err)
	}
	encoded := encodedBuf.Bytes()

	decoded := make([]byte, len(encoded))
	_, err = bnut.DecodingReader(bytes.NewBuffer(encoded), int64(len(encoded))).Read(decoded)
	if err != nil {
		t.Errorf("decoding reader returned an error: %s", err)
	}

	if !reflect.DeepEqual(decoded, original) {
		t.Errorf("decoded data is not equal to original data! Original: %q vs. decoded: %q", string(original), string(decoded))
	}
}
