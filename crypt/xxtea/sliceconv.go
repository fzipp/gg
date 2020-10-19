// Copyright 2020 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xxtea

import "encoding/binary"

const wordSize = 4

var endianness = binary.LittleEndian

func bytesToWords(b []byte) (w []uint32, rest []byte) {
	w = make([]uint32, len(b)/wordSize)
	for i := range w {
		w[i] = endianness.Uint32(b[i*wordSize : (i+1)*wordSize])
	}
	rest = b[len(b)-len(b)%wordSize:]
	return w, rest
}

func wordsToBytes(w []uint32) (b []byte) {
	b = make([]byte, len(w)*wordSize)
	for i := range w {
		endianness.PutUint32(b[i*wordSize:], w[i])
	}
	return b
}
