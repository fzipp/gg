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
	if s == nil {
		s = noScripting{}
	}
	return &Runner{ctx: newContext(s, t, startActor)}
}

// Start begins the given dialog at the "start" label and returns dialog
// choices when available.
// If the Options slice of the returned Choices object is empty the end
// of the dialog has been reached.
func (r *Runner) Start(d *Dialog) *Choices {
	return r.StartAt(d, labelStart)
}

// StartAt begins the given dialog at the specified label and returns choices.
// If the Options slice of the returned Choices object is empty the end
// of the dialog has been reached.
func (r *Runner) StartAt(d *Dialog, label string) *Choices {
	r.ctx.init(d)
	return r.ctx.runLabel(label)
}

type Talk interface {
	// Say makes the actor say a text.
	Say(actor, text string)
	// ShutUp makes all actors stop talking.
	ShutUp()
}

type Choices struct {
	Actor   string
	Options []*ChoiceOption
	ctx     *context
}

func (c *Choices) Choose(index int) *Choices {
	return c.ctx.choose(c.Options[index])
}

type ChoiceOption struct {
	Text   string
	choice *stmt.Choice
}

type noScripting struct{}

func (s noScripting) Eval(code string) (result any, err error) {
	// do nothing, always return true
	return true, nil
}
