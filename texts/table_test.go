// Copyright 2020 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package texts_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/fzipp/gg/texts"
)

const textTableTSV = `text_id	en
10001	hello
10002	world
10003	This is a test
`

func TestResolveTextsString(t *testing.T) {
	textTable, err := texts.From(strings.NewReader(textTableTSV))
	if err != nil {
		t.Fatalf("could not read text table: %s", err)
	}
	tests := []struct {
		text string
		want string
	}{
		{"@10001", "hello"},
		{"@10001, @10002", "hello, world"},
		{"@10001, @10002 @", "hello, world @"},
		{"@10003: @10001, @10002", "This is a test: hello, world"},
		{"@abc @def@", "@abc @def@"},
	}
	for _, tt := range tests {
		resolved, err := textTable.ResolveTextsString(tt.text)
		if err != nil {
			t.Errorf("could not writeResolved text string %q, got error: %s", tt.text, err)
		}
		if resolved != tt.want {
			t.Errorf("resolved text for %q was %q, want %q", tt.text, resolved, tt.want)
		}
	}
}

func TestFromFile(t *testing.T) {
	tests := []struct {
		path string
		want texts.Table
	}{
		{"testdata/TestTexts.tsv", texts.Table{
			12345: "hello, world",
			20001: "This is a test.",
			20003: "More text...",
			20002: `A text with \"quotes\"`,
		}},
	}
	for _, tt := range tests {
		textTable, err := texts.FromFile(tt.path)
		if err != nil {
			t.Errorf("could not load text table from file %q: %s", tt.path, err)
			continue
		}
		if !reflect.DeepEqual(textTable, tt.want) {
			t.Errorf("text table loaded from file %q was: %#v, want: %#v", tt.path, textTable, tt.want)
		}
	}
}
