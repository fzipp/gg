// Copyright 2020 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package condition

type Actor struct {
	Actor string
}

func (c *Actor) IsFulfilled(ctx Context) bool {
	return ctx.IsActor(c.Actor)
}

func (c *Actor) String() string {
	return "[" + c.Actor + "]"
}

type Code struct {
	Code string
}

func (c *Code) IsFulfilled(ctx Context) bool {
	return ctx.IsCodeTrue(c.Code)
}

func (c *Code) String() string {
	return "[" + c.Code + "]"
}

type Once struct{}

func (c *Once) IsFulfilled(ctx Context) bool {
	return ctx.IsOnce()
}

func (c *Once) String() string {
	return "[once]"
}

type ShowOnce struct{}

func (c *ShowOnce) IsFulfilled(ctx Context) bool {
	return ctx.IsShowOnce()
}

func (c *ShowOnce) String() string {
	return "[showonce]"
}

type OnceEver struct{}

func (c *OnceEver) IsFulfilled(ctx Context) bool {
	return ctx.IsOnceEver()
}

func (c *OnceEver) String() string {
	return "[onceever]"
}

type TempOnce struct{}

func (c *TempOnce) IsFulfilled(ctx Context) bool {
	return ctx.IsTempOnce()
}

func (c *TempOnce) String() string {
	return "[temponce]"
}
