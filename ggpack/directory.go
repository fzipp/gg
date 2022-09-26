// Copyright 2020 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ggpack

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"sort"
	"time"

	"github.com/fzipp/gg/ggdict"
)

type rootDirFile struct {
	dir    *directory
	offset int
}

func (f *rootDirFile) Stat() (fs.FileInfo, error) { return f.dir.info, nil }

func (f *rootDirFile) Read([]byte) (int, error) {
	return 0, &fs.PathError{Op: "read", Path: ".", Err: errors.New("is a directory")}
}

func (f *rootDirFile) Close() error { return nil }

func (f *rootDirFile) ReadDir(count int) ([]fs.DirEntry, error) {
	n := len(f.dir.entries) - f.offset
	if count > 0 && n > count {
		n = count
	}
	if n == 0 {
		if count <= 0 {
			return nil, nil
		}
		return nil, io.EOF
	}
	list := make([]fs.DirEntry, n)
	for i := range list {
		list[i] = f.dir.entries[f.offset+i]
	}
	f.offset += n
	return list, nil
}

type file struct {
	stat fs.FileInfo
	r    io.Reader
}

func (f *file) Stat() (fs.FileInfo, error)       { return f.stat, nil }
func (f *file) Read(b []byte) (n int, err error) { return f.r.Read(b) }
func (f *file) Close() error                     { return nil }

type fileInfo struct {
	name       string
	mode       fs.FileMode
	size       int64
	modTime    time.Time
	packOffset int64
}

func (fi *fileInfo) Name() string               { return fi.name }
func (fi *fileInfo) IsDir() bool                { return fi.Mode().IsDir() }
func (fi *fileInfo) Type() fs.FileMode          { return fi.Mode().Type() }
func (fi *fileInfo) Info() (fs.FileInfo, error) { return fi, nil }
func (fi *fileInfo) Size() int64                { return fi.size }
func (fi *fileInfo) Mode() fs.FileMode          { return fi.mode }
func (fi *fileInfo) ModTime() time.Time         { return fi.modTime }
func (fi *fileInfo) Sys() any                   { return nil }

type directory struct {
	info    fs.FileInfo
	entries []fs.DirEntry
	lookup  map[string]*fileInfo
}

const (
	keyFiles    = "files"
	keyFilename = "filename"
	keyOffset   = "offset"
	keySize     = "size"
)

func readDirectory(buf []byte, root *fileInfo, shortStringIndices bool) (*directory, error) {
	directoryDict, err := ggdict.Unmarshal(buf, shortStringIndices)
	if err != nil {
		return nil, fmt.Errorf("could not read directory: %w", err)
	}
	return directoryFrom(directoryDict, root)
}

func directoryFrom(dict map[string]any, root *fileInfo) (*directory, error) {
	files, ok := dict[keyFiles].([]any)
	if !ok {
		return nil, fmt.Errorf("%q is not an array", keyFiles)
	}
	dir := &directory{
		info:    root,
		lookup:  make(map[string]*fileInfo, len(files)),
		entries: make([]fs.DirEntry, 0, len(files)),
	}
	for _, fileEntry := range files {
		entryDict := fileEntry.(map[string]any)
		filename, ok := entryDict[keyFilename].(string)
		if !ok {
			return nil, fmt.Errorf("%q is not a string", keyFilename)
		}
		offset, ok := entryDict[keyOffset].(int)
		if !ok {
			return nil, fmt.Errorf("%q is not an int", keyOffset)
		}
		size, ok := entryDict[keySize].(int)
		if !ok {
			return nil, fmt.Errorf("%q is not an int", keySize)
		}
		fi := &fileInfo{
			name:       filename,
			mode:       0,
			size:       int64(size),
			modTime:    root.modTime,
			packOffset: int64(offset),
		}
		dir.lookup[filename] = fi
		dir.entries = append(dir.entries, fi)
	}
	sort.Slice(dir.entries, func(i, j int) bool {
		return dir.entries[i].Name() < dir.entries[j].Name()
	})
	return dir, nil
}
