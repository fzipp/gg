// Copyright 2020 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xor

import (
	"io"

	"github.com/fzipp/gg/crypt/internal/transform"
)

type decoder struct {
	key    *Key
	xorSum byte
}

func newDecoder(key *Key, expectedSize int64) transform.Transformer {
	return &decoder{key: key, xorSum: byte(expectedSize)}
}

func (d *decoder) Transform(dst, src []byte) {
	for i, b := range src {
		x := b ^ d.key.MagicBytes[i&0x0F] ^ byte(i)*d.key.Multiplier
		dst[i] = x ^ d.xorSum
		d.xorSum = x
	}
}

func DecodingReader(r io.Reader, key *Key, expectedSize int64) io.Reader {
	return transform.NewReader(r, newDecoder(key, expectedSize))
}
