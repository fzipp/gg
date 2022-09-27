// Copyright 2020 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package savegame

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"

	"github.com/fzipp/gg/crypt/xxtea"
	"github.com/fzipp/gg/ggdict"
)

var key = xxtea.Key{
	0xAEA4EDF3,
	0xAFF8332A,
	0xB5A2DBB4,
	0x9B4BA022,
}

var endianness = binary.LittleEndian

const lenFooter = 16

func Load(path string) (map[string]any, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("could not open savegame file: %w", err)
	}
	defer f.Close()
	return Read(f)
}

func Read(r io.Reader) (map[string]any, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("could not read savegame data: %w", err)
	}
	decrypted := xxtea.Decrypt(data, key)
	if !isChecksumOk(decrypted) {
		return nil, fmt.Errorf("invalid checksum for savegame data")
	}
	dict, err := ggdict.Unmarshal(decrypted, false) // Todo: RtMI?
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal savegame data: %w", err)
	}
	return dict, nil
}

func Save(path string, dict map[string]any) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("could not open savegame file: %w", err)
	}
	defer f.Close()
	return Write(f, dict)
}

func Write(w io.Writer, dict map[string]any) error {
	data := ggdict.Marshal(dict, false) // Todo: RtMI?
	data = zeroPad(data, 500_000)
	sum := checksum(data)
	footerBytes := make([]byte, lenFooter)
	endianness.PutUint32(footerBytes, sum)
	data = append(data, footerBytes...)
	encrypted := xxtea.Encrypt(data, key)
	_, err := w.Write(encrypted)
	if err != nil {
		return fmt.Errorf("could not write savegame data: %w", err)
	}
	return nil
}

func zeroPad(data []byte, minLen int) []byte {
	if len(data) >= minLen {
		return data
	}
	lenPadding := minLen - len(data)
	return append(data, make([]byte, lenPadding)...)
}

func isChecksumOk(data []byte) bool {
	checksumIndex := len(data) - lenFooter
	sumGot := checksum(data[:checksumIndex])
	sumWant := endianness.Uint32(data[checksumIndex:])
	return sumGot == sumWant
}

func checksum(data []byte) uint32 {
	sum := uint32(0x6583463)
	for _, b := range data {
		sum += uint32(b)
	}
	return sum
}
