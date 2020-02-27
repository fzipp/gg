// Copyright 2020 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package yack_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/fzipp/gg/yack"
	"github.com/fzipp/gg/yack/condition"
	"github.com/fzipp/gg/yack/stmt"
)

func TestParseLabels(t *testing.T) {
	tests := []struct {
		src  string
		want map[string]int
	}{
		{`
:label1
:label2
:label3`,
			map[string]int{
				"label1": 0,
				"label2": 0,
				"label3": 0,
			},
		},
		{`
:test_label1
-> test_label2

:test_label2
-> test_label3

:test_label3
-> test_label4

:test_label4`,
			map[string]int{
				"test_label1": 0,
				"test_label2": 1,
				"test_label3": 2,
				"test_label4": 3,
			},
		},
		{`
:init
!test()
testactor: "@20000"
-> exit

:main
1 "@20001" -> done
2 "@20002" -> done

:done
testactor: "@20003"`,
			map[string]int{
				"init": 0,
				"main": 3,
				"done": 5,
			},
		},
	}
	for _, tt := range tests {
		dlg, err := yack.Parse("test", strings.NewReader(tt.src))
		if err != nil {
			t.Errorf("parsing of dialog %q returned error: %s", tt.src, err)
			continue
		}
		if !reflect.DeepEqual(dlg.LabelIndex, tt.want) {
			t.Errorf("parsed label index for %q was: %#v, want: %#v", tt.src, dlg.LabelIndex, tt.want)
		}
	}
}

func TestParseStatement(t *testing.T) {
	tests := []struct {
		line string
		want yack.Statement
	}{
		{"-> exit", stmt.Goto{Label: "exit"}},
		{"-> main", stmt.Goto{Label: "main"}},
		{"-> test_label", stmt.Goto{Label: "test_label"}},
		{"-> testLabel", stmt.Goto{Label: "testLabel"}},

		{"!g.test_var <- NO", stmt.Execute{Code: "g.test_var <- NO"}},
		{"!g.test_var = YES", stmt.Execute{Code: "g.test_var = YES"}},
		{"!++g.test_var", stmt.Execute{Code: "++g.test_var"}},
		{"!test_var <- g.test_vars[0]", stmt.Execute{Code: "test_var <- g.test_vars[0]"}},
		{"!cameraFollow(currentActor)", stmt.Execute{Code: "cameraFollow(currentActor)"}},
		{"!actorTalkOffset(currentActor, 200, -500)", stmt.Execute{Code: "actorTalkOffset(currentActor, 200, -500)"}},
		{`!startActorIdle(testactor, 2.0, [ "a", "b", "c" ])`, stmt.Execute{Code: `startActorIdle(testactor, 2.0, [ "a", "b", "c" ])`}},
		{`!testFunc("\n\"\\")`, stmt.Execute{Code: `testFunc("\n\"\\")`}},

		{`1 "@12345" -> testLabel`, stmt.Choice{Index: 1, Text: "@12345", GotoLabel: "testLabel"}},
		{`2 "@12346" -> done`, stmt.Choice{Index: 2, Text: "@12346", GotoLabel: "done"}},
		{`3 $g.test_var -> testLabel`, stmt.Choice{Index: 3, Text: "$g.test_var", GotoLabel: "testLabel"}},
		{`4 $_testVar1 -> testLabel`, stmt.Choice{Index: 4, Text: "$_testVar1", GotoLabel: "testLabel"}},
		{`5 $Test.test_func_name(1) -> exit`, stmt.Choice{Index: 5, Text: "$Test.test_func_name(1)", GotoLabel: "exit"}},
		{`6 "$Test.test_func_name(2)" -> label1`, stmt.Choice{Index: 6, Text: "$Test.test_func_name(2)", GotoLabel: "label1"}},

		{"shutup", stmt.ShutUp{}},

		{"pause 0.5", stmt.Pause{Seconds: 0.5}},
		{"pause 0.432", stmt.Pause{Seconds: 0.432}},
		{"pause 1.0", stmt.Pause{Seconds: 1.0}},
		{"pause 4", stmt.Pause{Seconds: 4}},
		{"pause 8.0", stmt.Pause{Seconds: 8.0}},

		{"waitfor", stmt.WaitFor{Actor: ""}},
		{"waitfor testactor", stmt.WaitFor{Actor: "testactor"}},
		{"waitfor testactor2", stmt.WaitFor{Actor: "testactor2"}},
		{"waitfor currentActor", stmt.WaitFor{Actor: "currentActor"}},

		{"waitwhile Test.testMethod()", stmt.WaitWhile{CodeCondition: "Test.testMethod()"}},

		{"parrot NO", stmt.Parrot{Enabled: false}},
		{"parrot no", stmt.Parrot{Enabled: false}},
		{"parrot YES", stmt.Parrot{Enabled: true}},
		{"parrot yes", stmt.Parrot{Enabled: true}},

		{"dialog testactor", stmt.Dialog{Actor: "testactor"}},
		{"dialog testactor2", stmt.Dialog{Actor: "testactor2"}},

		{"override done", stmt.Override{Label: "done"}},
		{"override done2", stmt.Override{Label: "done2"}},

		{"allowobjects YES", stmt.AllowObjects{Allow: true}},
		{"allowobjects yes", stmt.AllowObjects{Allow: true}},
		{"allowobjects NO", stmt.AllowObjects{Allow: false}},
		{"allowobjects no", stmt.AllowObjects{Allow: false}},

		{"limit 3", stmt.Limit{N: 3}},
		{"limit 5", stmt.Limit{N: 5}},

		{`testactor: "@12345"`, stmt.Say{Actor: "testactor", Text: "@12345"}},
		{`testactor2: "@43057"`, stmt.Say{Actor: "testactor2", Text: "@43057"}},
		{`testactor: "This is a test."`, stmt.Say{Actor: "testactor", Text: "This is a test."}},
		{`testactor: "This is a test with escaped \"double quotes\"."`, stmt.Say{Actor: "testactor", Text: "This is a test with escaped \"double quotes\"."}},
		{`testactor: "This is a test with an escaped backslash: C:\\Program Files"`, stmt.Say{Actor: "testactor", Text: "This is a test with an escaped backslash: C:\\Program Files"}},
		{`testactor: "$g.test_var"`, stmt.Say{Actor: "testactor", Text: "$g.test_var"}},
		{`testactor: "^{test}"`, stmt.Say{Actor: "testactor", Text: "^{test}"}},
		{`testactor: "^{test_name}"`, stmt.Say{Actor: "testactor", Text: "^{test_name}"}},
		{`testactor: "@12345" -> done`, stmt.Say{Actor: "testactor", Text: "@12345", OptionalGotoLabel: "done"}},
		{`testactor: "This is a test." -> main`, stmt.Say{Actor: "testactor", Text: "This is a test.", OptionalGotoLabel: "main"}},
	}
	for _, tt := range tests {
		dlg, err := yack.Parse("test", strings.NewReader(tt.line))
		if err != nil {
			t.Errorf("parsing of statement %q returned error: %s", tt.line, err)
			continue
		}
		s := dlg.Statements[0].Statement
		if s != tt.want {
			t.Errorf("parsed statement %q was: %#v, want: %#v", tt.line, s, tt.want)
		}
	}
}

