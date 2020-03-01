// Copyright 2020 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package texts_test

import (
	"errors"
	"io"
	"reflect"
	"strings"
	"testing"

	"github.com/fzipp/gg/texts"
)

const textTableTSV = `text_id	en
10001	hello
10002	world
10003	This is a test
`

func TestResolveTextsString(t *testing.T) {
	textTable, err := texts.From(strings.NewReader(textTableTSV))
	if err != nil {
		t.Fatalf("could not read text table: %s", err)
	}
	tests := []struct {
		text string
		want string
	}{
		{"@10001", "hello"},
		{"@10001, @10002", "hello, world"},
		{"@10001, @10002 @", "hello, world @"},
		{"@10003: @10001, @10002", "This is a test: hello, world"},
		{"@abc @def@", "@abc @def@"},
		{"@90001", "@90001"},
	}
	for _, tt := range tests {
		resolved, err := textTable.ResolveTextsString(tt.text)
		if err != nil {
			t.Errorf("could not writeResolved text string %q, got error: %s", tt.text, err)
		}
		if resolved != tt.want {
			t.Errorf("resolved text for %q was %q, want %q", tt.text, resolved, tt.want)
		}
	}
}

func TestFromFile(t *testing.T) {
	tests := []struct {
		path string
		want texts.Table
	}{
		{"testdata/TestTexts.tsv", texts.Table{
			12345: "hello, world",
			20001: "This is a test.",
			20003: "More text...",
			20002: `A text with \"quotes\"`,
		}},
	}
	for _, tt := range tests {
		textTable, err := texts.FromFile(tt.path)
		if err != nil {
			t.Errorf("could not load text table from file %q: %s", tt.path, err)
			continue
		}
		if !reflect.DeepEqual(textTable, tt.want) {
			t.Errorf("text table loaded from file %q was: %#v, want: %#v", tt.path, textTable, tt.want)
		}
	}
}

func TestFromFileErrors(t *testing.T) {
	tests := []struct {
		path      string
		wantError string
	}{
		{"testdata/DoesNotExist.tsv", "could not open text table file: open testdata/DoesNotExist.tsv: no such file or directory"},
	}
	for _, tt := range tests {
		_, err := texts.FromFile(tt.path)
		if err == nil {
			t.Errorf("expected error when reading text table file %q, but no error returned", tt.path)
			continue
		}
		if err.Error() != tt.wantError {
			t.Errorf("error message when reading text table file %q was: %q, want: %q", tt.path, err.Error(), tt.wantError)
		}
	}
}

func TestFromErrors(t *testing.T) {
	tests := []struct {
		input     string
		wantError string
	}{
		{"text_id\ten\nabc def", "could not read TSV record: record on line 2: wrong number of fields"},
		{"text_id\ten\nabc\tdef", `could not parse text ID: strconv.Atoi: parsing "abc": invalid syntax`},
	}
	for _, tt := range tests {
		_, err := texts.From(strings.NewReader(tt.input))
		if err == nil {
			t.Errorf("expected error for parsing text table %q, but no error returned", tt.input)
			continue
		}
		if err.Error() != tt.wantError {
			t.Errorf("error message for parsing of text table %q was: %q, want: %q", tt.input, err.Error(), tt.wantError)
		}
	}
}

func TestResolveTextsErrors(t *testing.T) {
	tests := []struct {
		caseName  string
		w         io.Writer
		r         io.Reader
		wantError string
	}{
		{"read input error", &strings.Builder{}, errorReader{}, "could not read from input: test read error"},
		{"write output error", errorWriter{}, strings.NewReader("text_id\ten\n"), "test write error"},
	}
	textTable, err := texts.From(strings.NewReader("text_id\ten\n"))
	if err != nil {
		t.Fatalf("parsing of text table for tests returned error: %s", err)
	}
	for _, tt := range tests {
		err := textTable.ResolveTexts(tt.w, tt.r)
		if err == nil {
			t.Errorf("expected error for resolving texts, but no error returned")
			continue
		}
		if err.Error() != tt.wantError {
			t.Errorf("error message for test case %q was: %q, want: %q", tt.caseName, err.Error(), tt.wantError)
		}
	}
}

type errorReader struct{}

func (r errorReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("test read error")
}

type errorWriter struct{}

func (w errorWriter) Write(p []byte) (n int, err error) {
	return 0, errors.New("test write error")
}
