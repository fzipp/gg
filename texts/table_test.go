// Copyright 2020 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package texts_test

import (
	"strings"
	"testing"

	"github.com/fzipp/gg/texts"
)

const textTableTSV = `text_id	en
10001	hello
10002	world
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
