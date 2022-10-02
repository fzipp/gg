// Copyright 2020 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package twp encodes/decodes data with XOR encryption as used by
// Thimbleweed Park and Delores games.
package twp

import (
	"errors"
	"io"

	"github.com/fzipp/gg/crypt/internal/transform"
	"github.com/fzipp/gg/ggdict"
)

// Key is an XOR key for Thimbleweed Park or Delores.
type Key struct {
	MagicBytes []byte
	Multiplier byte
}

func (key *Key) DecodingReader(r io.Reader, expectedSize int64) io.Reader {
	return transform.NewReader(r, newDecoder(key, expectedSize))
}

func (key *Key) EncodingWriter(w io.Writer, expectedSize int64) io.Writer {
	return transform.NewWriter(w, newEncoder(key, expectedSize))
}

func (key *Key) NeedsLoading() bool {
	return false
}

func (key *Key) LoadFrom(execFile string) error {
	return errors.New("this Key does not need to be loaded")
}

func (key *Key) GGDictFormat() ggdict.Format {
	return ggdict.FormatThimbleweed
}
