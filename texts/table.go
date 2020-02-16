// Copyright 2020 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

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
				_, err = w.Write([]byte(t[textID]))
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

func (t Table) ResolveTextsString(s string) (string, error) {
	var sb strings.Builder
	err := t.ResolveTexts(&sb, strings.NewReader(s))
	if err != nil {
		return s, fmt.Errorf("could not resolve texts in string: %w", err)
	}
	return sb.String(), nil
}

func isNumeric(b byte) bool {
	return b >= '0' && b <= '9'
}
