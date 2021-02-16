// Copyright 2020 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package transform_test

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/fzipp/gg/crypt/internal/transform"
)

func TestReader(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"", ""},
		{"abcdefg", "ABCDEFG"},
		{"This is a test.", "THIS IS A TEST."},
	}
	for _, tt := range tests {
		r := transform.NewReader(strings.NewReader(tt.input), upperCaseTransformer{})
		buf, err := io.ReadAll(r)
		output := string(buf)
		if err != nil {
			t.Errorf("reading %q with transformer returned error: %s", tt.input, err)
			continue
		}
		if output != tt.want {
			t.Errorf("reading %q with transformer resulted in %q, want: %q", tt.input, output, tt.want)
		}
	}
}

type upperCaseTransformer struct{}

func (t upperCaseTransformer) Transform(dst, src []byte) {
	copy(dst, bytes.ToUpper(src))
}
