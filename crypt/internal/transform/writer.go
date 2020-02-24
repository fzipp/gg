// Copyright 2020 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package transform

import "io"

type writer struct {
	writer      io.Writer
	transformer Transformer
}

func (w *writer) Write(p []byte) (n int, err error) {
	dst := make([]byte, len(p))
	w.transformer.Transform(dst, p)
	return w.writer.Write(dst)
}

func NewWriter(w io.Writer, t Transformer) io.Writer {
	return &writer{
		writer:      w,
		transformer: t,
	}
}
