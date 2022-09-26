// Copyright 2020 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package xor encodes/decodes data with "unbreakable" XOR encryption.
package xor

import (
	"crypto/md5"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/fzipp/gg/crypt/internal/transform"
)

type KeyInterface interface {
	newDecoder(expectedSize int64) transform.Transformer
	DecodingReader(r io.Reader, expectedSize int64) io.Reader
	newEncoder(expectedSize int64) transform.Transformer
	EncodingWriter(w io.Writer, expectedSize int64) io.Writer
	NeedsLoading() bool
	LoadKey(packFileName string) error
	UsesShortKeyIndices() bool
}

type Key struct {
	MagicBytes []byte
	Multiplier byte
}

// Special Key for RtMI contains two byte arrays.
type MonkeyIslandKey struct {
	MagicBytes1 [256]byte
	MagicBytes2 [65536]byte
	Modifier    byte
	loaded      bool
}

func (key *Key) NeedsLoading() bool {
	return false
}

func (key *Key) LoadKey(packFileName string) error {
	return errors.New("this Key does not need to be loaded")
}

func (key *Key) UsesShortKeyIndices() bool {
	return false
}

func (key *MonkeyIslandKey) UsesShortKeyIndices() bool {
	return true
}

func (key *MonkeyIslandKey) NeedsLoading() bool {
	return !key.loaded
}

func (key *MonkeyIslandKey) LoadKey(packFileName string) error {
	if key.loaded {
		return errors.New("this Key does not need to be loaded")
	}

	packFileName, err := filepath.Abs(packFileName)
	if err != nil {
		return err
	}

	directory := filepath.Dir(packFileName)

	gameExecutableName := filepath.Join(directory, "Return to Monkey Island.exe")

	if _, err := os.Stat(gameExecutableName); errors.Is(err, os.ErrNotExist) {
		return err
	}

	data, err := ioutil.ReadFile(gameExecutableName)
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

// KnownKeys is a collection of XOR keys for ggpack files found in the wild.
// These keys differ slightly at MagicBytes[5] (0x5B vs. 0x56) and regarding
// the multiplier (0x6D vs. 0xAD). This is reflected in the names (e.g. "56ad")
// by which they can be referenced.
var KnownKeys = map[string]KeyInterface{
	"5b6d": &Key{
		MagicBytes: []byte{
			0x4F, 0xD0, 0xA0, 0xAC, 0x4A, 0x5B, 0xB9, 0xE5,
			0x93, 0x79, 0x45, 0xA5, 0xC1, 0xCB, 0x31, 0x93,
		},
		Multiplier: 0x6D,
	},
	"566d": &Key{
		MagicBytes: []byte{
			0x4F, 0xD0, 0xA0, 0xAC, 0x4A, 0x56, 0xB9, 0xE5,
			0x93, 0x79, 0x45, 0xA5, 0xC1, 0xCB, 0x31, 0x93,
		},
		Multiplier: 0x6D,
	},
	"5bad": &Key{
		MagicBytes: []byte{
			0x4F, 0xD0, 0xA0, 0xAC, 0x4A, 0x5B, 0xB9, 0xE5,
			0x93, 0x79, 0x45, 0xA5, 0xC1, 0xCB, 0x31, 0x93,
		},
		Multiplier: 0xAD,
	},
	"56ad": &Key{
		MagicBytes: []byte{
			0x4F, 0xD0, 0xA0, 0xAC, 0x4A, 0x56, 0xB9, 0xE5,
			0x93, 0x79, 0x45, 0xA5, 0xC1, 0xCB, 0x31, 0x93,
		},
		Multiplier: 0xAD,
	},
	"delores": &Key{
		MagicBytes: []byte{
			0x3F, 0x41, 0x41, 0x60, 0x95, 0x87, 0x4A, 0xE6,
			0x34, 0xC6, 0x3A, 0x86, 0x29, 0x27, 0x77, 0x8D,
			0x38, 0xB4, 0x96, 0xC9, 0x38, 0xB4, 0x96, 0xC9,
			0x00, 0xE0, 0x0A, 0xC6, 0x00, 0xE0, 0x0A, 0xC6,
			0x00, 0x3C, 0x1C, 0xC6, 0x00, 0x3C, 0x1C, 0xC6,
			0x00, 0xE4, 0x40, 0xC6, 0x00, 0xE4, 0x40, 0xC6,
		},
		Multiplier: 0x6D,
	},
	"rtmi": &MonkeyIslandKey{
		Modifier: 0x78,
		loaded:   false,
	},
}

// It is the default key since the author of this package happens to have
// only ggpack files encrypted with this key.
var DefaultKey = KnownKeys["56ad"]
