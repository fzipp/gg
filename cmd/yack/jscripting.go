// Copyright 2020 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/robertkrimen/otto"
)

type jscripting struct {
	vm      *otto.Otto
	verbose bool
}

func newJScripting() *jscripting {
	vm := otto.New()
	s := &jscripting{vm: vm}
	_, err := vm.Run(initCode)
	if err != nil {
		s.log(err)
	}
	return s
}

func (s *jscripting) Eval(code string) (result interface{}, err error) {
	code = squirrelToJS(code)
	val, err := s.vm.Run(code)
	if err != nil {
		s.log(err)
		return true, nil
	}
	return val.Export()
}

func (s *jscripting) log(err error) {
	if s.verbose {
		_, _ = fmt.Fprintf(os.Stderr, "Script error: %s\n", err)
	}
}

func squirrelToJS(code string) string {
	return strings.ReplaceAll(code, "<-", "=")
}

const initCode = `
var YES = true;
var NO = false;

var g = {};

var currentActor = 0;
var currentRoom = 0;

function random(min, max) {
	min = Math.ceil(min);
	max = Math.floor(max);
	return Math.floor(Math.random() * (max - min)) + min;
}

function randomOdds(p) {
	return Math.random() <= p;
}
`
