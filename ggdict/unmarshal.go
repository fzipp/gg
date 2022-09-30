// Copyright 2020 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ggdict

import (
	"errors"
	"fmt"
	"strconv"
)

func Unmarshal(data []byte, shortStringIndices bool) (map[string]any, error) {
	u := &unmarshaller{buf: data, shortStringIndices: shortStringIndices}

	signature := u.readRawInt32()
	if signature != formatSignature {
		return nil, fmt.Errorf("invalid format signature: %#x", signature)
	}

	// Unused, as far as known. Always 1. Maybe format version?
	_ = u.readRawInt32()

	offsetIndexStart := u.readRawInt32()
	ou := &unmarshaller{
		buf:                data,
		offset:             offsetIndexStart,
		shortStringIndices: shortStringIndices,
	}
	offsetIndex, err := ou.readValue()
	if err != nil {
		return nil, fmt.Errorf("could not read offset index: %w", err)
	}

	offsets, ok := offsetIndex.(offsets)
	if !ok {
		return nil, errors.New("read value is not an offset index")
	}
	u.offsetIndex = offsets
	root, err := u.readValue()
	if err != nil {
		return nil, fmt.Errorf("could not read root: %w", err)
	}
	dict, ok := root.(map[string]any)
	if !ok {
		return nil, errors.New("root is not a dictionary")
	}
	return dict, nil
}

type unmarshaller struct {
	buf                []byte
	offset             int
	offsetIndex        offsets
	shortStringIndices bool
}

func (u *unmarshaller) readValue() (any, error) {
	switch valueType := u.readTypeMarker(); valueType {
	case typeNull:
		return nil, nil
	case typeDictionary:
		return u.readDictionary()
	case typeArray:
		return u.readArray()
	case typeString, typeCoordinate, typeCoordinateList, typeHotspot:
		return u.readString(), nil
	case typeInteger:
		return u.readInteger()
	case typeFloat:
		return u.readFloat()
	case typeOffsets:
		return u.readOffsets(), nil
	default:
		return nil, fmt.Errorf("unknown value type: %d", valueType)
	}
}

func (u *unmarshaller) readTypeMarker() valueType {
	return valueType(u.readByte())
}

func (u *unmarshaller) readDictionary() (map[string]any, error) {
	length := u.readRawInt32()
	dictionary := make(map[string]any, length)
	for i := 0; i < length; i++ {
		key := u.readString()
		value, err := u.readValue()
		if err != nil {
			return nil, fmt.Errorf("could not read dictionary value for key %q: %w", key, err)
		}
		dictionary[key] = value
	}
	if u.readTypeMarker() != typeDictionary {
		return nil, fmt.Errorf("unterminated dictionary")
	}
	return dictionary, nil
}

func (u *unmarshaller) readArray() ([]any, error) {
	length := u.readRawInt32()
	array := make([]any, length)
	for i := 0; i < length; i++ {
		value, err := u.readValue()
		if err != nil {
			return nil, fmt.Errorf("could not read array value: %w", err)
		}
		array[i] = value
	}
	if u.readTypeMarker() != typeArray {
		return nil, fmt.Errorf("unterminated array")
	}
	return array, nil
}

func (u *unmarshaller) readString() string {
	var stringIndex int

	if u.shortStringIndices {
		stringIndex = u.readRawInt16()
	} else {
		stringIndex = u.readRawInt32()
	}

	startOffset := u.offsetIndex[stringIndex]
	endOffset := startOffset
	for endOffset < len(u.buf) && u.buf[endOffset] != 0 {
		endOffset++
	}
	return string(u.buf[startOffset:endOffset])
}

func (u *unmarshaller) readInteger() (int, error) {
	return strconv.Atoi(u.readString())
}

func (u *unmarshaller) readFloat() (float64, error) {
	return strconv.ParseFloat(u.readString(), 64)
}

func (u *unmarshaller) readOffsets() offsets {
	var offsets offsets
	for byteOrder.Uint32(u.buf[u.offset:]) != 0xFFFFFFFF {
		offsets = append(offsets, u.readRawInt32())
	}
	return offsets
}

func (u *unmarshaller) readRawInt32() int {
	i := int(byteOrder.Uint32(u.buf[u.offset:]))
	u.offset += 4
	return i
}

func (u *unmarshaller) readRawInt16() int {
	i := int(byteOrder.Uint16(u.buf[u.offset:]))
	u.offset += 2
	return i
}

func (u *unmarshaller) readByte() byte {
	b := u.buf[u.offset]
	u.offset++
	return b
}
