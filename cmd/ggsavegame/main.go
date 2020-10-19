// Copyright 2020 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// A tool to decrypt and encrypt Thimbleweed Park savegame files.
//
// Usage:
//     ggsavegame -to-json|-from-json savegame_file
//
// Flags:
//     -to-json    Converts the given savegame file to JSON format on
//                 standard output.
//     -from-json  Converts the given JSON file to savegame format on
//                 standard output. You might want to redirect it to a file,
//                 since it is a binary format.
//
// Examples:
//     ggsavegame -to-json Savegame1.save > Savegame1.json
//     ggsavegame -from-json Savegame1.json > Savegame1.save
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/fzipp/gg/savegame"
)

func usage() {
	fail(`A tool to decrypt and encrypt Thimbleweed Park savegame files.

Usage:
    ggsavegame -to-json|-from-json savegame_file

Flags:
    -to-json    Converts the given savegame file to JSON format on
                standard output.
    -from-json  Converts the given JSON file to savegame format on
                standard output. You might want to redirect it to a file,
                since it is a binary format.

Examples:
    ggsavegame -to-json Savegame1.save > Savegame1.json
    ggsavegame -from-json Savegame1.json > Savegame1.save`)
}

func main() {
	savegameFilePath := flag.String("to-json", "", "")
	jsonFilePath := flag.String("from-json", "", "")

	flag.Usage = usage
	flag.Parse()

	if *savegameFilePath == "" && *jsonFilePath == "" {
		usage()
	}
	if *savegameFilePath != "" && *jsonFilePath != "" {
		fail("-from-json and -to-json flags cannot be used together. See -help for more information.")
	}

	if *savegameFilePath != "" {
		toJSON(*savegameFilePath)
		return
	}

	if *jsonFilePath != "" {
		fromJSON(*jsonFilePath)
		return
	}
}

func toJSON(path string) {
	dict, err := savegame.Load(path)
	check(err)
	jsonData, err := json.MarshalIndent(dict, "", "  ")
	check(err)
	fmt.Println(string(jsonData))
}

func fromJSON(path string) {
	jsonData, err := ioutil.ReadFile(path)
	check(err)
	dict := make(map[string]interface{})
	err = json.Unmarshal(jsonData, &dict)
	check(err)
	err = savegame.Write(os.Stdout, dict)
	check(err)
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
