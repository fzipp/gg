// Copyright 2020 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stmt_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/fzipp/gg/yack/stmt"
)

func TestString(t *testing.T) {
	tests := []struct {
		statement stmt.Statement
		want      string
	}{
		{&stmt.ShutUp{}, "shutup"},
		{&stmt.Say{Actor: "testactor", Text: "hello, world"}, `testactor: "hello, world"`},
		{&stmt.Say{Actor: "testactor", Text: "hello, world", OptionalGotoLabel: "done"}, `testactor: "hello, world" -> done`},
		{&stmt.Pause{Seconds: 0.3}, "pause 0.3"},
		{&stmt.Parrot{Enabled: true}, "parrot YES"},
		{&stmt.Parrot{Enabled: false}, "parrot NO"},
		{&stmt.WaitFor{Actor: "testactor"}, "waitfor testactor"},
		{&stmt.WaitFor{}, "waitfor"},
		{&stmt.WaitWhile{CodeCondition: "g.test_var < 5"}, "waitwhile g.test_var < 5"},
		{&stmt.Dialog{Actor: "testactor2"}, "dialog testactor2"},
		{&stmt.Override{Label: "test_label"}, "override test_label"},
		{&stmt.AllowObjects{Allow: true}, "allowobjects YES"},
		{&stmt.AllowObjects{Allow: false}, "allowobjects NO"},
		{&stmt.Limit{N: 5}, "limit 5"},
		{&stmt.Goto{Label: "test_label2"}, "-> test_label2"},
		{&stmt.Execute{Code: "testFunction(2, 3)"}, "!testFunction(2, 3)"},
		{&stmt.Choice{Index: 1, Text: "hello, world", GotoLabel: "test_label3"}, `1 "hello, world" -> test_label3`},
		{&stmt.Choice{Index: 2, Text: "bye", GotoLabel: "done"}, `2 "bye" -> done`},
	}

	for _, tt := range tests {
		if s := tt.statement.String(); s != tt.want {
			t.Errorf("yack representation of statement %#v was: %q, want: %q", tt.statement, s, tt.want)
		}
	}
}

func TestExecute(t *testing.T) {
	statements := []stmt.Statement{
		&stmt.ShutUp{},
		&stmt.Say{Actor: "testactor", Text: "hello, world"},
		&stmt.Say{Actor: "testactor", Text: "hello, world", OptionalGotoLabel: "done"},
		&stmt.Pause{Seconds: 0.3},
		&stmt.Parrot{Enabled: true},
		&stmt.Parrot{Enabled: false},
		&stmt.WaitFor{Actor: "testactor"},
		&stmt.WaitFor{},
		&stmt.WaitWhile{CodeCondition: "g.test_var < 5"},
		&stmt.Dialog{Actor: "testactor2"},
		&stmt.Override{Label: "test_label"},
		&stmt.AllowObjects{Allow: true},
		&stmt.AllowObjects{Allow: false},
		&stmt.Limit{N: 5},
		&stmt.Goto{Label: "test_label2"},
		&stmt.Execute{Code: "testFunction(2, 3)"},
		&stmt.Choice{Index: 1, Text: "hello, world", GotoLabel: "test_label3"},
		&stmt.Choice{Index: 2, Text: "bye", GotoLabel: "done"},
	}
	wantCalls := `ShutUp()
Say("testactor", "hello, world")
Say("testactor", "hello, world")
Goto("done")
Pause(0.3)
Parrot(true)
Parrot(false)
WaitFor("testactor")
WaitFor("")
WaitWhile("g.test_var < 5")
Dialog("testactor2")
Override("test_label")
AllowObjects(true)
AllowObjects(false)
Limit(5)
Goto("test_label2")
Execute("testFunction(2, 3)")
Choice(1, "hello, world", "test_label3")
Choice(2, "bye", "done")
`

	ctx := &tracingTestContext{}
	for _, statement := range statements {
		statement.Execute(ctx)
	}
	if calls := ctx.callTrace.String(); calls != wantCalls {
		t.Errorf("call trace for statements was:\n%s, want:\n%s", calls, wantCalls)
	}
}

type tracingTestContext struct {
	callTrace strings.Builder
}

func (ctx *tracingTestContext) ShutUp() {
	ctx.callTrace.WriteString("ShutUp()\n")
}

func (ctx *tracingTestContext) Say(actor, text string) {
	ctx.callTrace.WriteString(fmt.Sprintf("Say(%q, %q)\n", actor, text))
}

func (ctx *tracingTestContext) Pause(seconds float64) {
	ctx.callTrace.WriteString(fmt.Sprintf("Pause(%v)\n", seconds))
}

func (ctx *tracingTestContext) Parrot(enabled bool) {
	ctx.callTrace.WriteString(fmt.Sprintf("Parrot(%t)\n", enabled))
}

func (ctx *tracingTestContext) WaitFor(actor string) {
	ctx.callTrace.WriteString(fmt.Sprintf("WaitFor(%q)\n", actor))
}

func (ctx *tracingTestContext) WaitWhile(codeCondition string) {
	ctx.callTrace.WriteString(fmt.Sprintf("WaitWhile(%q)\n", codeCondition))
}

func (ctx *tracingTestContext) Dialog(actor string) {
	ctx.callTrace.WriteString(fmt.Sprintf("Dialog(%q)\n", actor))
}

func (ctx *tracingTestContext) Override(label string) {
	ctx.callTrace.WriteString(fmt.Sprintf("Override(%q)\n", label))
}

func (ctx *tracingTestContext) AllowObjects(allow bool) {
	ctx.callTrace.WriteString(fmt.Sprintf("AllowObjects(%t)\n", allow))
}

func (ctx *tracingTestContext) Limit(n int) {
	ctx.callTrace.WriteString(fmt.Sprintf("Limit(%d)\n", n))
}

func (ctx *tracingTestContext) Goto(label string) {
	ctx.callTrace.WriteString(fmt.Sprintf("Goto(%q)\n", label))
}

func (ctx *tracingTestContext) Execute(code string) {
	ctx.callTrace.WriteString(fmt.Sprintf("Execute(%q)\n", code))
}

func (ctx *tracingTestContext) Choice(index int, text, gotoLabel string) {
	ctx.callTrace.WriteString(fmt.Sprintf("Choice(%d, %q, %q)\n", index, text, gotoLabel))
}
