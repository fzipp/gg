// Copyright 2020 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cond_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/fzipp/gg/yack/cond"
)

func TestString(t *testing.T) {
	tests := []struct {
		cond cond.Condition
		want string
	}{
		{&cond.Once{}, "[once]"},
		{&cond.ShowOnce{}, "[showonce]"},
		{&cond.OnceEver{}, "[onceever]"},
		{&cond.ShowOnceEver{}, "[showonceever]"},
		{&cond.TempOnce{}, "[temponce]"},
		{&cond.Actor{Actor: "testactor"}, "[testactor]"},
		{&cond.Code{Code: "g.test_var == NO"}, "[g.test_var == NO]"},
		{&cond.Code{Code: "testFunc()"}, "[testFunc()]"},
	}

	for _, tt := range tests {
		if s := tt.cond.String(); s != tt.want {
			t.Errorf("yack representation of condition %#v was: %q, want: %q", tt.cond, s, tt.want)
		}
	}
}

func TestIsFulfilled(t *testing.T) {
	conditions := []cond.Condition{
		&cond.Once{},
		&cond.ShowOnce{},
		&cond.OnceEver{},
		&cond.ShowOnceEver{},
		&cond.TempOnce{},
		&cond.Actor{Actor: "testactor"},
		&cond.Code{Code: "g.test_var == NO"},
		&cond.Code{Code: "testFunc()"},
	}
	wantCalls := `IsOnce()
IsShowOnce()
IsOnceEver()
IsShowOnceEver()
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

func (ctx *tracingTestContext) IsShowOnceEver() bool {
	ctx.callTrace.WriteString("IsShowOnceEver()\n")
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
