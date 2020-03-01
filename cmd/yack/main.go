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
	"bytes"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
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

	path := flag.Arg(0)
	file, err := os.Open(path)
	check(err)

	var buf bytes.Buffer
	err = textTable.ResolveTexts(&buf, file)
	check(err)
	dialog, err := yack.Parse(filepath.Base(path), &buf)
	check(err)
	err = file.Close()
	check(err)

	run(dialog, *startActorFlag, *startLabelFlag)
}

func run(dialog *yack.Dialog, startActor, startLabel string) {
	checkStartLabelExists(dialog, startLabel)
	talk := &consoleTalk{}
	runner := yack.NewRunner(dialog, noScripting{}, talk, startActor)
	runner.Init()
	choices := runner.StartAt(startLabel)
	for len(choices.Options) > 0 {
		printChoices(choices)
		prompt := choices.Actor + "> "
		input := userInput(prompt, 1, len(choices.Options))
		opt := choices.Options[input-1]
		fmt.Println()
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

func printChoices(choices *yack.Choices) {
	fmt.Println()
	for i, opt := range choices.Options {
		fmt.Printf("%d) %s\n", i+1, opt.Text)
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

var animationTagsRegexp = regexp.MustCompile(`\^?{.*}`)

func stripAnimationTags(text string) string {
	return animationTagsRegexp.ReplaceAllString(text, "")
}

func check(err error) {
	if err != nil {
		fail(err)
	}
}

func fail(message interface{}) {
	_, _ = fmt.Fprintln(os.Stderr, message)
	os.Exit(1)
}

type consoleTalk struct{}

func (t *consoleTalk) Say(actor, text string) {
	text = stripAnimationTags(text)
	if text == "" {
		return
	}
	fmt.Printf("%s: %s\n", actor, text)
	time.Sleep(time.Duration(len(text)) * 70 * time.Millisecond)
}

type noScripting struct{}

func (s noScripting) Eval(code string) (result interface{}, err error) {
	// do nothing, always return true
	return true, nil
}
