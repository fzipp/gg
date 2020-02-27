// Copyright 2020 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package yack_test

import (
	"testing"

	"github.com/fzipp/gg/yack"
	"github.com/fzipp/gg/yack/condition"
	"github.com/fzipp/gg/yack/stmt"
)

var testStatementsDialog = &yack.Dialog{
	Statements: []yack.ConditionalStatement{
		{Statement: stmt.ShutUp{}},
		{Statement: stmt.Say{Actor: "testactor", Text: "@12345"}},
		{Statement: stmt.Say{Actor: "testactor2", Text: "hello, world", OptionalGotoLabel: "done"}},
		{Statement: stmt.Pause{Seconds: 2.5}},
		{Statement: stmt.Execute{Code: "testFunc()"}},
		{Statement: stmt.Goto{Label: "main"}},
		{Statement: stmt.Choice{Index: 1, Text: "hello, world", GotoLabel: "greet"}},
		{Statement: stmt.Choice{Index: 1, Text: "lorem ipsum", GotoLabel: "more"}},
		{Statement: stmt.Choice{Index: 2, Text: "bye", GotoLabel: "done"}},
		{Statement: stmt.WaitFor{Actor: "testactor"}},
		{Statement: stmt.WaitWhile{CodeCondition: "g.test_var == NO"}},
		{Statement: stmt.Parrot{Enabled: false}},
		{Statement: stmt.Parrot{Enabled: true}},
		{Statement: stmt.Dialog{Actor: "testactor2"}},
		{Statement: stmt.Override{Label: "done"}},
		{Statement: stmt.AllowObjects{Allow: false}},
		{Statement: stmt.AllowObjects{Allow: true}},
		{Statement: stmt.Limit{N: 4}},
	},
	LabelIndex: map[string]int{
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
		{Statement: stmt.ShutUp{}},
		{Statement: stmt.ShutUp{}},
		{Statement: stmt.ShutUp{}},
		{Statement: stmt.ShutUp{}},
		{Statement: stmt.ShutUp{}},
		{Statement: stmt.ShutUp{}},
		{Statement: stmt.ShutUp{}},
	},
	LabelIndex: map[string]int{
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
		conditions yack.Conditions
		want       string
	}{
		{yack.Conditions{}, ""},
		{yack.Conditions{
			&condition.Actor{Actor: "testactor"},
		}, "[testactor]"},
		{yack.Conditions{
			&condition.Actor{Actor: "testactor2"},
			&condition.Once{},
			&condition.Code{Code: "g.test_var == YES"},
		}, "[testactor2] [once] [g.test_var == YES]"},
		{yack.Conditions{
			&condition.OnceEver{},
			&condition.Code{Code: "testFunc()"},
		}, "[onceever] [testFunc()]"},
		{yack.Conditions{
			&condition.TempOnce{},
		}, "[temponce]"},
		{yack.Conditions{
			&condition.ShowOnce{},
			&condition.Actor{Actor: "testactor"},
		}, "[showonce] [testactor]"},
		{yack.Conditions{
			&condition.Actor{Actor: "testactor"},
			&condition.ShowOnce{},
		}, "[testactor] [showonce]"},
	}
	for _, tt := range tests {
		if s := tt.conditions.String(); s != tt.want {
			t.Errorf("yack representation of conditions %#v was: %s, want: %s", tt.conditions, s, tt.want)
		}
	}
}
