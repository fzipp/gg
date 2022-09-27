// Copyright 2020 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package twp

import "github.com/fzipp/gg/crypt/internal/transform"

type decoder struct {
	key    *Key
	cursor byte
	xorSum byte
}

func newDecoder(key *Key, expectedSize int64) transform.Transformer {
	return &decoder{key: key, xorSum: byte(expectedSize)}
}

func (d *decoder) Transform(dst, src []byte) {
	for i, b := range src {
		x := b ^ d.key.MagicBytes[d.cursor&0x0F] ^ d.cursor*d.key.Multiplier
		dst[i] = x ^ d.xorSum
		d.xorSum = x
		d.cursor++
	}
}