func TestParseCondition(t *testing.T) {
	tests := []struct {
		line string
		want yack.Condition
	}{
		{"! [once]", &condition.Once{}},
		{"! [showonce]", &condition.ShowOnce{}},
		{"! [onceever]", &condition.OnceEver{}},
		{"! [temponce]", &condition.TempOnce{}},

		{"! [testactor]", &condition.Actor{Actor: "testactor"}},
		{"! [testactor2]", &condition.Actor{Actor: "testactor2"}},

		{"! [_test_var]", &condition.Code{Code: "_test_var"}},
		{"! [test_var]", &condition.Code{Code: "test_var"}},
		{"! [g.test_var == 1]", &condition.Code{Code: "g.test_var == 1"}},
		{"! [g.test_var]", &condition.Code{Code: "g.test_var"}},
		{"! [test.testVar]", &condition.Code{Code: "test.testVar"}},
		{"! [!test.testVar && isTest()]", &condition.Code{Code: "!test.testVar && isTest()"}},
		{"! [test_var == YES]", &condition.Code{Code: "test_var == YES"}},
		{"! [random(1,5) == 1]", &condition.Code{Code: "random(1,5) == 1"}},
		{"! [(g.test_var == YES) && Test.testVar == YES]", &condition.Code{Code: "(g.test_var == YES) && Test.testVar == YES"}},
		{"! [!_test_var && !(testFunction1(test2) || testFunction2(test2))]", &condition.Code{Code: "!_test_var && !(testFunction1(test2) || testFunction2(test2))"}},
	}
	for _, tt := range tests {
		dlg, err := yack.Parse("test", strings.NewReader(tt.line))
		if err != nil {
			t.Errorf("parsing of condition %q returned error: %s", tt.line, err)
			continue
		}
		c := dlg.Statements[0].Conditions[0]
		if !reflect.DeepEqual(c, tt.want) {
			t.Errorf("parsed condition %q was: %#v, want: %#v", tt.line, c, tt.want)
		}
	}
}

func TestParseErrors(t *testing.T) {
	tests := []struct {
		src       string
		msgPrefix string
	}{
		{"parrot AYE", "test:1:11: invalid boolean literal: AYE"},
		{"pause five", "test:1:7: expected Float, found Ident five"},
		{"limit 0.5", "test:1:7: expected Int, found Float 0.5"},
		{"invalidcommand", "test:1:15: invalid command: invalidcommand"},
		{`testactor: "invalidstring`, `test:1:12: literal not terminated`},
		{"pause 2)", `test:1:8: expected "\n", found ")"`},
	}
	for _, tt := range tests {
		_, err := yack.Parse("test", strings.NewReader(tt.src))
		if err == nil {
			t.Errorf("no error on parsing %q, but expected one", tt.src)
			continue
		}
		if !strings.HasPrefix(err.Error(), tt.msgPrefix) {
			t.Errorf("error on parsing %q was %q, want %q[...]", tt.src, err.Error(), tt.msgPrefix)
		}
	}
}

