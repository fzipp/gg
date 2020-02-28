// Copyright 2020 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ggdict

import (
	"sort"
	"strconv"
)

func Marshal(dict map[string]interface{}) []byte {
	m := newMarshaller()
	m.writeRawInt(formatSignature)
	m.writeRawInt(1)
	m.writeRawInt(0)
	m.writeValue(dict)
	m.writeKeys()
	return m.buf
}

type marshaller struct {
	buf      []byte
	offset   int
	keys     []string
	keyIndex map[string]int
}

func newMarshaller() *marshaller {
	return &marshaller{keyIndex: make(map[string]int)}
}

func (m *marshaller) writeValue(value interface{}) {
	switch v := value.(type) {
	case nil:
		m.writeNull()
	case map[string]interface{}:
		m.writeDictionary(v)
	case []interface{}:
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
	m.writeByte(byte(t))
}

func (m *marshaller) writeNull() {
	m.writeTypeMarker(typeNull)
}

func (m *marshaller) writeDictionary(d map[string]interface{}) {
	// sorted keys for reproducible results
	keys := make([]string, 0, len(d))
	for k := range d {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	m.writeTypeMarker(typeDictionary)
	m.writeRawInt(len(d))
	for _, k := range keys {
		m.writeKeyIndex(k)
		m.writeValue(d[k])
	}
	m.writeTypeMarker(typeDictionary)
}

func (m *marshaller) writeArray(a []interface{}) {
	m.writeTypeMarker(typeArray)
	m.writeRawInt(len(a))
	for _, v := range a {
		m.writeValue(v)
	}
	m.writeTypeMarker(typeArray)
}

func (m *marshaller) writeString(s string) {
	m.writeTypeMarker(typeString)
	m.writeKeyIndex(s)
}

func (m *marshaller) writeInteger(i int) {
	m.writeTypeMarker(typeInteger)
	m.writeKeyIndex(strconv.Itoa(i))
}

func (m *marshaller) writeFloat(f float64) {
	m.writeTypeMarker(typeFloat)
	m.writeKeyIndex(strconv.FormatFloat(f, 'g', -1, 64))
}

func (m *marshaller) writeKeyIndex(key string) {
	offset, ok := m.keyIndex[key]
	if !ok {
		offset = len(m.keys)
		m.keyIndex[key] = offset
		m.keys = append(m.keys, key)
	}
	m.writeRawInt(offset)
}

func (m *marshaller) writeKeys() {
	byteOrder.PutUint32(m.buf[8:], uint32(m.offset))

	m.writeTypeMarker(typeOffsets)
	keyOffset := m.offset
	lengths := make([]int, len(m.keys))
	for i, key := range m.keys {
		lengths[i] = len(key) + 1
		keyOffset += 4
	}
	keyOffset += 5
	for _, length := range lengths {
		m.writeRawInt(keyOffset)
		keyOffset += length
	}
	m.writeRawInt(0xFFFFFFFF)

	m.writeByte(0x8)
	for _, key := range m.keys {
		m.writeKey(key)
	}
}

func (m *marshaller) writeKey(key string) {
	m.buf = append(m.buf, []byte(key)...)
	m.offset += len(key)
	m.writeByte(0)
}

func (m *marshaller) writeRawInt(i int) {
	intBytes := make([]byte, 4)
	byteOrder.PutUint32(intBytes, uint32(i))
	m.buf = append(m.buf, intBytes...)
	m.offset += len(intBytes)
}

func (m *marshaller) writeByte(b byte) {
	m.buf = append(m.buf, b)
	m.offset++
}
