// Copyright 2020 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package transform

import (
	"io"
)

type reader struct {
	reader      io.Reader
	transformer Transformer
}

func (r *reader) Read(p []byte) (n int, err error) {
	n, err = r.reader.Read(p)
	r.transformer.Transform(p[:n], p[:n])
	return n, err
}

func NewReader(r io.Reader, t Transformer) io.Reader {
	return &reader{
		reader:      r,
		transformer: t,
	}
}
