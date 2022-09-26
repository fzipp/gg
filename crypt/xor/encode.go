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
	cursor byte
}

func (key *Key) newEncoder(expectedSize int64) transform.Transformer {
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

func (key *Key) EncodingWriter(w io.Writer, expectedSize int64) io.Writer {
	return transform.NewWriter(w, key.newEncoder(expectedSize))
}

type monkeyIslandEncoder struct {
	key    *MonkeyIslandKey
	cursor uint16
}

func (Key *MonkeyIslandKey) newEncoder(expectedSize int64) transform.Transformer {
	return &monkeyIslandEncoder{key: Key, cursor: uint16(uint16(expectedSize) + uint16(Key.Modifier))}
}

func (key *MonkeyIslandKey) EncodingWriter(w io.Writer, expectedSize int64) io.Writer {
	return transform.NewWriter(w, key.newEncoder(expectedSize))
}

func (d *monkeyIslandEncoder) Transform(dst, src []byte) {
	//var d.cursor uint16 = uint16(uint16(len(src)) + uint16(d.key.Modifier))
	for i, b := range src {
		x := b ^ d.key.MagicBytes1[uint8((uint8(d.cursor)+(d.key.Modifier))&0xFF)] ^ d.key.MagicBytes2[d.cursor]
		dst[i] = x
		d.cursor = uint16(d.cursor + uint16(d.key.MagicBytes1[uint8(d.cursor&0xFF)])&0xFFFF)
	}
}
