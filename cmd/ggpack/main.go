// Copyright 2020 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// A tool to inspect, unpack or create "ggpack" files.
//
// Usage:
//     ggpack -list|-extract|-create "filename_pattern" [-key name] ggpack_file
//
// Flags:
//     -list     List files in the pack matching the pattern.
//     -extract  Extract the files from the pack matching the pattern to
//               the current working directory.
//     -create   Create a new pack and add the files from the file system
//               matching the pattern.
//     -key      Name of the key to decrypt/encrypt the data via XOR.
//               Possible names: 56ad (default), 5bad, 566d, 5b6d
//
// Examples:
//     ggpack -list "*" MyPackage.ggpack1
//     ggpack -list "*.tsv" MyPackage.ggpack1
//     ggpack -extract "ExampleSheet.png" MyPackage.ggpack1
//     ggpack -extract "*.txt" MyPackage.ggpack1
//     ggpack -extract "*" MyPackage.ggpack1
//     ggpack -create "*" MyPackage.ggpack1`
package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/fzipp/gg/crypt/xor"
	"github.com/fzipp/gg/ggpack"
)

func usage() {
	fail(`A tool to inspect, unpack or create "ggpack" files.

Usage:
    ggpack -list|-extract|-create "filename_pattern" ggpack_file

Flags:
    -list     List files in the pack matching the pattern.
    -extract  Extract the files from the pack matching the pattern to
              the current working directory.
    -create   Create a new pack and add the files from the file system
              matching the pattern.
    -key      Name of the key to decrypt/encrypt the data via XOR.
              Possible names: 56ad (default), 5bad, 566d, 5b6d

Examples:
    ggpack -list "*" MyPackage.ggpack1
    ggpack -list "*.tsv" MyPackage.ggpack1
    ggpack -extract "ExampleSheet.png" MyPackage.ggpack1
    ggpack -extract "*.txt" MyPackage.ggpack1
    ggpack -extract "*" MyPackage.ggpack1
    ggpack -create "*" MyPackage.ggpack1`)
}

func main() {
	listPattern := flag.String("list", "", "List files in the pack matching the pattern.")
	extractPattern := flag.String("extract", "", "Extract the files from the pack matching the pattern to the current working directory.")
	createPattern := flag.String("create", "", "Create a new pack and add the files from the file system matching the pattern.")
	keyName := flag.String("key", "56ad", "Name of the key to decrypt/encrypt the data via XOR.")

	flag.Usage = usage
	flag.Parse()

	if flag.NArg() < 1 {
		usage()
		return
	}
	if flag.NArg() > 1 {
		fmt.Println(flag.Args())
		fail("Please specify only one pack_file argument. See -help for more information.")
		return
	}
	packFile := flag.Arg(0)

	patternFlags := []string{*listPattern, *extractPattern, *createPattern}
	var patterns []string
	for _, pattern := range patternFlags {
		if pattern != "" {
			patterns = append(patterns, pattern)
		}
	}

	if len(patterns) == 0 {
		fail("Please choose an operation via flag. See -help for more information.")
		return
	}

	if len(patterns) > 1 {
		fail("Please use only one operation flag, not multiple at the same time. See -help for more information.")
		return
	}

	pattern := patterns[0]
	key, ok := xor.KnownKeys[strings.ToLower(*keyName)]
	if !ok {
		fail("Unknown XOR key name: \"" + *keyName + "\"")
	}

	if *createPattern != "" {
		paths, err := filepath.Glob(pattern)
		check(err)
		err = create(packFile, paths, key)
		check(err)
		return
	}

	pack, err := ggpack.OpenUsingKey(packFile, key)
	check(err)
	defer pack.Close()

	filenames, err := filterFilenames(pack.List(), pattern)
	check(err)

	if *listPattern != "" {
		list(filenames)
	}
	if *extractPattern != "" {
		extractAll(pack, filenames)
	}
}

func filterFilenames(entries []ggpack.DirectoryEntry, pattern string) ([]string, error) {
	var filtered []string
	for _, entry := range entries {
		matches, err := filepath.Match(pattern, entry.Filename)
		if err != nil {
			return nil, fmt.Errorf("invalid filename pattern: %s", pattern)
		}
		if matches {
			filtered = append(filtered, entry.Filename)
		}
	}
	return filtered, nil
}

func list(filenames []string) {
	for _, filename := range filenames {
		fmt.Println(filename)
	}
}

func extractAll(pack *ggpack.Pack, filenames []string) {
	for _, filename := range filenames {
		extract(pack, filename)
	}
}

func extract(pack *ggpack.Pack, filename string) {
	packFile, _, err := pack.File(filename)
	check(err)
	diskFile, err := os.Create(filename)
	check(err)
	defer diskFile.Close()
	_, err = io.Copy(diskFile, packFile)
	check(err)
	err = diskFile.Sync()
	check(err)
}

func create(packFilePath string, paths []string, key *xor.Key) error {
	packFile, err := os.Create(packFilePath)
	if err != nil {
		return fmt.Errorf("could not create pack file: %w", err)
	}
	defer packFile.Close()
	packer, err := ggpack.NewPacker(packFile)
	if err != nil {
		return fmt.Errorf("could not initialize pack file: %w", err)
	}
	packer.SetKey(key)
	err = packer.WriteFiles(paths)
	if err != nil {
		return fmt.Errorf("could not write files to pack file: %w", err)
	}
	return packer.Finish()
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
