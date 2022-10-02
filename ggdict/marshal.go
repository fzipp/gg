// Copyright 2020 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ggdict

import (
	"sort"
	"strconv"
)

func Marshal(dict map[string]any, f Format) []byte {
	m := newMarshaller(f)
	m.writeRawUint32(formatSignature)
	m.writeRawUint32(1)
	m.writeRawUint32(0)
	m.writeValue(dict)
	m.writeStringOffsets()
	m.writeStrings()
	return m.buf
}

type marshaller struct {
	buf           []byte
	offset        int
	strings       []string
	stringIndices map[string]int
	format        Format
}

func newMarshaller(f Format) *marshaller {
	return &marshaller{
		stringIndices: make(map[string]int),
		format:        f,
	}
}

func (m *marshaller) writeValue(value any) {
	switch v := value.(type) {
	case nil:
		m.writeNull()
	case map[string]any:
		m.writeDictionary(v)
	case []any:
		m.writeArray(v)
	case string:
		m.writeString(v)
	case int:
		m.writeInteger(v)
	case int32:
		m.writeInteger(int(v))
	case int64:
		m.writeInteger(int(v))
	case uint32:
		m.writeInteger(int(v))
	case uint64:
		m.writeInteger(int(v))
	case float64:
		m.writeFloat(v)
	case float32:
		m.writeFloat(float64(v))
	}
}

func (m *marshaller) writeTypeMarker(t valueType) {
	m.writeRawByte(byte(t))
}

func (m *marshaller) writeNull() {
	m.writeTypeMarker(typeNull)
}

func (m *marshaller) writeDictionary(d map[string]any) {
	m.writeTypeMarker(typeDictionary)
	m.writeRawUint32(len(d))
	// sorted keys for reproducible results
	for _, k := range sortedKeys(d) {
		m.writeStringIndex(k)
		m.writeValue(d[k])
	}
	m.writeTypeMarker(typeDictionary)
}

func sortedKeys(m map[string]any) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func (m *marshaller) writeArray(a []any) {
	m.writeTypeMarker(typeArray)
	m.writeRawUint32(len(a))
	for _, v := range a {
		m.writeValue(v)
	}
	m.writeTypeMarker(typeArray)
}

func (m *marshaller) writeString(s string) {
	m.writeTypeMarker(typeString)
	m.writeStringIndex(s)
}

func (m *marshaller) writeInteger(i int) {
	m.writeTypeMarker(typeInteger)
	m.writeStringIndex(strconv.Itoa(i))
}

func (m *marshaller) writeFloat(f float64) {
	m.writeTypeMarker(typeFloat)
	m.writeStringIndex(strconv.FormatFloat(f, 'g', -1, 64))
}

func (m *marshaller) writeStringIndex(s string) {
	idx, ok := m.stringIndices[s]
	if !ok {
		idx = len(m.strings)
		m.stringIndices[s] = idx
		m.strings = append(m.strings, s)
	}
	if m.format.ShortStringIndices {
		m.writeRawUint16(idx)
	} else {
		m.writeRawUint32(idx)
	}
}

func (m *marshaller) writeStringOffsets() {
	m.writeStringOffsetsStart(m.offset)
	m.writeTypeMarker(typeStringOffsets)
	strOffset := m.offset
	lengths := make([]int, len(m.strings))
	for i, key := range m.strings {
		lengths[i] = len(key) + 1
		strOffset += 4
	}
	strOffset += 5
	for _, length := range lengths {
		m.writeRawUint32(strOffset)
		strOffset += length
	}
	m.writeRawUint32(0xFFFFFFFF)
}

func (m *marshaller) writeStringOffsetsStart(offset int) {
	byteOrder.PutUint32(m.buf[8:], uint32(offset))
}

func (m *marshaller) writeStrings() {
	m.writeTypeMarker(typeStrings)
	for _, s := range m.strings {
		m.writeRawString(s)
	}
}

func (m *marshaller) writeRawString(s string) {
	m.writeRawBytes(append([]byte(s), 0))
}

func (m *marshaller) writeRawUint32(i int) {
	b := make([]byte, 4)
	byteOrder.PutUint32(b, uint32(i))
	m.writeRawBytes(b)
}

func (m *marshaller) writeRawUint16(i int) {
	b := make([]byte, 2)
	byteOrder.PutUint16(b, uint16(i))
	m.writeRawBytes(b)
}

func (m *marshaller) writeRawBytes(b []byte) {
	m.buf = append(m.buf, b...)
	m.offset += len(b)
}

func (m *marshaller) writeRawByte(b byte) {
	m.buf = append(m.buf, b)
	m.offset++
}
