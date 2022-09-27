// Copyright 2020 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ggdict_test

import (
	"reflect"
	"testing"

	"github.com/fzipp/gg/ggdict"
)

func TestMarshal(t *testing.T) {
	tests := []struct {
		dict map[string]any
		want []byte
	}{
		{nil, []byte{
			// format signature
			0x1, 0x2, 0x3, 0x4,
			// always 1
			0x1, 0x0, 0x0, 0x0,
			// string offsets start offset (18)
			0x12, 0x0, 0x0, 0x0,
			// dictionary type start marker
			0x2,
			// length of dictionary (0)
			0x0, 0x0, 0x0, 0x0,
			// no dictionary entries
			// dictionary end marker
			0x2,
			// string offsets type start marker
			0x7,
			// no string offsets
			// string offsets end marker
			0xff, 0xff, 0xff, 0xff,
			// string values start marker
			0x8,
			// no strings
		}},
		{map[string]any{
			"key_array": []any{"test", 14, 3.2},
			"key_dictionary": map[string]any{
				"key_a": 26,
				"key_b": 54.8,
				"key_c": "test",
			},
			"key_float64": 0.5,
			"key_float32": float32(0.5),
			"key_int":     4,
			"key_int32":   int32(5),
			"key_int64":   int64(6),
			"key_null":    nil,
			"key_string":  "test",
			"key_uint32":  uint32(7),
			"key_uint64":  uint64(8),
		}, []byte{
			// format signature
			0x1, 0x2, 0x3, 0x4,
			// always 1
			0x1, 0x0, 0x0, 0x0,
			// string offsets start offset (157)
			0x9d, 0x0, 0x0, 0x0,
			// dictionary type start marker
			0x2,
			// length of dictionary (11)
			0xb, 0x0, 0x0, 0x0,

			0x0, 0x0, 0x0, 0x0, // string table offset 0: "key_array"
			0x3,                // array type marker
			0x3, 0x0, 0x0, 0x0, // array length (3)
			0x4,                // [0] string type marker
			0x1, 0x0, 0x0, 0x0, // [0] string table offset 1: "test"
			0x5,                // [1] int type marker
			0x2, 0x0, 0x0, 0x0, // [1] string table offset 2: "14"
			0x6,                // [2] float type marker
			0x3, 0x0, 0x0, 0x0, // [2] string table offset 3: "3.2"
			0x3, // array end marker

			0x4, 0x0, 0x0, 0x0, // string table offset 4: "key_dictionary"
			0x2,                // dictionary type marker
			0x3, 0x0, 0x0, 0x0, // dictionary length (3)
			0x5, 0x0, 0x0, 0x0, // string table offset 5: "key_a"
			0x5,                // int type marker
			0x6, 0x0, 0x0, 0x0, // string table offset 6: "26"
			0x7, 0x0, 0x0, 0x0, // string table offset 7: "key_b"
			0x6,                // float type marker
			0x8, 0x0, 0x0, 0x0, // string table offset 8: "54.8"
			0x9, 0x0, 0x0, 0x0, // string table offset 9: "key_c"
			0x4,                // string type marker
			0x1, 0x0, 0x0, 0x0, // string table offset 1: "test"
			0x2, // dictionary end marker

			0xa, 0x0, 0x0, 0x0, // string table offset 10: "key_float32"
			0x6,                // float type marker
			0xb, 0x0, 0x0, 0x0, // string table offset 11: "0.5"

			0xc, 0x0, 0x0, 0x0, // string table offset 12: "key_float64"
			0x6,                // float type marker
			0xb, 0x0, 0x0, 0x0, // string table offset 11: "0.5"

			0xd, 0x0, 0x0, 0x0, // string table offset 13: "key_int"
			0x5,                // int type marker
			0xe, 0x0, 0x0, 0x0, // string table offset 14: "4"

			0xf, 0x0, 0x0, 0x0, // string table offset 15: "key_int32"
			0x5,                 // int type marker
			0x10, 0x0, 0x0, 0x0, // string table offset 16: "5"

			0x11, 0x0, 0x0, 0x0, // string table offset 17: "key_int64"
			0x5,                 // int type marker
			0x12, 0x0, 0x0, 0x0, // string table offset 18: "6"

			0x13, 0x0, 0x0, 0x0, // string table offset 19: "key_null"
			0x1, // null type marker

			0x14, 0x0, 0x0, 0x0, // string table offset 20: "key_string"
			0x4,                // string type marker
			0x1, 0x0, 0x0, 0x0, // string table offset 1: "test"

			0x15, 0x0, 0x0, 0x0, // string table offset 21: "key_uint32"
			0x5,                 // int type marker
			0x16, 0x0, 0x0, 0x0, // string table offset 22: "7"

			0x17, 0x0, 0x0, 0x0, // string table offset 23: "key_uint64"
			0x5,                 // int type marker
			0x18, 0x0, 0x0, 0x0, // string table offset 24: "8"

			0x2, // dictionary end marker

			// string offsets type start marker
			0x7,
			// string offsets
			0x7, 0x1, 0x0, 0x0,
			0x11, 0x1, 0x0, 0x0,
			0x16, 0x1, 0x0, 0x0,
			0x19, 0x1, 0x0, 0x0,
			0x1d, 0x1, 0x0, 0x0,
			0x2c, 0x1, 0x0, 0x0,
			0x32, 0x1, 0x0, 0x0,
			0x35, 0x1, 0x0, 0x0,
			0x3b, 0x1, 0x0, 0x0,
			0x40, 0x1, 0x0, 0x0,
			0x46, 0x1, 0x0, 0x0,
			0x52, 0x1, 0x0, 0x0,
			0x56, 0x1, 0x0, 0x0,
			0x62, 0x1, 0x0, 0x0,
			0x6a, 0x1, 0x0, 0x0,
			0x6c, 0x1, 0x0, 0x0,
			0x76, 0x1, 0x0, 0x0,
			0x78, 0x1, 0x0, 0x0,
			0x82, 0x1, 0x0, 0x0,
			0x84, 0x1, 0x0, 0x0,
			0x8d, 0x1, 0x0, 0x0,
			0x98, 0x1, 0x0, 0x0,
			0xa3, 0x1, 0x0, 0x0,
			0xa5, 0x1, 0x0, 0x0,
			0xb0, 0x1, 0x0, 0x0,
			// string offsets end marker
			0xff, 0xff, 0xff, 0xff,

			// string values start marker
			0x8,
			// "key_array\x00"
			0x6b, 0x65, 0x79, 0x5f, 0x61, 0x72, 0x72, 0x61, 0x79, 0x0,
			// "test\x00"
			0x74, 0x65, 0x73, 0x74, 0x0,
			// "14\x00"
			0x31, 0x34, 0x0,
			// "3.2\x00"
			0x33, 0x2e, 0x32, 0x0,
			// "key_dictionary\x00"
			0x6b, 0x65, 0x79, 0x5f, 0x64, 0x69, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x61, 0x72, 0x79, 0x0,
			// "key_a\x00"
			0x6b, 0x65, 0x79, 0x5f, 0x61, 0x0,
			// "26\x00"
			0x32, 0x36, 0x0,
			// "key_b\x00"
			0x6b, 0x65, 0x79, 0x5f, 0x62, 0x0,
			// "54.8\x00"
			0x35, 0x34, 0x2e, 0x38, 0x0,
			// "key_c\x00"
			0x6b, 0x65, 0x79, 0x5f, 0x63, 0x0,
			// "key_float32\x00"
			0x6b, 0x65, 0x79, 0x5f, 0x66, 0x6c, 0x6f, 0x61, 0x74, 0x33, 0x32, 0x0,
			// "0.5\x00"
			0x30, 0x2e, 0x35, 0x0,
			// "key_float64\x00"
			0x6b, 0x65, 0x79, 0x5f, 0x66, 0x6c, 0x6f, 0x61, 0x74, 0x36, 0x34, 0x0,
			// "key_int\x00"
			0x6b, 0x65, 0x79, 0x5f, 0x69, 0x6e, 0x74, 0x0,
			// "4\x00"
			0x34, 0x0,
			// "key_int32\x00"
			0x6b, 0x65, 0x79, 0x5f, 0x69, 0x6e, 0x74, 0x33, 0x32, 0x0,
			// "5\x00"
			0x35, 0x0,
			// "key_int64\x00"
			0x6b, 0x65, 0x79, 0x5f, 0x69, 0x6e, 0x74, 0x36, 0x34, 0x0,
			// "6\x00"
			0x36, 0x0,
			// "key_null\x00"
			0x6b, 0x65, 0x79, 0x5f, 0x6e, 0x75, 0x6c, 0x6c, 0x0,
			// "key_string\x00"
			0x6b, 0x65, 0x79, 0x5f, 0x73, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x0,
			// "key_uint32\x00"
			0x6b, 0x65, 0x79, 0x5f, 0x75, 0x69, 0x6e, 0x74, 0x33, 0x32, 0x0,
			// "7\x00"
			0x37, 0x0,
			// "key_uint64\x00"
			0x6b, 0x65, 0x79, 0x5f, 0x75, 0x69, 0x6e, 0x74, 0x36, 0x34, 0x0,
			// "8\x00"
			0x38, 0x0,
		}},
	}
	for _, tt := range tests {
		if data := ggdict.Marshal(tt.dict, false); !reflect.DeepEqual(data, tt.want) {
			t.Errorf("ggdict marshalling of %#v was:\n%#v, want:\n%#v", tt.dict, data, tt.want)
		}
	}
}
