// Copyright 2020 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ggdict_test

import (
	"reflect"
	"testing"

	"github.com/fzipp/gg/ggdict"
)

func TestRoundTrip(t *testing.T) {
	dict := map[string]any{
		"name":    "Test",
		"count":   4,
		"numbers": []any{0.5, 3, 2.6, 1.4},
		"subobject": map[string]any{
			"title": "Test 2",
			"id":    0,
		},
		"nothing": nil,
	}
	data := ggdict.Marshal(dict)
	newDict, err := ggdict.Unmarshal(data)
	if err != nil {
		t.Errorf("Unmarshal returned an error: %s", err)
		return
	}
	if !reflect.DeepEqual(dict, newDict) {
		t.Errorf("Marshal/unmarshal round trip resulted in\n%#v, want:\n%#v", newDict, dict)
	}
}
