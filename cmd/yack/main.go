// Copyright 2020 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Runs yack dialogs on the console. Scripting language statements are ignored.
// Scripting language conditions always evaluate to true.
//
// Usage:
//     yack [-l start_label] [-a start_actor] [-t texts_file] yack_file
//
// Flags:
//     -l    The start label. Default: "start"
//     -a    The start actor. Default: "you"
//     -t    A text table file in TSV (tab-separated values) format to look up
//           text IDs (i.e. "@12345") and replace them with actual texts.
//
// Examples:
//     yack ExampleDialog.yack
//     yack -t ExampleTexts.tsv ExampleDialog.yack
//     yack -l introduction -t ExampleTexts.tsv ExampleDialog.yack
package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/fzipp/gg/texts"
	"github.com/fzipp/gg/yack"
)

func usage() {
	fail(`Runs yack dialogs on the console. Scripting language statements are ignored.
Scripting language conditions always evaluate to true.

Flags:
    -l  The start label. Default: "start"
    -a  The start actor. Default: "you"
    -t  A text table file in TSV (tab-separated values) format to look up
        text IDs (i.e. "@12345") and replace them with actual texts.

Usage:
    yack [-l start_label] [-a start_actor] [-t texts_file] yack_file

Examples:
    yack ExampleDialog.yack
	yack -t ExampleTexts.tsv ExampleDialog.yack
	yack -l introduction -t ExampleTexts.tsv ExampleDialog.yack`)
}

func main() {
	textsFileFlag := flag.String("t", "", "path to a text database file (TSV format)")
	startLabelFlag := flag.String("l", "start", "start label")
	startActorFlag := flag.String("a", "you", "start actor name")

	flag.Usage = usage
	flag.Parse()

	if flag.NArg() < 1 {
		usage()
	}

	var err error
	var textTable texts.Table

	if *textsFileFlag != "" {
		textTable, err = texts.FromFile(*textsFileFlag)
		if err != nil {
			fail(err)
		}
	}

	dialog, err := yack.Load(flag.Arg(0))
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
	}
	run(dialog, *startActorFlag, *startLabelFlag, textTable)
}

func run(dialog *yack.Dialog, startActor, startLabel string, textTable texts.Table) {
	checkStartLabelExists(dialog, startLabel)
	talk := &consoleTalk{textTable}
	runner := yack.NewRunner(dialog, noScripting{}, talk, startActor)
	runner.Init()
	choices := runner.StartAt(startLabel)
	for len(choices.Options) > 0 {
		printChoices(choices, textTable)
		input := userInput(1, len(choices.Options))
		opt := choices.Options[input-1]
		choices = choices.Choose(opt)
	}
}

func checkStartLabelExists(dialog *yack.Dialog, startLabel string) {
	if _, ok := dialog.LabelIndex[startLabel]; !ok {
		if startLabel == "start" {
			fail(fmt.Sprintf("Label %q not found. Try passing a different start label with the -l flag.", startLabel))
		}
		fail(fmt.Sprintf("Label %q not found.", startLabel))
	}
}

func printChoices(choices *yack.Choices, textTable texts.Table) {
	fmt.Println()
	for i, opt := range choices.Options {
		fmt.Printf("%d) %s\n", i+1, maybeResolve(opt.Text, textTable))
	}
}

func userInput(min, max int) int {
	for {
		fmt.Print("> ")
		var no int
		_, err := fmt.Scanf("%d", &no)
		if err != nil || no < min || no > max {
			continue
		}
		return no
	}
}

func fail(message interface{}) {
	_, _ = fmt.Fprintln(os.Stderr, message)
	os.Exit(1)
}

type consoleTalk struct {
	textTable texts.Table
}

func (t *consoleTalk) Say(actor, text string) {
	fmt.Printf("%s: %s\n", actor, maybeResolve(text, t.textTable))
	time.Sleep(time.Duration(len(text)) * 500 * time.Millisecond)
}

type noScripting struct{}

func (s noScripting) Eval(code string) (result interface{}, err error) {
	// do nothing, always return true
	return true, nil
}

func maybeResolve(text string, textTable texts.Table) string {
	if textTable == nil {
		return text
	}
	text, err := textTable.ResolveTextsString(text)
	if err != nil {
		fail(err)
	}
	return text
}
