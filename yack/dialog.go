// Copyright 2020 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package yack

import (
	"sort"
	"strings"

	"github.com/fzipp/gg/yack/cond"
	"github.com/fzipp/gg/yack/stmt"
)

// Dialog is a parsed yack dialog. It consists of a sequence of statements,
// each one guarded by zero or more conditions, and a label index. Each label
// points via index to a statement in the statements slice.
type Dialog struct {
	Statements []ConditionalStatement
	Labels     map[string]int
}

// String formats the dialog in yack syntax.
func (d *Dialog) String() string {
	var sb strings.Builder
	ll := d.buildLabelsLookup()
	for i, statement := range d.Statements {
		ll.writeLabels(&sb, i)
		sb.WriteString(statement.String() + "\n")
	}
	ll.writeLabels(&sb, len(d.Statements))
	return sb.String()
}

func (d *Dialog) buildLabelsLookup() labelsLookup {
	ll := make(labelsLookup)
	for label, i := range d.Labels {
		ll[i] = append(ll[i], label)
	}
	return ll
}

type labelsLookup map[int][]string

func (ll labelsLookup) writeLabels(sb *strings.Builder, index int) {
	if labels, ok := ll[index]; ok {
		sort.Strings(labels)
		for i, label := range labels {
			if i == 0 {
				sb.WriteRune('\n')
			}
			sb.WriteString(":" + label + "\n")
		}
	}
}

// A ConditionalStatement is a statement guarded by zero or more conditions.
type ConditionalStatement struct {
	Statement  stmt.Statement
	Conditions []cond.Condition
}

// String formats the conditional statement in yack syntax, e.g.
// "statement [condition1] [condition2] [condition3]"
func (c ConditionalStatement) String() string {
	s := c.Statement.String()
	cs := conditions(c.Conditions).String()
	if cs == "" {
		return s
	}
	return s + " " + cs
}

// conditions are zero or more conditions to guard a statement in a dialog script.
type conditions []cond.Condition

// AreFulfilled returns true if all conditions are fulfilled.
func (c conditions) AreFulfilled(ctx cond.Context) bool {
	for _, cond := range c {
		if !cond.IsFulfilled(ctx) {
			return false
		}
	}
	return true
}

// String formats a set of conditions in yack syntax, e.g.
// "[condition1] [condition2] [condition3]"
func (c conditions) String() string {
	cs := make([]string, len(c))
	for i, cond := range c {
		cs[i] = cond.String()
	}
	return strings.Join(cs, " ")
}
