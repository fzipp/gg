package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
)

type TextTable map[int]string

func LoadTextsFromFile(path string) (TextTable, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("could not open text table file: %w", err)
	}
	defer file.Close()
	return LoadTexts(file)
}

func LoadTexts(r io.Reader) (TextTable, error) {
	csvReader := csv.NewReader(r)
	csvReader.Comma = '\t'
	csvReader.LazyQuotes = true
	texts := make(TextTable)
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

func (tt TextTable) InsertTexts(w io.Writer, r io.Reader) error {
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
				_, err = w.Write([]byte(tt[textID]))
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
