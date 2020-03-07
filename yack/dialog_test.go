// Copyright 2020 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package yack_test

import (
	"testing"

	"github.com/fzipp/gg/yack"
	"github.com/fzipp/gg/yack/cond"
	"github.com/fzipp/gg/yack/stmt"
)

var testStatementsDialog = &yack.Dialog{
	Statements: []yack.ConditionalStatement{
		{Statement: &stmt.ShutUp{}},
		{Statement: &stmt.Say{Actor: "testactor", Text: "@12345"}},
		{Statement: &stmt.Say{Actor: "testactor2", Text: "hello, world", OptionalGotoLabel: "done"}},
		{Statement: &stmt.Pause{Seconds: 2.5}},
		{Statement: &stmt.Execute{Code: "testFunc()"}},
		{Statement: &stmt.Goto{Label: "main"}},
		{Statement: &stmt.Choice{Index: 1, Text: "hello, world", GotoLabel: "greet"}},
		{Statement: &stmt.Choice{Index: 1, Text: "lorem ipsum", GotoLabel: "more"}},
		{Statement: &stmt.Choice{Index: 2, Text: "bye", GotoLabel: "done"}},
		{Statement: &stmt.WaitFor{Actor: "testactor"}},
		{Statement: &stmt.WaitWhile{CodeCondition: "g.test_var == NO"}},
		{Statement: &stmt.Parrot{Enabled: false}},
		{Statement: &stmt.Parrot{Enabled: true}},
		{Statement: &stmt.Dialog{Actor: "testactor2"}},
		{Statement: &stmt.Override{Label: "done"}},
		{Statement: &stmt.AllowObjects{Allow: false}},
		{Statement: &stmt.AllowObjects{Allow: true}},
		{Statement: &stmt.Limit{N: 4}},
	},
	Labels: map[string]int{
		"start": 0,
	},
}

func TestDialogStringStatements(t *testing.T) {
	tests := []struct {
		dialog *yack.Dialog
		want   string
	}{
		{testStatementsDialog, `
:start
shutup
testactor: "@12345"
testactor2: "hello, world" -> done
pause 2.5
!testFunc()
-> main
1 "hello, world" -> greet
1 "lorem ipsum" -> more
2 "bye" -> done
waitfor testactor
waitwhile g.test_var == NO
parrot NO
parrot YES
dialog testactor2
override done
allowobjects NO
allowobjects YES
limit 4
`},
	}
	for _, tt := range tests {
		if s := tt.dialog.String(); s != tt.want {
			t.Errorf("yack representation of dialog %#v was: %s, want: %s", tt.dialog, s, tt.want)
		}
	}
}

var testLabelsDialog = &yack.Dialog{
	Statements: []yack.ConditionalStatement{
		{Statement: &stmt.ShutUp{}},
		{Statement: &stmt.ShutUp{}},
		{Statement: &stmt.ShutUp{}},
		{Statement: &stmt.ShutUp{}},
		{Statement: &stmt.ShutUp{}},
		{Statement: &stmt.ShutUp{}},
		{Statement: &stmt.ShutUp{}},
	},
	Labels: map[string]int{
		"init":   0,
		"start":  2,
		"main":   5,
		"topic1": 5,
		"topic2": 6,
		"done":   7,
	},
}

func TestDialogStringLabels(t *testing.T) {
	tests := []struct {
		dialog *yack.Dialog
		want   string
	}{
		{testLabelsDialog, `
:init
shutup
shutup

:start
shutup
shutup
shutup

:main
:topic1
shutup

:topic2
shutup

:done
`},
	}
	for _, tt := range tests {
		if s := tt.dialog.String(); s != tt.want {
			t.Errorf("yack representation of dialog %#v was: %s, want: %s", tt.dialog, s, tt.want)
		}
	}
}

func TestConditionsString(t *testing.T) {
	tests := []struct {
		conditions []cond.Condition
		want       string
	}{
		{[]cond.Condition{}, "shutup"},
		{[]cond.Condition{
			&cond.Actor{Actor: "testactor"},
		}, "shutup [testactor]"},
		{[]cond.Condition{
			&cond.Actor{Actor: "testactor2"},
			&cond.Once{},
			&cond.Code{Code: "g.test_var == YES"},
		}, "shutup [testactor2] [once] [g.test_var == YES]"},
		{[]cond.Condition{
			&cond.OnceEver{},
			&cond.Code{Code: "testFunc()"},
		}, "shutup [onceever] [testFunc()]"},
		{[]cond.Condition{
			&cond.TempOnce{},
		}, "shutup [temponce]"},
		{[]cond.Condition{
			&cond.ShowOnce{},
			&cond.Actor{Actor: "testactor"},
		}, "shutup [showonce] [testactor]"},
		{[]cond.Condition{
			&cond.Actor{Actor: "testactor"},
			&cond.ShowOnce{},
		}, "shutup [testactor] [showonce]"},
	}
	for _, tt := range tests {
		condStmt := yack.ConditionalStatement{
			Statement:  &stmt.ShutUp{},
			Conditions: tt.conditions,
		}
		if s := condStmt.String(); s != tt.want {
			t.Errorf("yack representation of conditional statement %#v was: %s, want: %s", condStmt, s, tt.want)
		}
	}
}
