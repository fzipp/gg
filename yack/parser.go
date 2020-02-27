// Copyright 2020 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package yack

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/scanner"
	"unicode"

	"github.com/fzipp/gg/yack/condition"
	"github.com/fzipp/gg/yack/stmt"
)

// Load reads and parses a dialog from a yack file. File IO or syntax errors
// are returned as error.
func Load(path string) (*Dialog, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("could not open yack file '%s': %w", path, err)
	}
	defer file.Close()
	return Parse(filepath.Base(path), file)
}

// Parse parses a dialog from a source in yack format. Syntax errors are
// returned as error. The filename parameter is only used as prefix for the
// error messages.
func Parse(filename string, src io.Reader) (*Dialog, error) {
	var p parser
	return p.parse(filename, src)
}

type parser struct {
	errors  errorList
	scanner scanner.Scanner
	pos     scanner.Position
	tok     rune
	lit     string
}

func (p *parser) next() {
	p.tok = p.scanner.Scan()
	p.pos = p.scanner.Position
	p.lit = p.scanner.TokenText()
}

func (p *parser) parse(filename string, src io.Reader) (*Dialog, error) {
	p.scanner.Init(src)
	p.scanner.Filename = filename
	p.scanner.Whitespace ^= 1 << '\n'
	p.scanner.Error = func(s *scanner.Scanner, msg string) {
		p.error(s.Position, msg)
	}
	p.next()
	d := &Dialog{
		Statements: make([]ConditionalStatement, 0),
		LabelIndex: make(map[string]int),
	}
	for p.tok != scanner.EOF {
		if p.tok == ':' {
			p.next()
			label := p.parseIdentifier()
			d.LabelIndex[label] = len(d.Statements)
			p.expectCommentOrNewLine()
			continue
		}
		condStmt := p.parseConditionalStatement()
		if condStmt.Statement != nil {
			d.Statements = append(d.Statements, condStmt)
		}
	}
	return d, p.errors.Err()
}

func (p *parser) parseConditionalStatement() ConditionalStatement {
	statement := p.parseStatement()
	conditions := p.parseConditions()
	p.expectCommentOrNewLine()
	return ConditionalStatement{
		Statement:  statement,
		Conditions: conditions,
	}
}

func (p *parser) parseStatement() Statement {
	switch p.tok {
	case '!':
		p.tok = p.scanner.Next()
		code := p.parseCode()
		return stmt.Execute{Code: code}
	case '-':
		label := p.parseGoto()
		return stmt.Goto{Label: label}
	case scanner.Int:
		index := p.parseInt()
		var text string
		if p.tok == '$' {
			text = p.parseCode()
		} else {
			text = p.parseString()
		}
		gotoLabel := p.parseGoto()
		return stmt.Choice{Index: index, Text: text, GotoLabel: gotoLabel}
	case scanner.Ident:
		ident := p.lit
		if p.scanner.Peek() == ':' {
			p.next()
			p.next()
			text := p.parseString()
			gotoLabel := ""
			if p.tok == '-' {
				gotoLabel = p.parseGoto()
			}
			return stmt.Say{
				Actor: ident, Text: text,
				OptionalGotoLabel: gotoLabel,
			}
		}
		switch ident {
		case "shutup":
			p.next()
			return stmt.ShutUp{}
		case "pause":
			p.next()
			seconds := p.parseFloat()
			return stmt.Pause{Seconds: seconds}
		case "waitfor":
			p.next()
			actor := ""
			if p.tok == scanner.Ident {
				actor = p.parseIdentifier()
			}
			return stmt.WaitFor{Actor: actor}
		case "waitwhile":
			p.tok = p.scanner.Next()
			code := p.parseCode()
			return stmt.WaitWhile{CodeCondition: code}
		case "parrot":
			p.next()
			enabled := p.parseBool()
			return stmt.Parrot{Enabled: enabled}
		case "dialog":
			p.next()
			actor := p.parseIdentifier()
			return stmt.Dialog{Actor: actor}
		case "override":
			p.next()
			label := p.parseIdentifier()
			return stmt.Override{Label: label}
		case "allowobjects":
			p.next()
			allow := p.parseBool()
			return stmt.AllowObjects{Allow: allow}
		case "limit":
			p.next()
			n := p.parseInt()
			return stmt.Limit{N: n}
		default:
			p.next()
			p.error(p.pos, fmt.Sprintf("invalid command: %s", ident))
		}
	}
	return nil
}

func (p *parser) parseConditions() Conditions {
	var conditions Conditions
	for p.tok == '[' {
		p.tok = p.scanner.Next()
		cond := p.parseCondition()
		p.expect(']')
		conditions = append(conditions, cond)
	}
	return conditions
}

func (p *parser) parseGoto() (label string) {
	p.expect('-')
	p.expect('>')
	return p.parseIdentifier()
}

func (p *parser) parseIdentifier() string {
	name := p.lit
	p.expect(scanner.Ident)
	return name
}

