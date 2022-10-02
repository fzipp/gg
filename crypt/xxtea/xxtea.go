// Copyright 2020 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package xxtea is an implementation of the XXTEA (Corrected Block TEA) block
// cipher as described in:
// David J. Wheeler and Roger M. Needham (October 1998). "Correction to XTEA".
// Computer Laboratory, Cambridge University
// https://www.movable-type.co.uk/scripts/xxtea.pdf
package xxtea

// Key is a 128-bit key.
type Key [4]uint32

func Encrypt(p []byte, k Key) []byte {
	v, rest := bytesToWords(p)
	encrypt(v, k)
	return append(wordsToBytes(v), rest...)
}

func Decrypt(p []byte, k Key) []byte {
	v, rest := bytesToWords(p)
	decrypt(v, k)
	return append(wordsToBytes(v), rest...)
}

const delta = 0x9e3779b9

func encrypt(v []uint32, k Key) {
	n := len(v)
	if n <= 1 {
		return
	}
	q := 6 + 52/n
	sum := 0
	z := v[n-1]
	for ; q > 0; q-- {
		sum += delta
		e := (sum >> 2) & 3
		for p := 0; p < n-1; p++ {
			y := v[p+1]
			v[p] += mx(y, z, sum, p, e, k)
			z = v[p]
		}
		y := v[0]
		v[n-1] += mx(y, z, sum, n-1, e, k)
		z = v[n-1]
	}
}

func decrypt(v []uint32, k Key) {
	n := len(v)
	if n <= 1 {
		return
	}
	q := 6 + 52/n
	sum := q * delta
	y := v[0]
	for ; q > 0; q-- {
		e := (sum >> 2) & 3
		for p := n - 1; p > 0; p-- {
			z := v[p-1]
			v[p] -= mx(y, z, sum, p, e, k)
			y = v[p]
		}
		z := v[n-1]
		v[0] -= mx(y, z, sum, 0, e, k)
		y = v[0]
		sum -= delta
	}
}

func mx(y, z uint32, sum, p, e int, k Key) uint32 {
	return ((z>>5 ^ y<<2) + (y>>3 ^ z<<4)) ^ ((uint32(sum) ^ y) + (k[(p&3)^e] ^ z))
}
