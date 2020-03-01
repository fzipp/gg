// Copyright 2020 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package texts replaces ID placeholders like @12345 with actual texts
// from a text table.
package texts

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type Table map[int]string

func FromFile(path string) (Table, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("could not open text table file: %w", err)
	}
	defer file.Close()
	return From(file)
}

func From(r io.Reader) (Table, error) {
	csvReader := csv.NewReader(r)
	csvReader.Comma = '\t'
	csvReader.LazyQuotes = true
	texts := make(Table)
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

const idMarker = '@'

func (t Table) ResolveTexts(w io.Writer, r io.Reader) error {
	eof := false
	withinTextID := false
	var bufferedTextID strings.Builder
	br := bufio.NewReader(r)
	bw := bufio.NewWriter(w)
	for {
		ch, _, err := br.ReadRune()
		if err == io.EOF {
			eof = true
		}
		if err != nil && !eof {
			return fmt.Errorf("could not read from input: %w", err)
		}
		if ch == idMarker {
			withinTextID = true
			bufferedTextID.Reset()
			continue
		}
		if withinTextID {
			if isNumeric(ch) {
				bufferedTextID.WriteRune(ch)
				continue
			}
			if bufferedTextID.Len() > 0 {
				textID, err := strconv.Atoi(bufferedTextID.String())
				if err != nil {
					return fmt.Errorf("could not parse text ID: %w", err)
				}
				text, ok := t[textID]
				if !ok {
					text = string(idMarker) + bufferedTextID.String()
				}
				_, err = bw.WriteString(text)
				if err != nil {
					return fmt.Errorf("could not write replacement text to output: %w", err)
				}
			} else {
				_, err := bw.WriteRune(idMarker)
				if err != nil {
					return fmt.Errorf("could not write text ID marker to output: %w", err)
				}
			}
			withinTextID = false
		}
		if eof {
			break
		}
		_, err = bw.WriteRune(ch)
		if err != nil {
			return fmt.Errorf("could not write to output: %w", err)
		}
	}
	return bw.Flush()
}

func (t Table) ResolveTextsString(s string) (string, error) {
	var sb strings.Builder
	err := t.ResolveTexts(&sb, strings.NewReader(s))
	if err != nil {
		return s, fmt.Errorf("could not resolve texts in string: %w", err)
	}
	return sb.String(), nil
}

func isNumeric(r rune) bool {
	return r >= '0' && r <= '9'
}
