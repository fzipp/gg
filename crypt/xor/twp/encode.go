// Copyright 2020 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package twp

import "github.com/fzipp/gg/crypt/internal/transform"

type encoder struct {
	key    *Key
	xorSum byte
	cursor byte
}

func newEncoder(key *Key, expectedSize int64) transform.Transformer {
	return &encoder{key: key, xorSum: byte(expectedSize)}
}

func (e *encoder) Transform(dst, src []byte) {
	for i, b := range src {
		x := b ^ e.xorSum
		dst[i] = x ^ e.key.MagicBytes[e.cursor&0x0F] ^ e.cursor*e.key.Multiplier
		e.xorSum = x
		e.cursor++
	}
}
