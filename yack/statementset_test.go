// Copyright 2020 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package yack

import (
	"testing"

	"github.com/fzipp/gg/yack/stmt"
)

func TestStatementSet(t *testing.T) {
	statements := []Statement{
		&stmt.Pause{Seconds: 3},
		&stmt.ShutUp{},
		&stmt.ShutUp{},
	}
	set := newStatementSet()
	for _, s := range statements {
		if set.contains(s) {
			t.Fatalf("set contains %#v, but it should not", s)
		}
	}
	for _, s := range statements {
		set.add(s)
	}
	for _, s := range statements {
		if !set.contains(s) {
			t.Fatalf("set does not contain %#v after add(), but it should", s)
		}
	}
	set.clear()
	for _, s := range statements {
		if set.contains(s) {
			t.Fatalf("set contains %#v after clear(), but it should not", s)
		}
	}
}
