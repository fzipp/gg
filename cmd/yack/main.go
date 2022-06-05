// Copyright 2020 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// A tool to run yack dialogs on the console.
//
// The scripting language used for code statements and code conditions is not
// Squirrel, but JavaScript. However, Squirrel's "<-" assignment operator can
// be used and simply gets replaced with "=". YES and NO are also pre-defined
// as true and false.
//
// Usage:
//
//	yack [-t texts_file] [-l start_label] [-a start_actor] yack_file
//
// Flags:
//
//	-t  A text table file in TSV (tab-separated values) format to look up
//	    text IDs (i.e. "@12345") and replace them with actual texts.
//	-l  The start label. Default: "start"
//	-a  The start actor. Default: "you"
//	-d  Show debug information, such as animation tags and script errors.
//
// Examples:
//
//	yack ExampleDialog.yack
//	yack -t ExampleTexts.tsv ExampleDialog.yack
//	yack -t ExampleTexts.tsv -l main ExampleDialog.yack
package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"regexp"
	"time"

	"github.com/fzipp/gg/texts"
	"github.com/fzipp/gg/yack"
)

func usage() {
	fail(`Runs yack dialogs on the console.

The scripting language used for code statements and code conditions is not
Squirrel, but JavaScript. However, Squirrel's "<-" assignment operator can
be used and simply gets replaced with "=". YES and NO are also pre-defined
as true and false.

Flags:
    -t  A text table file in TSV (tab-separated values) format to look up
        text IDs (i.e. "@12345") and replace them with actual texts.
    -l  The start label. Default: "start"
    -a  The start actor. Default: "you"
    -d  Show debug information, such as animation tags and script errors.

Usage:
    yack [-t texts_file] [-l start_label] [-a start_actor] [-d] yack_file

Examples:
    yack ExampleDialog.yack
    yack -t ExampleTexts.tsv ExampleDialog.yack
    yack -t ExampleTexts.tsv -l main ExampleDialog.yack`)
}

func main() {
	textsFileFlag := flag.String("t", "", "path to a text database file (TSV format)")
	startLabelFlag := flag.String("l", "start", "start label")
	startActorFlag := flag.String("a", "you", "start actor name")
	debugFlag := flag.Bool("d", false, "show debugMode information")

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

	path := flag.Arg(0)
	file, err := os.Open(path)
	check(err)

	var buf bytes.Buffer
	err = textTable.ResolveTexts(&buf, file)
	check(err)
	dialog, err := yack.Read(&buf)
	check(err)
	err = file.Close()
	check(err)

	run(dialog, *startActorFlag, *startLabelFlag, *debugFlag)
}

func run(dialog *yack.Dialog, startActor, startLabel string, debugMode bool) {
	checkStartLabelExists(dialog, startLabel)
	talk := &consoleTalk{debugMode: debugMode}
	scripting := newJScripting()
	scripting.verbose = debugMode
	runner := yack.NewRunner(scripting, talk, startActor)
	choices := runner.StartAt(dialog, startLabel)
	for len(choices.Options) > 0 {
		printChoices(choices, debugMode)
		prompt := choices.Actor + "> "
		input := userInput(prompt, 1, len(choices.Options))
		fmt.Println()
		choices = choices.Choose(input - 1)
	}
}

func checkStartLabelExists(dialog *yack.Dialog, startLabel string) {
	if _, ok := dialog.Labels[startLabel]; !ok {
		if startLabel == "start" {
			fail(fmt.Sprintf("Label %q not found. Try passing a different start label with the -l flag.", startLabel))
		}
		fail(fmt.Sprintf("Label %q not found.", startLabel))
	}
}

func printChoices(choices *yack.Choices, debugMode bool) {
	fmt.Println()
	for i, opt := range choices.Options {
		text := opt.Text
		if !debugMode {
			text = stripAnimationTags(text)
		}
		fmt.Printf("%d) %s\n", i+1, text)
	}
}

func userInput(prompt string, min, max int) int {
	for {
		fmt.Print(prompt)
		var no int
		_, err := fmt.Scanf("%d", &no)
		if err != nil || no < min || no > max {
			if min == max {
				fmt.Printf("Please choose option %d. You don't have a choice anyway.\n", min)
			} else {
				fmt.Printf("Please choose an option between %d and %d.\n", min, max)
			}
			continue
		}
		return no
	}
}

var (
	animationTagsRegexp   = regexp.MustCompile(`\^?{.*}`)
	parenthesesTagsRegexp = regexp.MustCompile(`^\(.*\)`)
)

func stripAnimationTags(text string) string {
	text = animationTagsRegexp.ReplaceAllString(text, "")
	return parenthesesTagsRegexp.ReplaceAllString(text, "")
}

func check(err error) {
	if err != nil {
		fail(err)
	}
}

func fail(message any) {
	_, _ = fmt.Fprintln(os.Stderr, message)
	os.Exit(1)
}

type consoleTalk struct {
	debugMode bool
}

func (t *consoleTalk) Say(actor, text string) {
	if !t.debugMode {
		text = stripAnimationTags(text)
	}
	if text == "" {
		return
	}
	fmt.Printf("%s: %s\n", actor, text)
	time.Sleep(time.Duration(len(text)) * 70 * time.Millisecond)
}

func (t *consoleTalk) ShutUp() {
	// do nothing
}
