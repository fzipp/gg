// Copyright 2020 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package yack

type statementSet struct {
	statements map[Statement]struct{}
}

func newStatementSet() *statementSet {
	s := &statementSet{}
	s.clear()
	return s
}

func (s *statementSet) add(stmt Statement) {
	s.statements[stmt] = struct{}{}
}

func (s *statementSet) contains(stmt Statement) bool {
	_, ok := s.statements[stmt]
	return ok
}

func (s *statementSet) clear() {
	s.statements = make(map[Statement]struct{})
}
