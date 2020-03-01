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

func NewRunner(s Scripting, t Talk, startActor string) *Runner {
	return &Runner{ctx: newContext(s, t, startActor)}
}

func (r *Runner) Start(d *Dialog) *Choices {
	return r.StartAt(d, labelStart)
}

func (r *Runner) StartAt(d *Dialog, label string) *Choices {
	r.ctx.init(d)
	return r.ctx.runLabel(label)
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
