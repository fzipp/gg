// Copyright 2020 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package yack reads and runs yack dialogs.
//
// The EBNF grammar for yack files:
//
//	SourceFile = { Line } .
//
//	Line = [ (LabelTag | ConditionalStatement) ] [ comment ] newline .
//
//	LabelTag = ":" label .
//
//	ConditionalStatement = Statement { Condition } .
//
//	Condition = "[" ("once" | "showonce" | "onceever" | "showonceever" |
//	                 "temponce" | actor_name | Code) "]" .
//
//	Statement = SayStmt | GotoStmt | ExecuteStmt | ChoiceStmt | ShutUpStmt |
//	            PauseStmt | WaitForStmt | WaitWhileStmt | ParrotStmt | DialogStmt |
//	            OverrideStmt | AllowObjectsStmt | LimitStmt .
//
//	SayStmt          = actor_name ":" string_lit [ GotoStmt ] .
//	GotoStmt         = "->" label .
//	ExecuteStmt      = "!" Code .
//	ChoiceStmt       = int_lit (string_lit | ("$" Code)) GotoStmt .
//	ShutUpStmt       = "shutup" .
//	PauseStmt        = "pause" float_lit .
//	WaitForStmt      = "waitfor" [ actor_name ] .
//	WaitWhileStmt    = "waitwhile" Code .
//	ParrotStmt       = "parrot" bool_lit .
//	DialogStmt       = "dialog" actor_name .
//	OverrideStmt     = "override" label .
//	AllowObjectsStmt = "allowobjects" bool_lit .
//	LimitStmt        = "limit" int_lit .
//
//	Code = /* Scripting language expression, e.g. Squirrel */ .
//
//	comment        = ";" { unicode_char } .
//	actor_name     = unicode_letter { unicode_letter | unicode_digit } .
//	label          = letter_uscore { letter_uscore | unicode_digit } .
//	letter_uscore  = unicode_letter | "_" .
//	bool_lit       = "yes" | "YES" | "no" | "NO" .
//	int_lit        = decimal_digits .
//	float_lit      = decimal_digits | decimal_digits "." [ decimal_digits ] |
//	                 "." decimal_digits .
//	decimal_digits = { decimal_digit } .
//	decimal_digit  = "0" … "9" .
//	string_lit     = `"` { unicode_value } `"` .
//	unicode_value  = unicode_char | escaped_char .
//	escaped_char   = `\` ( "n" | `\` | `"` ) .
//
//	newline        = /* the Unicode code point U+000A */ .
//	unicode_char   = /* an arbitrary Unicode code point except newline */ .
//	unicode_letter = /* a Unicode code point classified as "Letter" */ .
//	unicode_digit  = /* a Unicode code point classified as "Number, decimal digit" */ .
package yack
