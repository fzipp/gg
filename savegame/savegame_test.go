// Copyright 2020 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package savegame

import (
	"reflect"
	"testing"
)

func TestZeroPad(t *testing.T) {
	tests := []struct {
		data   []byte
		minLen int
		want   []byte
	}{
		{[]byte(""), 5, []byte("\x00\x00\x00\x00\x00")},
		{[]byte("a"), 5, []byte("a\x00\x00\x00\x00")},
		{[]byte("ab"), 5, []byte("ab\x00\x00\x00")},
		{[]byte("abc"), 5, []byte("abc\x00\x00")},
		{[]byte("hello"), 10, []byte("hello\x00\x00\x00\x00\x00")},
		{[]byte("hello, world"), 10, []byte("hello, world")},
	}
	for _, tt := range tests {
		got := zeroPad(tt.data, tt.minLen)
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("zero padding for %q with minimum length %d - got: %q, want: %q",
				tt.data, tt.minLen, got, tt.want)
		}
	}
}

func TestIsChecksumOk(t *testing.T) {
	tests := []struct {
		data []byte
		want bool
	}{
		{[]byte("\x63\x34\x58\x06\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00"), true},
		{[]byte("a\xc4\x34\x58\x06\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00"), true},
		{[]byte("ab\x26\x35\x58\x06\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00"), true},
		{[]byte("abc\x89\x35\x58\x06\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00"), true},
		{[]byte("hello, world\xeb\x38\x58\x06\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00"), true},

		{[]byte("\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00"), false},
		{[]byte("abc\x63\x35\x58\x06\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00"), false},
		{[]byte("hello, world\xeb\x38\x58\x07\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00"), false},
	}
	for _, tt := range tests {
		got := isChecksumOk(tt.data)
		if got != tt.want {
			t.Errorf("checksum ok for %q? got: %v, want: %v",
				tt.data, got, tt.want)
		}
	}
}

func TestChecksum(t *testing.T) {
	tests := []struct {
		data []byte
		want uint32
	}{
		{[]byte(""), 0x06583463},
		{[]byte("a"), 0x065834c4},
		{[]byte("ab"), 0x06583526},
		{[]byte("abc"), 0x06583589},
		{[]byte("hello, world"), 0x065838eb},
	}
	for _, tt := range tests {
		got := checksum(tt.data)
		if got != tt.want {
			t.Errorf("checksum for %q - got: %#08x, want: %#08x",
				tt.data, got, tt.want)
		}
	}
}
