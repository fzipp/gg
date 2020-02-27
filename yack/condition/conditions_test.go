// Copyright 2020 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package condition_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/fzipp/gg/yack"
	"github.com/fzipp/gg/yack/condition"
)

func TestString(t *testing.T) {
	tests := []struct {
		cond yack.Condition
		want string
	}{
		{&condition.Once{}, "[once]"},
		{&condition.ShowOnce{}, "[showonce]"},
		{&condition.OnceEver{}, "[onceever]"},
		{&condition.TempOnce{}, "[temponce]"},
		{&condition.Actor{Actor: "testactor"}, "[testactor]"},
		{&condition.Code{Code: "g.test_var == NO"}, "[g.test_var == NO]"},
		{&condition.Code{Code: "testFunc()"}, "[testFunc()]"},
	}

	for _, tt := range tests {
		if s := tt.cond.String(); s != tt.want {
			t.Errorf("yack representation of condition %#v was: %q, want: %q", tt.cond, s, tt.want)
		}
	}
}

func TestIsFulfilled(t *testing.T) {
	conditions := yack.Conditions{
		&condition.Once{},
		&condition.ShowOnce{},
		&condition.OnceEver{},
		&condition.TempOnce{},
		&condition.Actor{Actor: "testactor"},
		&condition.Code{Code: "g.test_var == NO"},
		&condition.Code{Code: "testFunc()"},
	}
	wantCalls := `IsOnce()
IsShowOnce()
IsOnceEver()
IsTempOnce()
IsActor("testactor")
IsCodeTrue("g.test_var == NO")
IsCodeTrue("testFunc()")
`

	ctx := &tracingTestContext{ret: true}
	for _, cond := range conditions {
		if f := cond.IsFulfilled(ctx); f != ctx.ret {
			t.Errorf("condition %#v fulfilled? was: %t, want: %t", cond, f, ctx.ret)
		}
	}
	if calls := ctx.callTrace.String(); calls != wantCalls {
		t.Errorf("call trace for conditions was:\n%s, want:\n%s", calls, wantCalls)
	}
}

type tracingTestContext struct {
	callTrace strings.Builder
	ret       bool
}

func (ctx *tracingTestContext) IsOnce() bool {
	ctx.callTrace.WriteString("IsOnce()\n")
	return ctx.ret
}

func (ctx *tracingTestContext) IsShowOnce() bool {
	ctx.callTrace.WriteString("IsShowOnce()\n")
	return ctx.ret
}

func (ctx *tracingTestContext) IsOnceEver() bool {
	ctx.callTrace.WriteString("IsOnceEver()\n")
	return ctx.ret
}

func (ctx *tracingTestContext) IsTempOnce() bool {
	ctx.callTrace.WriteString("IsTempOnce()\n")
	return ctx.ret
}

func (ctx *tracingTestContext) IsCodeTrue(code string) bool {
	ctx.callTrace.WriteString(fmt.Sprintf("IsCodeTrue(%q)\n", code))
	return ctx.ret
}

func (ctx *tracingTestContext) IsActor(actor string) bool {
	ctx.callTrace.WriteString(fmt.Sprintf("IsActor(%q)\n", actor))
	return ctx.ret
}
