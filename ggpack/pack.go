// Copyright 2020 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ggpack

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"

	"github.com/fzipp/gg/crypt/bnut"
	"github.com/fzipp/gg/crypt/xor"
)

type Pack struct {
	reader    io.ReadSeeker
	directory directory
	xorKey    *xor.Key
}

func Open(path string) (*Pack, error) {
	return OpenUsingKey(path, xor.DefaultKey)
}

// OpenUsingKey is the same as Open, but uses a different key than the
// default key (xor.DefaultKey) for XOR decryption of the pack.
func OpenUsingKey(path string, key *xor.Key) (*Pack, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("could not open file '%s': %w", path, err)
	}
	pack := Pack{reader: file, xorKey: key}
	pack.directory, err = pack.readDirectory()
	if err != nil {
		file.Close()
		return nil, fmt.Errorf("could not read pack directory: %w", err)
	}
	return &pack, nil
}

func (p *Pack) Close() error {
	if closer, ok := p.reader.(io.Closer); ok {
		return closer.Close()
	}
	return nil
}

func (p *Pack) List() []DirectoryEntry {
	entries := make([]DirectoryEntry, 0, len(p.directory))
	for filename, entry := range p.directory {
		entries = append(entries, DirectoryEntry{
			Filename: filename,
			Size:     entry.Size,
		})
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Filename < entries[j].Filename
	})
	return entries
}

func (p *Pack) File(filename string) (r io.Reader, size int64, err error) {
	entry, exists := p.directory[filename]
	if !exists {
		return nil, 0, fmt.Errorf("file '%s' does not exist in pack", filename)
	}
	isBnut := filepath.Ext(filename) == ".bnut"
	r, err = p.entryReader(entry, isBnut)
	if err != nil {
		return nil, 0, fmt.Errorf("could not read file '%s' in pack", filename)
	}
	return r, entry.Size, nil
}

func (p *Pack) entryReader(entry entry, isBnut bool) (io.Reader, error) {
	_, err := p.reader.Seek(entry.Offset, io.SeekStart)
	if err != nil {
		return nil, fmt.Errorf("could not seek offset: %w", err)
	}
	limitedReader := &io.LimitedReader{R: p.reader, N: entry.Size}
	decodingReader := xor.DecodingReader(limitedReader, p.xorKey, entry.Size)
	if isBnut {
		return bnut.DecodingReader(decodingReader, entry.Size), nil
	}
	return decodingReader, nil
}

func (p *Pack) readDirectory() (directory, error) {
	entry, err := p.readDirectoryEntry()
	if err != nil {
		return nil, err
	}
	r, err := p.entryReader(entry, false)
	if err != nil {
		return nil, err
	}
	buf := make([]byte, entry.Size)
	_, err = io.ReadFull(r, buf)
	if err != nil {
		return nil, fmt.Errorf("could not read directory bytes: %w", err)
	}
	return readDirectory(buf)
}

func (p *Pack) readDirectoryEntry() (entry, error) {
	var data struct {
		Offset, Size uint32
	}
	if err := binary.Read(p.reader, binary.LittleEndian, &data); err != nil {
		return entry{}, fmt.Errorf("could not read directory offset and size: %w", err)
	}
	return entry{
		Offset: int64(data.Offset),
		Size:   int64(data.Size),
	}, nil
}
