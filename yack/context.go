// Copyright 2020 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package yack

import (
	"sort"
	"time"

	"github.com/fzipp/gg/yack/stmt"
)

type context struct {
	dialog    *Dialog
	scripting Scripting
	talk      Talk

	stmtCounter int

	currentActor   string
	parrot         bool
	objectsAllowed bool
	limit          int

	choices map[int]*stmt.Choice

	executed map[Statement]bool
	shown    map[Statement]bool
}

func newContext(d *Dialog, s Scripting, t Talk, startActor string) *context {
	return &context{
		dialog:       d,
		scripting:    s,
		talk:         t,
		currentActor: startActor,
		parrot:       true,
		limit:        6,
		choices:      make(map[int]*stmt.Choice),
		executed:     make(map[Statement]bool),
		shown:        make(map[Statement]bool),
	}
}

func (ctx *context) Say(actor, text string) {
	ctx.talk.Say(actor, text)
}

func (ctx *context) Pause(seconds float64) {
	time.Sleep(time.Duration(seconds * float64(time.Second)))
}

func (ctx *context) Parrot(enabled bool) {
	ctx.parrot = enabled
}

func (ctx *context) Dialog(actor string) {
	ctx.currentActor = actor
}

func (ctx *context) AllowObjects(allow bool) {
	ctx.objectsAllowed = allow
}

func (ctx *context) Limit(n int) {
	ctx.limit = n
}

func (ctx *context) Goto(label string) {
	if label == labelExit {
		ctx.stmtCounter = len(ctx.dialog.Statements)
		return
	}
	ctx.stmtCounter = ctx.dialog.LabelIndex[label]
}

func (ctx *context) Execute(code string) {
	_, _ = ctx.scripting.Eval(code)
}

func (ctx *context) Choice(index int, text, gotoLabel string) {
	if ctx.parrot {
		ctx.Say(ctx.currentActor, text)
	}
	ctx.Goto(gotoLabel)
}

func (ctx *context) ShutUp() {
	// TODO: implement
}

func (ctx *context) WaitFor(actor string) {
	// TODO: implement
}

func (ctx *context) WaitWhile(code string) {
	// TODO: implement
}

func (ctx *context) Override(label string) {
	// TODO: implement
}

func (ctx *context) IsActor(actor string) bool {
	return ctx.currentActor == actor
}

func (ctx *context) IsCodeTrue(code string) bool {
	result, err := ctx.scripting.Eval(code)
	if err != nil {
		return false
	}
	yes, ok := result.(bool)
	return ok && yes
}

func (ctx *context) IsOnce() bool {
	s, ok := ctx.currentStatement()
	return ok && !ctx.executed[s.Statement]
}

func (ctx *context) IsShowOnce() bool {
	s, ok := ctx.currentStatement()
	return ok && !ctx.shown[s.Statement]
}

func (ctx *context) IsOnceEver() bool {
	// TODO
	return ctx.IsOnce()
}

func (ctx *context) IsTempOnce() bool {
	// TODO
	return ctx.IsOnce()
}

func (ctx *context) run() *Choices {
	for ctx.stmtCounter < len(ctx.dialog.Statements) {
		condStmt, ok := ctx.currentStatement()
		if !ok {
			break
		}
		counterBeforeExec := ctx.stmtCounter
		if condStmt.Conditions.AreFulfilled(ctx) {
			if choice, ok := condStmt.Statement.(*stmt.Choice); ok {
				ctx.addChoice(choice)
			} else {
				ctx.execute(condStmt.Statement)
			}
		}
		if ctx.stmtCounter == counterBeforeExec {
			ctx.stmtCounter++
		}
		if !ctx.choicesReady() {
			continue
		}
		return &Choices{
			Actor:   ctx.currentActor,
			Options: ctx.choiceOptions(),
			ctx:     ctx,
		}
	}
	return &Choices{Actor: ctx.currentActor, ctx: ctx}
}

func (ctx *context) execute(statement Statement) {
	statement.Execute(ctx)
	ctx.executed[statement] = true
}

func (ctx *context) currentStatement() (ConditionalStatement, bool) {
	if ctx.stmtCounter >= len(ctx.dialog.Statements) {
		return ConditionalStatement{}, false
	}
	return ctx.dialog.Statements[ctx.stmtCounter], true
}

func (ctx *context) choicesReady() bool {
	currentStmt, ok := ctx.currentStatement()
	if !ok {
		return false
	}
	_, isCurrentStmtChoice := currentStmt.Statement.(*stmt.Choice)
	return len(ctx.choices) > 0 && !isCurrentStmtChoice
}

func (ctx *context) choiceOptions() []*ChoiceOption {
	options := make([]*ChoiceOption, len(ctx.choices))
	i := 0
	for _, choice := range ctx.choices {
		options[i] = &ChoiceOption{
			Text:   ctx.evalText(choice.Text),
			choice: choice,
		}
		i++
	}
	sort.Slice(options, func(i, j int) bool {
		return options[i].choice.Index < options[j].choice.Index
	})
	return options
}

func (ctx *context) choose(opt *ChoiceOption) *Choices {
	ctx.clearChoices()
	ctx.execute(opt.choice)
	return ctx.run()
}

func (ctx *context) clearChoices() {
	ctx.choices = make(map[int]*stmt.Choice)
}

func (ctx *context) addChoice(choice *stmt.Choice) {
	if _, exists := ctx.choices[choice.Index]; exists {
		return
	}
	ctx.choices[choice.Index] = choice
	ctx.shown[choice] = true
}

func (ctx *context) evalText(text string) string {
	if len(text) == 0 || text[0] != '$' {
		return text
	}
	code := text[1:]
	result, err := ctx.scripting.Eval(code)
	if err != nil {
		return "(script error)"
	}
	evaluatedText, ok := result.(string)
	if !ok {
		return "(script error: not a string)"
	}
	return evaluatedText
}