func TestComments(t *testing.T) {
	tests := []struct {
		line string
		want yack.Statement
	}{
		{"!", stmt.Execute{Code: ""}},
		{"!;", stmt.Execute{Code: ""}},
		{"!;This is a comment", stmt.Execute{Code: ""}},
		{"!;; This is a comment", stmt.Execute{Code: ""}},
		{"!abc def", stmt.Execute{Code: "abc def"}},
		{"!abc def; This is a comment", stmt.Execute{Code: "abc def"}},
		{"!abc def ; This is a comment", stmt.Execute{Code: "abc def"}},
		{"!abc def;;This is a comment;;", stmt.Execute{Code: "abc def"}},
		{`!abc def ";This is not a comment"`, stmt.Execute{Code: `abc def ";This is not a comment"`}},
		{`!abc def ";This is not a comment"; This is a comment`, stmt.Execute{Code: `abc def ";This is not a comment"`}},
		{`!abc def ";\";This\" is not a comment"; This is a comment`, stmt.Execute{Code: `abc def ";\";This\" is not a comment"`}},
	}
	for _, tt := range tests {
		dlg, err := yack.Parse("test", strings.NewReader(tt.line))
		if err != nil {
			t.Errorf("parsing of statement with comment %q returned error: %s", tt.line, err)
			continue
		}
		s := dlg.Statements[0].Statement
		if s != tt.want {
			t.Errorf("parsed statement with comment %q was: %#v, want: %#v", tt.line, s, tt.want)
		}
	}
}

func TestLoad(t *testing.T) {
	path := "testdata/load.yack"
	want := &yack.Dialog{
		Statements: []yack.ConditionalStatement{
			{Statement: stmt.Execute{Code: "g.test_var <- YES"}},
			{Statement: stmt.Execute{Code: `testFunc("test")`}},
			{Statement: stmt.Goto{Label: "main"}},
			{Statement: stmt.Execute{Code: "testFunc()"}},
			{
				Statement:  stmt.Choice{Index: 1, Text: "@30001", GotoLabel: "test_label1"},
				Conditions: yack.Conditions{&condition.Once{}},
			},
			{
				Statement: stmt.Choice{Index: 2, Text: "@30002", GotoLabel: "test_label2"},
				Conditions: yack.Conditions{
					&condition.ShowOnce{},
					&condition.Code{Code: "g.test_var == NO"},
				},
			},
			{
				Statement:  stmt.Choice{Index: 2, Text: "@30003", GotoLabel: "test_label3"},
				Conditions: yack.Conditions{&condition.Code{Code: "g.test_var == YES"}},
			},
			{Statement: stmt.Choice{Index: 3, Text: "@30004", GotoLabel: "done"}},
			{Statement: stmt.Goto{Label: "exit"}},
			{Statement: stmt.Say{Actor: "testactor2", Text: "@40001", OptionalGotoLabel: ""}},
			{Statement: stmt.Say{Actor: "testactor", Text: "@40002", OptionalGotoLabel: ""}},
			{Statement: stmt.Goto{Label: "main"}},
			{Statement: stmt.Say{Actor: "testactor2", Text: "@40003", OptionalGotoLabel: ""}},
			{Statement: stmt.Goto{Label: "main"}},
			{
				Statement:  stmt.Say{Actor: "testactor2", Text: "@40004", OptionalGotoLabel: ""},
				Conditions: yack.Conditions{&condition.Code{Code: "test_var"}},
			},
			{
				Statement:  stmt.Say{Actor: "testactor2", Text: "@40004", OptionalGotoLabel: ""},
				Conditions: yack.Conditions{&condition.Code{Code: "!test_var"}},
			},
			{Statement: stmt.Goto{Label: "main"}},
			{Statement: stmt.Say{Actor: "testactor", Text: "@40005", OptionalGotoLabel: ""}},
		},
		LabelIndex: map[string]int{
			"init":        0,
			"start":       2,
			"main":        3,
			"test_label1": 9,
			"test_label2": 12,
			"test_label3": 14,
			"done":        17,
		},
	}
	dialog, err := yack.Load(path)
	if err != nil {
		t.Fatalf("loading of yack file %q returned error: %s", path, err)
	}
	if !reflect.DeepEqual(dialog, want) {
		t.Fatalf("parsed dialog from file %q was:\n%#v, want:\n%#v", path, dialog, want)
	}
}

func TestLoadError(t *testing.T) {
	tests := []struct {
		path string
	}{
		{"testdata/syntaxerror.yack"},
		{"testdata/doesnotexist.yack"},
	}
	for _, tt := range tests {
		_, err := yack.Load(tt.path)
		if err == nil {
			t.Fatalf("no error on parsing file %q, but expected one", tt.path)
		}
	}
}
