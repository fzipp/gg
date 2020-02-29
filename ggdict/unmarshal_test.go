// Copyright 2020 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ggdict_test

import (
	"testing"

	"github.com/fzipp/gg/ggdict"
)

func TestUnmarshalErrors(t *testing.T) {
	tests := []struct {
		data      []byte
		wantError string
	}{
		{[]byte{
			0x4, 0x3, 0x2, 0x1, // (invalid) format signature
		}, "invalid format signature: 0x1020304"},
		{[]byte{
			0x1, 0x2, 0x3, 0x4, // format signature
			0x1, 0x0, 0x0, 0x0, // always 1
			0xc, 0x0, 0x0, 0x0, // string offsets start offset (12)
			0x0, // (wrong) offsets type start marker
		}, "could not read offset index: unknown value type: 0"},
		{[]byte{
			0x1, 0x2, 0x3, 0x4, // format signature
			0x1, 0x0, 0x0, 0x0, // always 1
			0xc, 0x0, 0x0, 0x0, // string offsets start offset (12)
			0x1, // (wrong) offsets type start marker
		}, "read value is not an offset index"},
		{[]byte{
			0x1, 0x2, 0x3, 0x4, // format signature
			0x1, 0x0, 0x0, 0x0, // always 1
			0xd, 0x0, 0x0, 0x0, // string offsets start offset (13)
			0x1,                    // (wrong) dictionary type start marker
			0x7,                    // string offsets type start marker
			0xff, 0xff, 0xff, 0xff, // string offsets end marker
			0x8, // no strings
		}, "root is not a dictionary"},
		{[]byte{
			0x1, 0x2, 0x3, 0x4, // format signature
			0x1, 0x0, 0x0, 0x0, // always 1
			0x12, 0x0, 0x0, 0x0, // string offsets start offset (18)
			0x2,                // dictionary type start marker
			0x0, 0x0, 0x0, 0x0, // length of dictionary (0)
			0x0,                    // (wrong) dictionary end marker
			0x7,                    // string offsets type start marker
			0xff, 0xff, 0xff, 0xff, // string offsets end marker
			0x8, // no strings
		}, "could not read root: unterminated dictionary"},
		{[]byte{
			0x1, 0x2, 0x3, 0x4, // format signature
			0x1, 0x0, 0x0, 0x0, // always 1
			0x1c, 0x0, 0x0, 0x0, // string offsets start offset (28)
			0x2,                // dictionary type start marker
			0x1, 0x0, 0x0, 0x0, // length of dictionary (1)

			0x0, 0x0, 0x0, 0x0, // string table offset 0: "a"
			0x3,                // array type marker
			0x0, 0x0, 0x0, 0x0, // array length (0)
			0x0, // (wrong) array end marker

			0x2, // dictionary end marker
			0x7, // string offsets type start marker
			0x25, 0x0, 0x0, 0x0,
			0xff, 0xff, 0xff, 0xff, // string offsets end marker
			0x61, 0x0, // "a\x00"
			0x8,
		}, `could not read root: could not read dictionary value for key "a": unterminated array`},
		{[]byte{
			0x1, 0x2, 0x3, 0x4, // format signature
			0x1, 0x0, 0x0, 0x0, // always 1
			0x1d, 0x0, 0x0, 0x0, // string offsets start offset (29)
			0x2,                // dictionary type start marker
			0x1, 0x0, 0x0, 0x0, // length of dictionary (1)

			0x0, 0x0, 0x0, 0x0, // string table offset 0: "a"
			0x3,                // array type marker
			0x1, 0x0, 0x0, 0x0, // array length (1)
			0x0, // invalid array value
			0x3, // array end marker

			0x2, // dictionary end marker
			0x7, // string offsets type start marker
			0x26, 0x0, 0x0, 0x0,
			0xff, 0xff, 0xff, 0xff, // string offsets end marker
			0x61, 0x0, // "a\x00"
			0x8,
		}, `could not read root: could not read dictionary value for key "a": could not read array value: unknown value type: 0`},
	}
	for _, tt := range tests {
		_, err := ggdict.Unmarshal(tt.data)
		if err == nil {
			t.Errorf("expected error for unmarshalling of %#v, but no error returned", tt.data)
			continue
		}
		if err.Error() != tt.wantError {
			t.Errorf("error message for unmarshalling of %#v was: %q, want: %q", tt.data, err.Error(), tt.wantError)
		}
	}
}