func (p *parser) parseInt() int {
	intLit := p.lit
	p.expect(scanner.Int)
	i, err := strconv.Atoi(intLit)
	if err != nil {
		p.error(p.pos, fmt.Sprintf("invalid integer literal: %s", intLit))
	}
	return i
}

func (p *parser) parseFloat() float64 {
	floatLit := p.lit
	if p.tok == scanner.Float || p.tok == scanner.Int {
		p.next()
	} else {
		p.expect(scanner.Float)
	}
	f, err := strconv.ParseFloat(floatLit, 64)
	if err != nil {
		p.error(p.pos, fmt.Sprintf("invalid number literal: %s", floatLit))
	}
	return f
}

func (p *parser) parseString() string {
	stringLit := p.lit
	p.expect(scanner.String)
	s, err := strconv.Unquote(stringLit)
	if err != nil {
		p.error(p.pos, fmt.Sprintf("invalid string literal: %s", stringLit))
	}
	return s
}

func (p *parser) parseBool() bool {
	boolLit := p.parseIdentifier()
	b, err := convertToBool(boolLit)
	if err != nil {
		p.error(p.pos, fmt.Sprintf("invalid boolean literal: %s", boolLit))
	}
	return b
}

func (p *parser) parseComment() string {
	var sb strings.Builder
	p.tok = p.scanner.Next()
	for p.tok != scanner.EOF && p.tok != '\n' {
		sb.WriteRune(p.tok)
		p.tok = p.scanner.Next()
	}
	p.scanner.Position = p.scanner.Pos()
	return sb.String()
}

func (p *parser) parseCode() string {
	var sb strings.Builder
	inString := false
	inEscapeSequence := false
	openSquareBrackets := 0
	var prevTok rune
	for p.tok != scanner.EOF && p.tok != '\n' {
		if !inString {
			if p.tok == ';' {
				// comment starts here
				break
			}
			if p.tok == '-' && p.scanner.Peek() == '>' {
				// goto starts here
				break
			}
			if p.tok == '[' && unicode.IsSpace(prevTok) && !unicode.IsSpace(p.scanner.Peek()) {
				// conditions start here after code of '!' statement
				break
			}
			if p.tok == ']' && openSquareBrackets == 0 {
				// condition ends here for code of condition
				break
			}
			if p.tok == '[' {
				openSquareBrackets++
			}
			if p.tok == ']' {
				openSquareBrackets--
			}
		}
		switch p.tok {
		case '\\':
			if inString && !inEscapeSequence {
				inEscapeSequence = true
			}
		case '"':
			if !inEscapeSequence {
				inString = !inString
			}
			fallthrough
		default:
			inEscapeSequence = false
		}
		sb.WriteRune(p.tok)
		prevTok = p.tok
		p.tok = p.scanner.Next()
	}
	p.scanner.Position = p.scanner.Pos()
	return strings.TrimSpace(sb.String())
}

func (p *parser) parseCondition() Condition {
	c := p.parseCode()
	switch c {
	case "once":
		return &condition.Once{}
	case "showonce":
		return &condition.ShowOnce{}
	case "onceever":
		return &condition.OnceEver{}
	case "temponce":
		return &condition.TempOnce{}
	}
	if isActorName(c) {
		return &condition.Actor{Actor: c}
	}
	return &condition.Code{Code: c}
}

func (p *parser) expectCommentOrNewLine() {
	if p.tok != ';' && p.tok != '\n' && p.tok != scanner.EOF {
		p.expect('\n')
		return
	}
	if p.tok == ';' {
		_ = p.parseComment()
	}
	if p.tok == '\n' {
		p.next()
	}
}

func (p *parser) expect(tok rune) scanner.Position {
	pos := p.pos
	if p.tok != tok {
		p.errorExpected(pos, scanner.TokenString(tok))
	}
	p.next() // make progress in any case
	return pos
}

func (p *parser) errorExpected(pos scanner.Position, msg string) {
	msg = "expected " + msg
	if pos.Offset == p.pos.Offset {
		// the error happened at the current position;
		// make the error message more specific
		msg += ", found " + scanner.TokenString(p.tok)
		if p.tok < 0 {
			msg += " " + p.lit
		}
	}
	p.error(pos, msg)
}

func (p *parser) error(pos scanner.Position, msg string) {
	p.errors = append(p.errors, newError(pos, msg))
}

func convertToBool(s string) (bool, error) {
	switch s {
	case "YES", "yes":
		return true, nil
	case "NO", "no":
		return false, nil
	}
	return false, fmt.Errorf("invalid boolean literal: '%s'", s)
}

func isActorName(ident string) bool {
	for i, r := range ident {
		if i == 0 && !unicode.IsLetter(r) {
			return false
		}
		if !unicode.IsLetter(r) && !unicode.IsNumber(r) {
			return false
		}
	}
	return ident != ""
}
