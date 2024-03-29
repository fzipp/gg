// Copyright 2020 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xor

import (
	"io"

	"github.com/fzipp/gg/crypt/xor/rtmi"
	"github.com/fzipp/gg/crypt/xor/twp"
	"github.com/fzipp/gg/ggdict"
)

// Key is an XOR key for ggpack files.
type Key interface {
	DecodingReader(r io.Reader, expectedSize int64) io.Reader
	EncodingWriter(w io.Writer, expectedSize int64) io.Writer

	// NeedsLoading returns true if the key needs to be loaded
	// from the executable file via LoadFrom.
	NeedsLoading() bool
	// LoadFrom loads the key from the executable file.
	// This is only necessary if NeedsLoading returns true.
	LoadFrom(execFile string) error

	GGDictFormat() ggdict.Format
}

// KnownKeys is a collection of XOR keys for ggpack files found in the wild.
//
// For Thimbleweed Park multiple keys are known. They differ slightly at
// MagicBytes[5] (0x5B vs. 0x56) and regarding the multiplier (0x6D vs. 0xAD).
// This is reflected in the names (e.g. "thimbleweed-56ad") by which they can
// be referenced.
var KnownKeys = map[string]Key{
	// Thimbleweed Park
	"thimbleweed": &twp.Key{
		MagicBytes: []byte{
			0x4F, 0xD0, 0xA0, 0xAC, 0x4A, 0x56, 0xB9, 0xE5,
			0x93, 0x79, 0x45, 0xA5, 0xC1, 0xCB, 0x31, 0x93,
		},
		Multiplier: 0xAD,
	},
	// Thimbleweed Park
	"thimbleweed-5b6d": &twp.Key{
		MagicBytes: []byte{
			0x4F, 0xD0, 0xA0, 0xAC, 0x4A, 0x5B, 0xB9, 0xE5,
			0x93, 0x79, 0x45, 0xA5, 0xC1, 0xCB, 0x31, 0x93,
		},
		Multiplier: 0x6D,
	},
	// Thimbleweed Park
	"thimbleweed-566d": &twp.Key{
		MagicBytes: []byte{
			0x4F, 0xD0, 0xA0, 0xAC, 0x4A, 0x56, 0xB9, 0xE5,
			0x93, 0x79, 0x45, 0xA5, 0xC1, 0xCB, 0x31, 0x93,
		},
		Multiplier: 0x6D,
	},
	// Thimbleweed Park
	"thimbleweed-5bad": &twp.Key{
		MagicBytes: []byte{
			0x4F, 0xD0, 0xA0, 0xAC, 0x4A, 0x5B, 0xB9, 0xE5,
			0x93, 0x79, 0x45, 0xA5, 0xC1, 0xCB, 0x31, 0x93,
		},
		Multiplier: 0xAD,
	},
	// Delores
	"delores": &twp.Key{
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
	// Return to Monkey Island
	"monkey": &rtmi.Key{
		// The magic bytes of this key need to be loaded from
		// the executable file via Key.LoadFrom.
		Modifier: 0x78,
	},
}

// It is the default key since the author of this package happens to have
// only ggpack files encrypted with this key.
var DefaultKey = KnownKeys["thimbleweed"]
