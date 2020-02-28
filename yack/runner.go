// Copyright 2020 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package yack

import (
	"github.com/fzipp/gg/yack/stmt"
)

const (
	labelInit  = "init"
	labelStart = "start"
	labelExit  = "exit"
)

type Runner struct {
	ctx *context
}

func NewRunner(d *Dialog, s Scripting, t Talk, startActor string) *Runner {
	return &Runner{ctx: newContext(d, s, t, startActor)}
}

func (r *Runner) Init() {
	_ = r.StartAt(labelInit)
}

func (r *Runner) Start() *Choices {
	return r.StartAt(labelStart)
}

func (r *Runner) StartAt(label string) *Choices {
	r.ctx.Goto(label)
	return r.ctx.run()
}

type Talk interface {
	Say(actor, text string)
}

type Choices struct {
	Actor   string
	Options []*ChoiceOption
	ctx     *context
}

func (c *Choices) Choose(opt *ChoiceOption) *Choices {
	return c.ctx.choose(opt)
}

type ChoiceOption struct {
	Text   string
	choice *stmt.Choice
}
