// Copyright 2020 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bnut

import (
	"io"

	"github.com/fzipp/gg/crypt/internal/transform"
)

func DecodingReader(r io.Reader, expectedSize int64) io.Reader {
	return transform.NewReader(r, newTransformer(expectedSize))
}
