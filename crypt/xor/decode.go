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
	cursor byte
	xorSum byte
}

func (key *Key) newDecoder(expectedSize int64) transform.Transformer {
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

func (key *Key) DecodingReader(r io.Reader, expectedSize int64) io.Reader {
	return transform.NewReader(r, key.newDecoder(expectedSize))
}

type monkeyIslandDecoder struct {
	key    *MonkeyIslandKey
	cursor uint16
}

func (Key *MonkeyIslandKey) newDecoder(expectedSize int64) transform.Transformer {
	return &monkeyIslandDecoder{key: Key, cursor: uint16(uint16(expectedSize) + uint16(Key.Modifier))}
}

func (key *MonkeyIslandKey) DecodingReader(r io.Reader, expectedSize int64) io.Reader {
	return transform.NewReader(r, key.newDecoder(expectedSize))
}

func (d *monkeyIslandDecoder) Transform(dst, src []byte) {
	//var d.cursor uint16 = uint16(uint16(len(src)) + uint16(d.key.Modifier))
	for i, b := range src {
		x := b ^ d.key.MagicBytes1[uint8((uint8(d.cursor)+(d.key.Modifier))&0xFF)] ^ d.key.MagicBytes2[d.cursor]
		dst[i] = x
		d.cursor = uint16(d.cursor + uint16(d.key.MagicBytes1[uint8(d.cursor&0xFF)])&0xFFFF)
	}
}
