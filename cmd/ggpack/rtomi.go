// Copyright 2022 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/fzipp/gg/crypt/xor"
)

func loadKeyIfNecessary(key xor.Key, packFile string) {
	if !key.NeedsLoading() {
		return
	}
	execFile, err := locateExecFile(packFile)
	if err != nil {
		fail("Could not find game executable file. Please make sure that your pack file is located in the same directory as the game's executable.")
	}
	err = key.LoadFrom(execFile)
	if err != nil {
		fail("XOR key could not be loaded from the game's executable.")
	}
}

func locateExecFile(packFile string) (string, error) {
	packFile, err := filepath.Abs(packFile)
	if err != nil {
		return "", err
	}
	execFileNames := []string{
		"Return to Monkey Island.exe",
		"Return to Monkey Island",
	}
	for _, name := range execFileNames {
		execFile := filepath.Join(filepath.Dir(packFile), name)
		_, err = os.Stat(execFile)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				continue
			}
			return "", err
		}
		return execFile, nil
	}
	return "", errors.New("game executable file not found")
}
