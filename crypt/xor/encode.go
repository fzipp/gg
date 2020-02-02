// Copyright 2020 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xor

import (
	"io"

	"github.com/fzipp/gg/crypt/internal/transform"
)

type encoder struct {
	key    *Key
	xorSum byte
}

func newEncoder(key *Key, expectedSize int64) transform.Transformer {
	return &encoder{key: key, xorSum: byte(expectedSize)}
}

func (e *encoder) Transform(dst, src []byte) {
	for i, b := range src {
		x := b ^ e.xorSum
		dst[i] = x ^ e.key.MagicBytes[i&0x0F] ^ byte(i)*e.key.Multiplier
		e.xorSum = x
	}
}

func EncodingWriter(w io.Writer, key *Key, expectedSize int64) io.Writer {
	return transform.NewWriter(w, newEncoder(key, expectedSize))
}
