// Copyright 2020 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bnut

import (
	"io"

	"github.com/fzipp/gg/crypt/internal/transform"
)

type decoder struct {
	cursor int
}

func newDecoder(expectedSize int64) transform.Transformer {
	return &decoder{cursor: int(expectedSize & 0xff)}
}

func (d *decoder) Transform(dst, src []byte) {
	for i := 0; i < len(src); i++ {
		dst[i] = src[i] ^ cryptKey[d.cursor]
		d.cursor = (d.cursor + 1) % len(cryptKey)
	}
}

func DecodingReader(r io.Reader, expectedSize int64) io.Reader {
	return transform.NewReader(r, newDecoder(expectedSize))
}
