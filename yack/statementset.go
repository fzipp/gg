// Copyright 2020 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package yack

import "github.com/fzipp/gg/yack/stmt"

type statementSet struct {
	statements map[stmt.Statement]struct{}
}

func newStatementSet() *statementSet {
	s := &statementSet{}
	s.clear()
	return s
}

func (set *statementSet) add(s stmt.Statement) {
	set.statements[s] = struct{}{}
}

func (set *statementSet) contains(s stmt.Statement) bool {
	_, ok := set.statements[s]
	return ok
}

func (set *statementSet) clear() {
	set.statements = make(map[stmt.Statement]struct{})
}
