// Copyright 2020 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// A tool to convert back and forth between the GGDictionary format and JSON.
//
// The GGDictionary format is used by the Thimbleweed Park point-and-click
// adventure game engine. It encodes a key-value data structure like JSON, but as
// a binary format. For example, for Thimbleweed Park *.wimpy and *Animation.json
// files are stored in this format within a "ggpack" file.
//
// Usage:
//
//	ggdict -to-json|-from-json [-monkey-island] path
//
// Flags:
//
//	-to-json        Converts the given GGDictionary file to JSON format on
//	                standard output.
//	-from-json      Converts the given JSON file to GGDictionary format on
//	                standard output. You might want to redirect it to a file,
//	                since it is a binary format.
//	-monkey-island  Use the new format used in Return to Monkey Island.
//
// Examples:
//
//	ggdict -to-json Example.wimpy > Example.wimpy.json
//	ggdict -from-json Example.wimpy.json > Example.wimpy
//	ggdict -to-json ExampleAnimation.json > ExampleAnimation.really.json
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/fzipp/gg/ggdict"
)

func usage() {
	fail(`A tool to convert back and forth between the GGDictionary format and JSON.

The GGDictionary format encodes a key-value data structure like JSON, but as
a binary format. For example, for Thimbleweed Park *.wimpy and *Animation.json
files are stored in this format within a "ggpack" file.

Usage:
    ggdict -to-json|-from-json [-monkey-island] path

Flags:
    -to-json        Converts the given GGDictionary file to JSON format on
                    standard output.
    -from-json      Converts the given JSON file to GGDictionary format on
                    standard output. You might want to redirect it to a file,
                    since it is a binary format.
    -monkey-island  Use the new format used in Return to Monkey Island.

Examples:
    ggdict -to-json Example.wimpy > Example.wimpy.json
    ggdict -from-json Example.wimpy.json > Example.wimpy
    ggdict -to-json ExampleAnimation.json > ExampleAnimation.really.json`)
}

func main() {
	ggdictFilePath := flag.String("to-json", "", "")
	jsonFilePath := flag.String("from-json", "", "")
	monkeyIslandMode := flag.Bool("monkey-island", false, "")

	flag.Usage = usage
	flag.Parse()

	if *ggdictFilePath == "" && *jsonFilePath == "" {
		usage()
	}
	if *ggdictFilePath != "" && *jsonFilePath != "" {
		fail("-from-json and -to-json flags cannot be used together. See -help for more information.")
	}

	var format ggdict.Format
	if *monkeyIslandMode {
		format = ggdict.FormatMonkey
	} else {
		format = ggdict.FormatThimbleweed
	}

	if *ggdictFilePath != "" {
		toJSON(*ggdictFilePath, format)
		return
	}

	if *jsonFilePath != "" {
		fromJSON(*jsonFilePath, format)
		return
	}
}

func toJSON(path string, format ggdict.Format) {
	buf, err := os.ReadFile(path)
	check(err)
	dict, err := ggdict.Unmarshal(buf, format)
	check(err)
	jsonData, err := json.MarshalIndent(dict, "", "  ")
	check(err)
	fmt.Println(string(jsonData))
}

func fromJSON(path string, format ggdict.Format) {
	jsonData, err := os.ReadFile(path)
	check(err)
	dict := make(map[string]any)
	err = json.Unmarshal(jsonData, &dict)
	check(err)
	_, err = os.Stdout.Write(ggdict.Marshal(dict, format))
	check(err)
}

func check(err error) {
	if err != nil {
		fail(err)
	}
}

func fail(message any) {
	_, _ = fmt.Fprintln(os.Stderr, message)
	os.Exit(1)
}
