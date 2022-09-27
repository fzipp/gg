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
	"time"

	"github.com/fzipp/gg/crypt/bnut"
	"github.com/fzipp/gg/crypt/xor"
)

// Pack provides read access to the contents of a ggpack file.
// It implements the fs.FS, fs.ReadDirFS and io.Closer interfaces.
type Pack struct {
	reader    io.ReadSeeker
	modTime   time.Time
	directory *directory
	xorKey    xor.Key
}

func Open(path string) (*Pack, error) {
	return OpenUsingKey(path, xor.DefaultKey)
}

// OpenUsingKey is the same as Open, but uses a different key than the
// default key (xor.DefaultKey) for XOR decryption of the pack.
func OpenUsingKey(path string, key xor.Key) (*Pack, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("could not open file '%s': %w", path, err)
	}
	stat, err := f.Stat()
	if err != nil {
		return nil, fmt.Errorf("could get stat of file '%s': %w", path, err)
	}
	pack := Pack{reader: f, modTime: stat.ModTime(), xorKey: key}
	pack.directory, err = pack.readDirectory(key.UsesShortKeyIndices())
	if err != nil {
		f.Close()
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

// ReadDir reads the named directory
// and returns a list of directory entries sorted by filename.
//
// The only directory in a Pack is the root directory ".".
func (p *Pack) ReadDir(name string) ([]fs.DirEntry, error) {
	if !fs.ValidPath(name) {
		return nil, &fs.PathError{Op: "readdir", Path: name, Err: fs.ErrInvalid}
	}
	if name != "." {
		return nil, &fs.PathError{Op: "readdir", Path: name, Err: fs.ErrNotExist}
	}
	return (&rootDirFile{dir: p.directory}).ReadDir(0)
}

func (p *Pack) Open(name string) (fs.File, error) {
	if !fs.ValidPath(name) {
		return nil, &fs.PathError{Op: "open", Path: name, Err: fs.ErrInvalid}
	}
	if name == "." {
		return &rootDirFile{dir: p.directory}, nil
	}
	fi, exists := p.directory.lookup[name]
	if !exists {
		return nil, &fs.PathError{Op: "open", Path: name, Err: fs.ErrNotExist}
	}
	r, err := p.fileReader(fi)
	if err != nil {
		return nil, &fs.PathError{Op: "open", Path: name, Err: err}
	}
	return &file{stat: fi, r: r}, nil
}

func (p *Pack) fileReader(fi *fileInfo) (io.Reader, error) {
	_, err := p.reader.Seek(fi.packOffset, io.SeekStart)
	if err != nil {
		return nil, fmt.Errorf("could not seek offset: %w", err)
	}
	limitedReader := &io.LimitedReader{R: p.reader, N: fi.size}
	decodingReader := p.xorKey.DecodingReader(limitedReader, fi.size)
	switch filepath.Ext(fi.name) {
	case ".bank":
		// FMOD bank files are not XOR encrypted
		return limitedReader, nil
	case ".bnut":
		return bnut.DecodingReader(decodingReader, fi.size), nil
	}
	return decodingReader, nil
}

func (p *Pack) readDirectory(shortStringIndices bool) (*directory, error) {
	root, err := p.readRootInfo()
	if err != nil {
		return nil, err
	}
	r, err := p.fileReader(root)
	if err != nil {
		return nil, err
	}
	buf := make([]byte, root.size)
	_, err = io.ReadFull(r, buf)
	if err != nil {
		return nil, fmt.Errorf("could not read directory bytes: %w", err)
	}
	return readDirectory(buf, root, shortStringIndices)
}

func (p *Pack) readRootInfo() (*fileInfo, error) {
	var data struct {
		Offset, Size uint32
	}
	if err := binary.Read(p.reader, binary.LittleEndian, &data); err != nil {
		return nil, fmt.Errorf("could not read directory offset and size: %w", err)
	}
	return &fileInfo{
		name:       ".",
		mode:       fs.ModeDir,
		size:       int64(data.Size),
		modTime:    p.modTime,
		packOffset: int64(data.Offset),
	}, nil
}
