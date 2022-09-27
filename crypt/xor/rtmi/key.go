// Copyright 2022 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package rtmi encodes/decodes data with XOR encryption as used by
// Return to Monkey Island.
package rtmi

import (
	"crypto/md5"
	"errors"
	"io"
	"os"

	"github.com/fzipp/gg/crypt/internal/transform"
)

// Key is an XOR key for Return to Monkey Island.
type Key struct {
	MagicBytes1 [256]byte
	MagicBytes2 [65536]byte
	Modifier    byte
	loaded      bool
}

func (key *Key) DecodingReader(r io.Reader, expectedSize int64) io.Reader {
	return transform.NewReader(r, newDecoder(key, expectedSize))
}

func (key *Key) EncodingWriter(w io.Writer, expectedSize int64) io.Writer {
	return transform.NewWriter(w, newEncoder(key, expectedSize))
}

func (key *Key) UsesShortKeyIndices() bool {
	return true
}

func (key *Key) NeedsLoading() bool {
	return !key.loaded
}

func (key *Key) LoadFrom(execFile string) error {
	if key.loaded {
		return errors.New("this key does not need to be loaded")
	}

	data, err := os.ReadFile(execFile)
	if err != nil {
		return err
	}

	isMatch := func(firstValue byte, length int, checksum [16]byte, startIndex int) bool {
		if data[startIndex] != firstValue {
			return false
		}
		sum := md5.Sum(data[startIndex:(length + startIndex)])
		for i := 0; i < 16; i++ {
			if sum[i] != checksum[i] {
				return false
			}
		}
		return true
	}

	found1 := false
	found2 := false

	for i := 0; i < len(data)-256; i++ {
		if isMatch(0xD5, 256, [16]byte{0xB1, 0x90, 0xC4, 0x21, 0xFE, 0x7F, 0xEA, 0xFC, 0x77, 0xC5, 0x17, 0xA2, 0x32, 0xAB, 0xBB, 0x4C}, i) {
			for x := 0; x < 256; x++ {
				key.MagicBytes1[x] = data[i+x]
			}
			found1 = true
			break
		}
	}

	for i := 0; i < len(data)-65536; i++ {
		if isMatch(0xF7, 65536, [16]byte{0x7F, 0xAA, 0xF6, 0x57, 0x4F, 0x27, 0xEB, 0xD9, 0xD2, 0x74, 0x4C, 0xC6, 0x8E, 0x41, 0x15, 0xC8}, i) {
			for x := 0; x < 65536; x++ {
				key.MagicBytes2[x] = data[i+x]
			}
			found2 = true
			break
		}
	}

	if !found1 || !found2 {
		return errors.New("one or both keys could not be found")
	}

	return nil
}
