// Copyright 2022 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rtmi

import "github.com/fzipp/gg/crypt/internal/transform"

type decoder struct {
	key    *Key
	cursor uint16
}

func newDecoder(key *Key, expectedSize int64) transform.Transformer {
	return &decoder{key: key, cursor: uint16(expectedSize) + uint16(key.Modifier)}
}

func (d *decoder) Transform(dst, src []byte) {
	for i, b := range src {
		x := b ^ d.key.MagicBytes1[((uint8(d.cursor)+(d.key.Modifier))&0xFF)] ^ d.key.MagicBytes2[d.cursor]
		dst[i] = x
		d.cursor = d.cursor + uint16(d.key.MagicBytes1[uint8(d.cursor&0xFF)])&0xFFFF
	}
}
