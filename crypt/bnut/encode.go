// Copyright 2020 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bnut

import (
	"io"

	"github.com/fzipp/gg/crypt/internal/transform"
)

type encoder struct {
	cursor int
}

func newEncoder(expectedSize int64) transform.Transformer {
	return &encoder{cursor: int(expectedSize & 0xff)}
}

func (e *encoder) Transform(dst, src []byte) {
	for i := 0; i < len(src); i++ {
		e.cursor = (e.cursor + 1) % len(cryptKey)
		dst[i] ^= cryptKey[e.cursor]
	}
}

func EncodingWriter(w io.Writer, expectedSize int64) io.Writer {
	return transform.NewWriter(w, newEncoder(expectedSize))
}
