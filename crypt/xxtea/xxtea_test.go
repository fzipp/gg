// Copyright 2020 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xxtea_test

import (
	"reflect"
	"testing"

	"github.com/fzipp/gg/crypt/xxtea"
)

var key = xxtea.Key{
	0xAEA4EDF3,
	0xAFF8332A,
	0xB5A2DBB4,
	0x9B4BA022,
}

func TestEncrypt(t *testing.T) {
	tests := []struct {
		input []byte
		key   xxtea.Key
		want  []byte
	}{
		{
			[]byte("hello, world"), key,
			[]byte{
				0x54, 0xC3, 0xFB, 0xB8, 0xF5, 0xAA,
				0x3F, 0x3C, 0x5B, 0x91, 0xC3, 0x98,
			},
		},
		{
			[]byte("abcdefgh"), key,
			[]byte{0x9D, 0x5F, 0x1C, 0x05, 0xEB, 0x20, 0xB4, 0x4A},
		},
		{
			[]byte("abcdefg"), key,
			[]byte("abcdefg"),
		},
		{
			[]byte("abcdefghij"), key,
			[]byte{0x9D, 0x5F, 0x1C, 0x05, 0xEB, 0x20, 0xB4, 0x4A, 0x69, 0x6A},
		},
	}
	for _, tt := range tests {
		got := xxtea.Encrypt(tt.input, tt.key)
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("xxtea encryption of %q with key %X, got: %X, want: %X",
				tt.input, tt.key, got, tt.want)
		}
	}
}

func TestDecrypt(t *testing.T) {
	tests := []struct {
		input []byte
		key   xxtea.Key
		want  []byte
	}{
		{
			[]byte{
				0x54, 0xC3, 0xFB, 0xB8, 0xF5, 0xAA,
				0x3F, 0x3C, 0x5B, 0x91, 0xC3, 0x98,
			}, key,
			[]byte("hello, world"),
		},
		{
			[]byte{0x9D, 0x5F, 0x1C, 0x05, 0xEB, 0x20, 0xB4, 0x4A}, key,
			[]byte("abcdefgh"),
		},
		{
			[]byte("abcdefg"), key,
			[]byte("abcdefg"),
		},
		{
			[]byte{0x9D, 0x5F, 0x1C, 0x05, 0xEB, 0x20, 0xB4, 0x4A, 0x69, 0x6A}, key,
			[]byte("abcdefghij"),
		},
	}
	for _, tt := range tests {
		got := xxtea.Decrypt(tt.input, tt.key)
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("xxtea decryption of %q with key %X, got: %X, want: %X",
				tt.input, tt.key, got, tt.want)
		}
	}
}
