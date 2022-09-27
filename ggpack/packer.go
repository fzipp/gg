// Copyright 2020 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ggpack

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/fzipp/gg/crypt/bnut"
	"github.com/fzipp/gg/crypt/xor"
	"github.com/fzipp/gg/ggdict"
)

type Packer struct {
	writer   io.WriteSeeker
	offset   int64
	xorKey   xor.Key
	files    []any
	finished bool
}

func NewPacker(w io.WriteSeeker) (*Packer, error) {
	n, err := w.Write(make([]byte, 8))
	if err != nil {
		return nil, err
	}
	return &Packer{writer: w, offset: int64(n), xorKey: xor.DefaultKey}, nil
}

// SetKey sets the key for XOR encryption, if a different key than the default
// key (xor.DefaultKey) should be used.
// The key should be set before any write operations.
func (p *Packer) SetKey(key xor.Key) {
	p.xorKey = key
}

func (p *Packer) WriteFiles(paths []string) error {
	for _, path := range paths {
		err := p.WriteFile(path)
		if err != nil {
			return fmt.Errorf("could not write '%s' to pack file: %w", path, err)
		}
	}
	return nil
}

func (p *Packer) WriteFile(path string) error {
	return p.WriteFileAs(filepath.Base(path), path)
}

func (p *Packer) WriteFileAs(filenameInPack, sourceFilePath string) error {
	file, err := os.Open(sourceFilePath)
	if err != nil {
		return fmt.Errorf("could not open file '%s': %w", sourceFilePath, err)
	}
	defer file.Close()
	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("could not obdtain file stats: %w", err)
	}
	return p.Write(filenameInPack, file, fileInfo.Size())
}

func (p *Packer) Write(filenameInPack string, r io.Reader, size int64) error {
	if p.finished {
		return errors.New("attempted write to already finished pack")
	}

	fileOffset := p.offset
	var w io.Writer = p.writer
	w = p.xorKey.EncodingWriter(w, size)
	if filepath.Ext(filenameInPack) == ".bnut" {
		w = bnut.EncodingWriter(w, size)
	}
	n, err := io.CopyN(w, r, size)
	p.offset += n
	if err != nil {
		return fmt.Errorf("could not copy file data to pack: %w", err)
	}

	p.files = append(p.files, map[string]any{
		"filename": filenameInPack,
		"offset":   int(fileOffset),
		"size":     int(size),
	})

	return nil
}

func (p *Packer) Finish() error {
	if p.finished {
		return errors.New("pack already finished")
	}

	directory := map[string]any{
		"files": p.files,
	}
	dirOffset := p.offset
	data := ggdict.Marshal(directory, p.xorKey.UsesShortKeyIndices())
	size := len(data)
	n, err := p.xorKey.EncodingWriter(p.writer, int64(size)).Write(data)
	p.offset += int64(n)
	if err != nil {
		return err
	}
	p.finished = true

	_, err = p.writer.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}
	byteOrder := binary.LittleEndian
	err = binary.Write(p.writer, byteOrder, uint32(dirOffset))
	if err != nil {
		return err
	}
	err = binary.Write(p.writer, byteOrder, uint32(size))
	if err != nil {
		return err
	}

	return nil
}
