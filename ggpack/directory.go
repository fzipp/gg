// Copyright 2020 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ggpack

import (
	"fmt"
	"io"
	"io/fs"
	"time"

	"github.com/fzipp/gg/ggdict"
)

type packFile struct {
	stat fs.FileInfo
	r    io.Reader
}

func (f packFile) Stat() (fs.FileInfo, error)       { return f.stat, nil }
func (f packFile) Read(b []byte) (n int, err error) { return f.r.Read(b) }
func (f packFile) Close() error                     { return nil }

type fileDirEntry struct {
	name    string
	size    int64
	modTime time.Time
}

func (d fileDirEntry) Name() string               { return d.name }
func (d fileDirEntry) IsDir() bool                { return false }
func (d fileDirEntry) Type() fs.FileMode          { return 0 }
func (d fileDirEntry) Info() (fs.FileInfo, error) { return d, nil }
func (d fileDirEntry) Size() int64                { return d.size }
func (d fileDirEntry) Mode() fs.FileMode          { return d.Type() }
func (d fileDirEntry) ModTime() time.Time         { return d.modTime }
func (d fileDirEntry) Sys() interface{}           { return nil }

type directory map[string]entry

type entry struct {
	Offset int64
	Size   int64
}

const (
	keyFiles    = "files"
	keyFilename = "filename"
	keyOffset   = "offset"
	keySize     = "size"
)

func readDirectory(buf []byte) (directory, error) {
	directoryDict, err := ggdict.Unmarshal(buf)
	if err != nil {
		return nil, fmt.Errorf("could not read directory: %w", err)
	}
	return directoryFrom(directoryDict)
}

func directoryFrom(dict map[string]interface{}) (directory, error) {
	files, ok := dict[keyFiles].([]interface{})
	if !ok {
		return nil, fmt.Errorf("%q is not an array", keyFiles)
	}
	directory := make(directory, len(files))
	for _, fileEntry := range files {
		entryDict := fileEntry.(map[string]interface{})
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
		directory[filename] = entry{
			Offset: int64(offset),
			Size:   int64(size),
		}
	}
	return directory, nil
}
