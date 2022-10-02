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
	"github.com/fzipp/gg/ggdict"
)

// Key is an XOR key for Return to Monkey Island.
type Key struct {
	// MagicBytes1 is the short (256 bytes) key.
	MagicBytes1 []byte
	// MagicBytes2 is the long (65536 bytes) key.
	MagicBytes2 []byte
	Modifier    byte
}

func (key *Key) DecodingReader(r io.Reader, expectedSize int64) io.Reader {
	return transform.NewReader(r, newDecoder(key, expectedSize))
}

func (key *Key) EncodingWriter(w io.Writer, expectedSize int64) io.Writer {
	return transform.NewWriter(w, newEncoder(key, expectedSize))
}

func (key *Key) GGDictFormat() ggdict.Format {
	return ggdict.FormatMonkey
}

func (key *Key) NeedsLoading() bool {
	return key.MagicBytes1 == nil || key.MagicBytes2 == nil
}

func (key *Key) LoadFrom(execFile string) error {
	data, err := os.ReadFile(execFile)
	if err != nil {
		return err
	}
	key.MagicBytes1 = extractKey(data, 256, 0xD5, &[...]byte{
		0xB1, 0x90, 0xC4, 0x21, 0xFE, 0x7F, 0xEA, 0xFC,
		0x77, 0xC5, 0x17, 0xA2, 0x32, 0xAB, 0xBB, 0x4C,
	})
	key.MagicBytes2 = extractKey(data, 65536, 0xF7, &[...]byte{
		0x7F, 0xAA, 0xF6, 0x57, 0x4F, 0x27, 0xEB, 0xD9,
		0xD2, 0x74, 0x4C, 0xC6, 0x8E, 0x41, 0x15, 0xC8,
	})
	if key.NeedsLoading() {
		return errors.New("one or both keys could not be found")
	}
	return nil
}

func extractKey(data []byte, length int, firstByte byte, md5sum *[16]byte) []byte {
	for i := 0; i < len(data)-length; i++ {
		if isMatch(data[i:], length, firstByte, md5sum) {
			key := make([]byte, length)
			copy(key, data[i:i+length])
			return key
		}
	}
	return nil
}

func isMatch(data []byte, length int, firstByte byte, md5sum *[16]byte) bool {
	return data[0] == firstByte && md5.Sum(data[:length]) == *md5sum
}
