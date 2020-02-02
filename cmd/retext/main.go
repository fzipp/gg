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
	"bufio"
	"bytes"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
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

	texts, err := loadTextsFromFile(*textsFilePath)
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
		err := replaceTexts(w, os.Stdin, texts)
		check(err)
	}
}

func processFile(w io.Writer, path string, texts textTable) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("could not open source file: %w", err)
	}
	defer file.Close()
	return replaceTexts(w, file, texts)
}

const idMarker = '@'

func replaceTexts(w io.Writer, r io.Reader, texts textTable) error {
	withinTextID := false
	var bufferedTextID []byte
	br := bufio.NewReader(r)
	for {
		b, err := br.ReadByte()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("could not read from input: %w", err)
		}
		if b == idMarker {
			withinTextID = true
			bufferedTextID = bufferedTextID[:0]
			continue
		}
		if withinTextID {
			if isNumeric(b) {
				bufferedTextID = append(bufferedTextID, b)
				continue
			}
			if len(bufferedTextID) > 0 {
				textID, err := strconv.Atoi(string(bufferedTextID))
				if err != nil {
					return fmt.Errorf("could not parse text ID: %w", err)
				}
				_, err = w.Write([]byte(texts[textID]))
				if err != nil {
					return fmt.Errorf("could not write replacement text to output: %w", err)
				}
			} else {
				_, err := w.Write([]byte{idMarker})
				if err != nil {
					return  fmt.Errorf("could not write text ID marker to output: %w", err)
				}
			}
			withinTextID = false
		}
		_, err = w.Write([]byte{b})
		if err != nil {
			return fmt.Errorf("could not write to output: %w", err)
		}
	}
	return nil
}

func isNumeric(b byte) bool {
	return b >= '0' && b <= '9'
}

type textTable map[int]string

func loadTextsFromFile(filepath string) (textTable, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("could not open text table file: %w", err)
	}
	defer file.Close()
	return loadTexts(file)
}

func loadTexts(r io.Reader) (textTable, error) {
	csvReader := csv.NewReader(r)
	csvReader.Comma = '\t'
	csvReader.LazyQuotes = true
	texts := make(textTable)
	lineNumber := 0
	for {
		lineNumber++
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("could not read TSV record: %w", err)
		}
		if lineNumber == 1 {
			// skip header
			continue
		}
		id, err := strconv.Atoi(record[0])
		if err != nil {
			return nil, fmt.Errorf("could not parse text ID: %w", err)
		}
		texts[id] = record[1]
	}
	return texts, nil
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
