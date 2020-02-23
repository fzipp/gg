// Copyright 2020 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package yack

import (
	"strings"

	"github.com/fzipp/gg/yack/condition"
	"github.com/fzipp/gg/yack/stmt"
)

// Dialog is a parsed yack dialog. It consists of a sequence of statements,
// each one guarded by zero or more conditions, and a label index. Each label
// points via index to a statement in the statements slice.
type Dialog struct {
	Statements []ConditionalStatement
	LabelIndex map[string]int
}

// String formats the dialog in yack syntax.
func (d *Dialog) String() string {
	labels := make(map[int]string)
	for label, i := range d.LabelIndex {
		labels[i] = label
	}
	var sb strings.Builder
	for i, statement := range d.Statements {
		if label, ok := labels[i]; ok {
			sb.WriteString("\n:" + label + "\n")
		}
		sb.WriteString(statement.String() + "\n")
	}
	return sb.String()
}

// A ConditionalStatement is a statement guarded by zero or more conditions.
type ConditionalStatement struct {
	Statement  Statement
	Conditions Conditions
}

// String formats the conditional statement in yack syntax, e.g.
// "statement [condition1] [condition2] [condition3]"
func (c ConditionalStatement) String() string {
	s := c.Statement.String()
	cs := c.Conditions.String()
	if cs == "" {
		return s
	}
	return  s + " " + cs
}

// Statement is an executable statement in a dialog script.
type Statement interface {
	Execute(ctx stmt.Context)
	String() string
}

// Condition is a condition to guard statement in a dialog script.
type Condition interface {
	IsFulfilled(ctx condition.Context) bool
	String() string
}

// Conditions are zero or more conditions to guard a statement in a dialog script.
type Conditions []Condition

func (c Conditions) AreFulfilled(ctx condition.Context) bool {
	for _, cond := range c {
		if !cond.IsFulfilled(ctx) {
			return false
		}
	}
	return true
}

// String formats a set of conditions in yack syntax, e.g.
// "[condition1] [condition2] [condition3]"
func (c Conditions) String() string {
	cs := make([]string, len(c))
	for i, cond := range c {
		cs[i] = cond.String()
	}
	return strings.Join(cs, " ")
}
