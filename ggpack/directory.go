// Copyright 2020 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ggpack

import (
	"fmt"

	"github.com/fzipp/gg/ggdict"
)

type DirectoryEntry struct {
	Filename string
	Size     int64
}

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
