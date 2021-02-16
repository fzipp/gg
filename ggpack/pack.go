// Copyright 2020 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package ggpack reads and writes ggpack files.
package ggpack

import (
	"encoding/binary"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/fzipp/gg/crypt/bnut"
	"github.com/fzipp/gg/crypt/xor"
)

type Pack struct {
	reader    io.ReadSeeker
	modTime   time.Time
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
	stat, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("could get stat of file '%s': %w", path, err)
	}
	pack := Pack{reader: file, modTime: stat.ModTime(), xorKey: key}
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

func (p *Pack) ReadDir(name string) ([]fs.DirEntry, error) {
	if name != "." {
		return nil, fs.ErrNotExist // TODO: is this the best error?
	}
	list := make([]fs.DirEntry, 0, len(p.directory))
	for filename, entry := range p.directory {
		list = append(list, fileDirEntry{
			name:    filename,
			size:    entry.Size,
			modTime: p.modTime,
		})
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].Name() < list[j].Name()
	})
	return list, nil
}

func (p *Pack) Open(name string) (fs.File, error) {
	entry, exists := p.directory[name]
	if !exists {
		return nil, fs.ErrNotExist
	}
	isBnut := filepath.Ext(name) == ".bnut"
	r, err := p.entryReader(entry, isBnut)
	if err != nil {
		return nil, fmt.Errorf("could not read file '%s' in pack", name)
	}
	return packFile{
		stat: fileDirEntry{name: name, size: entry.Size, modTime: p.modTime},
		r:    r,
	}, nil
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
