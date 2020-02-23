// Copyright 2020 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package yack

import (
	"errors"
	"testing"
	"text/scanner"
)

func TestErr(t *testing.T) {
	tests := []struct {
		list errorList
		want error
	}{
		{nil, nil},
		{errorList{}, nil},
		{
			errorList{errors.New("e")},
			errors.New("e"),
		},
		{
			errorList{errors.New("a"), errors.New("b")},
			errors.New("a (and 1 more errors)"),
		},
		{
			errorList{errors.New("x"), errors.New("y"), errors.New("z")},
			errors.New("x (and 2 more errors)"),
		},
	}
	for _, tt := range tests {
		err := tt.list.Err()
		if !errorsEqual(err, tt.want) {
			t.Errorf("error list error was %q, want: %q", err, tt.want)
		}
	}
}

func TestError(t *testing.T) {
	tests := []struct {
		list errorList
		want string
	}{
		{nil, "no errors"},
		{errorList{}, "no errors"},
		{
			errorList{errors.New("e")},
			"e",
		},
		{
			errorList{errors.New("a"), errors.New("b")},
			"a (and 1 more errors)",
		},
		{
			errorList{errors.New("x"), errors.New("y"), errors.New("z")},
			"x (and 2 more errors)",
		},
	}
	for _, tt := range tests {
		message := tt.list.Error()
		if message != tt.want {
			t.Errorf("error list error message was %q, want: %q", message, tt.want)
		}
	}
}

func TestNewError(t *testing.T) {
	tests := []struct {
		pos  scanner.Position
		msg  string
		want error
	}{
		{
			scanner.Position{
				Filename: "test",
				Offset:   0,
				Line:     1,
				Column:   1,
			}, "an error message",
			errors.New("test:1:1: an error message"),
		},
		{
			scanner.Position{
				Filename: "test.yack",
				Offset:   44,
				Line:     2,
				Column:   14,
			}, "parse error",
			errors.New("test.yack:2:14: parse error"),
		},
	}
	for _, tt := range tests {
		err := newError(tt.pos, tt.msg)
		if !errorsEqual(err, tt.want) {
			t.Errorf("new error was %q, want: %q", err, tt.want)
		}
	}
}

func errorsEqual(a, b error) bool {
	if a == nil && b != nil {
		return false
	}
	if a != nil && b == nil {
		return false
	}
	if a != nil && b != nil {
		return a.Error() == b.Error()
	}
	return true
}
