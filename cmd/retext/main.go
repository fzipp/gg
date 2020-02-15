// Copyright 2020 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Replaces ID placeholders like @12345 in files with texts from a text
// table file in TSV (tab-separated values) format referenced via these IDs.
//
// Usage:
//     retext -t text_table_file [-w] [path ...]
//
// Flags:
//     -t  The text table file. This flag is mandatory. The expected format
//         of the file is TSV (tab-separated values). The first column contains
//         the numeric text IDs (without the @ marker), the second column the
//         texts. The first row is ignored and can contain column headers.
//     -w  Write result to the source file instead of standard output.
//
// Example text table file "texts.tsv":
//
//    text_id en
//    20001	Hi, how are you?
//    20002	Thanks, I'm fine.
//
// Example input file "story.txt":
//
//    She asked "@20001" and he answered "@20002"
//
// Example usage:
//
//    retext -t texts.tsv story.txt > story_complete.txt
//
// Example output file "story_complete.txt":
//
//    She asked "Hi, how are you?" and he answered "Thanks, I'm fine."
//
// Example bulk processing:
//
//    retext -t texts.tsv -w *.txt
package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

func usage() {
	fail(`Replaces ID placeholders like @12345 in files with texts from a text
table file in TSV (tab-separated values) format referenced via these IDs.

Usage:
    retext -t text_table_file [-o output_file | -w] [path ...]

Flags:
    -t  The text table file. This flag is mandatory. The expected format
        of the file is TSV (tab-separated values). The first column contains
        the numeric text IDs (without the @ marker), the second column the
        texts. The first row is ignored and can contain column headers.
    -w  Write result to the source file instead of standard output.

Example text table file:

text_id en
20001	Where there's a will there's a way.
20002	There's no smoke without fire.`)
}

func main() {
	textsFilePath := flag.String("t", "", "texts file in TSV format")
	replaceSource := flag.Bool("w", false, "write result to (source) file instead of stdout")

	flag.Usage = usage
	flag.Parse()

	inputFiles := flag.Args()

	if *textsFilePath == "" {
		fail("Please specify a text table file via -t. See -help for more information.")
		return
	}
	if *replaceSource && (len(inputFiles) == 0) {
		fail("Cannot use -w with standard input. See -help for more information.")
		return
	}

	texts, err := LoadTextsFromFile(*textsFilePath)
	check(err)

	var w io.Writer = os.Stdout

	for _, inputFile := range inputFiles {
		if *replaceSource {
			buf := &bytes.Buffer{}
			err = processFile(buf, inputFile, texts)
			check(err)
			err = ioutil.WriteFile(inputFile, buf.Bytes(), 0644)
			check(err)
		} else {
			err := processFile(w, inputFile, texts)
			check(err)
		}
	}
	if len(inputFiles) == 0 {
		err := texts.InsertTexts(w, os.Stdin)
		check(err)
	}
}

func processFile(w io.Writer, path string, texts TextTable) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("could not open source file: %w", err)
	}
	defer file.Close()
	return texts.InsertTexts(w, file)
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
