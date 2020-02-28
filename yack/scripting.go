// Copyright 2020 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package yack

// Scripting provides an interface for the evaluation of expressions in a
// scripting language. It can be used to plug in an interpreter or VM for
// any scripting language, e.g. Squirrel, Lua or JavaScript.
type Scripting interface {
	Eval(code string) (result interface{}, err error)
}
