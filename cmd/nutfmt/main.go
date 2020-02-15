// Copyright 2020 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Re-indents Squirrel scripting language source code. This is not a full-blown
// formatter like gofmt and does not aspire to be one. It was created for the
// limited purpose to be used in the context of the other tools in this module.
//
// Usage:
//     nutfmt [-w] [path ...]
//
// Flags:
//     -w  Write result to the source file instead of standard output.
//
// Examples:
//     nutfmt Example.bnut | less
//     nutfmt -w Example.bnut
package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"
)

func usage() {
	fail(`Re-indents Squirrel scripting language source code. This is not a full-blown
formatter like gofmt and does not aspire to be one.

Usage:
    nutfmt [-w] [path ...]

Flags:
    -w  Write result to the source file instead of standard output.

Examples:
    nutfmt Example.bnut | less
    nutfmt -w Example.bnut`)
}

func main() {
	if len(os.Args) < 2 {
		usage()
	}

	file, err := os.Open(os.Args[1])
	check(err)
	defer file.Close()

	err = prettify(os.Stdout, file)
	check(err)
}

func prettify(w io.Writer, r io.Reader) error {
	indent := "    "
	indentLevel := 0
	onNewLine := true
	var prev rune

	bw := bufio.NewWriter(w)
	br := bufio.NewReader(r)

	for {
		curr, _, err := br.ReadRune()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if isIndentDecreaser(curr) {
			indentLevel--
		}
		if onNewLine {
			if !unicode.IsSpace(curr) || curr == '\n' {
				_, err := bw.WriteString(strings.Repeat(indent, indentLevel))
				if err != nil {
					return fmt.Errorf("could not write indentation: %w", err)
				}
				_, err = bw.WriteRune(curr)
				if err != nil {
					return fmt.Errorf("could not write rune: %w", err)
				}
				onNewLine = false
			}
		} else {
			_, err := bw.WriteRune(curr)
			if err != nil {
				return fmt.Errorf("could not write rune: %w", err)
			}
		}
		if isIndentIncreaser(curr) {
			indentLevel++
		}

		prev = curr
		if prev == '\n' {
			onNewLine = true
		}
	}

	return bw.Flush()
}

func isIndentIncreaser(r rune) bool {
	return r == '{' || r == '(' || r == '['
}

func isIndentDecreaser(r rune) bool {
	return r == '}' || r == ')' || r == ']'
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
