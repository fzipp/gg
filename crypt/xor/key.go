// Copyright 2020 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xor

type Key struct {
	MagicBytes [16]byte
	Multiplier byte
}

// KnownKeys is a collection of XOR keys for ggpack files found in the wild.
// These keys differ slightly at MagicBytes[5] (0x5B vs. 0x56) and regarding
// the multiplier (0x6D vs. 0xAD). This is reflected in the names (e.g. "56ad")
// by which they can be referenced.
var KnownKeys = map[string]*Key {
	"5b6d": {
		MagicBytes: [...]byte{
			0x4F, 0xD0, 0xA0, 0xAC, 0x4A, 0x5B, 0xB9, 0xE5,
			0x93, 0x79, 0x45, 0xA5, 0xC1, 0xCB, 0x31, 0x93,
		},
		Multiplier: 0x6D,
	},
	"566d": {
		MagicBytes: [...]byte{
			0x4F, 0xD0, 0xA0, 0xAC, 0x4A, 0x56, 0xB9, 0xE5,
			0x93, 0x79, 0x45, 0xA5, 0xC1, 0xCB, 0x31, 0x93,
		},
		Multiplier: 0x6D,
	},
	"5bad": {
		MagicBytes: [...]byte{
			0x4F, 0xD0, 0xA0, 0xAC, 0x4A, 0x5B, 0xB9, 0xE5,
			0x93, 0x79, 0x45, 0xA5, 0xC1, 0xCB, 0x31, 0x93,
		},
		Multiplier: 0xAD,
	},
	"56ad": {
		MagicBytes: [...]byte{
			0x4F, 0xD0, 0xA0, 0xAC, 0x4A, 0x56, 0xB9, 0xE5,
			0x93, 0x79, 0x45, 0xA5, 0xC1, 0xCB, 0x31, 0x93,
		},
		Multiplier: 0xAD,
	},
}

// It is the default key since the author of this package happens to have
// only ggpack files encrypted with this key.
var DefaultKey = KnownKeys["56ad"]
