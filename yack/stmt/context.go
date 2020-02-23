// Copyright 2020 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stmt

type Context interface {
	ShutUp()
	Say(actor, text string)
	Pause(seconds float64)
	Parrot(enabled bool)
	WaitFor(actor string)
	WaitWhile(code string)
	Dialog(actor string)
	Override(label string)
	AllowObjects(allow bool)
	Limit(n int)
	Goto(label string)
	Execute(code string)
	Choice(index int, text, gotoLabel string)
}
