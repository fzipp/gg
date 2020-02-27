// Copyright 2020 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stmt

import "strconv"

type ShutUp struct{}

func (s ShutUp) Execute(ctx Context) {
	ctx.ShutUp()
}

func (s ShutUp) String() string {
	return "shutup"
}

type Say struct {
	Actor             string
	Text              string
	OptionalGotoLabel string
}

func (s Say) Execute(ctx Context) {
	ctx.Say(s.Actor, s.Text)
	if s.OptionalGotoLabel != "" {
		ctx.Goto(s.OptionalGotoLabel)
	}
}

func (s Say) String() string {
	gotoLabel := ""
	if s.OptionalGotoLabel != "" {
		gotoLabel = " -> " + s.OptionalGotoLabel
	}
	return s.Actor + ": " + strconv.Quote(s.Text) + gotoLabel
}

type Pause struct {
	Seconds float64
}

func (s Pause) Execute(ctx Context) {
	ctx.Pause(s.Seconds)
}

func (s Pause) String() string {
	return "pause " + strconv.FormatFloat(s.Seconds, 'g', -1, 64)
}

type Execute struct {
	Code string
}

func (s Execute) Execute(ctx Context) {
	ctx.Execute(s.Code)
}

func (s Execute) String() string {
	return "!" + s.Code
}

type Goto struct {
	Label string
}

func (s Goto) Execute(ctx Context) {
	ctx.Goto(s.Label)
}

func (s Goto) String() string {
	return "-> " + s.Label
}

type Choice struct {
	Index     int
	Text      string
	GotoLabel string
}

func (s Choice) Execute(ctx Context) {
	ctx.Choice(s.Index, s.Text, s.GotoLabel)
}

func (s Choice) String() string {
	return strconv.Itoa(s.Index) + " " + strconv.Quote(s.Text) + " -> " + s.GotoLabel
}

type WaitFor struct {
	Actor string
}

func (s WaitFor) Execute(ctx Context) {
	ctx.WaitFor(s.Actor)
}

func (s WaitFor) String() string {
	waitfor := "waitfor"
	if s.Actor == "" {
		return waitfor
	}
	return waitfor + " " + s.Actor
}

type WaitWhile struct {
	CodeCondition string
}

func (s WaitWhile) Execute(ctx Context) {
	ctx.WaitWhile(s.CodeCondition)
}

func (s WaitWhile) String() string {
	return "waitwhile " + s.CodeCondition
}

type Parrot struct {
	Enabled bool
}

func (s Parrot) Execute(ctx Context) {
	ctx.Parrot(s.Enabled)
}

func (s Parrot) String() string {
	return "parrot " + boolToString(s.Enabled)
}

type Dialog struct {
	Actor string
}

func (s Dialog) Execute(ctx Context) {
	ctx.Dialog(s.Actor)
}

func (s Dialog) String() string {
	return "dialog " + s.Actor
}

type Override struct {
	Label string
}

func (s Override) Execute(ctx Context) {
	ctx.Override(s.Label)
}

func (s Override) String() string {
	return "override " + s.Label
}

type AllowObjects struct {
	Allow bool
}

func (s AllowObjects) Execute(ctx Context) {
	ctx.AllowObjects(s.Allow)
}

func (s AllowObjects) String() string {
	return "allowobjects " + boolToString(s.Allow)
}

type Limit struct {
	N int
}

func (s Limit) Execute(ctx Context) {
	ctx.Limit(s.N)
}

func (s Limit) String() string {
	return "limit " + strconv.Itoa(s.N)
}

func boolToString(b bool) string {
	if b {
		return "YES"
	}
	return "NO"
}
