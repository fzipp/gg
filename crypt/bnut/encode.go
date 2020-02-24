// Copyright 2020 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bnut

import (
	"io"

	"github.com/fzipp/gg/crypt/internal/transform"
)

type transformer struct {
	cursor int
}

func newTransformer(expectedSize int64) transform.Transformer {
	return &transformer{cursor: int(expectedSize & 0xff)}
}

func (t *transformer) Transform(dst, src []byte) {
	for i := 0; i < len(src); i++ {
		dst[i] = src[i] ^ cryptKey[t.cursor]
		t.cursor = (t.cursor + 1) % len(cryptKey)
	}
}

func EncodingWriter(w io.Writer, expectedSize int64) io.Writer {
	return transform.NewWriter(w, newTransformer(expectedSize))
}
