// Copyright 2020 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ggdict

import (
	"errors"
	"fmt"
	"strconv"
)

func Unmarshal(data []byte, f Format) (map[string]any, error) {
	u := &unmarshaller{
		buf:    data,
		format: f,
	}

	signature := u.readRawUint32()
	if signature != formatSignature {
		return nil, fmt.Errorf("invalid format signature: %#x", signature)
	}

	// Unused, as far as known. Always 1. Maybe format version?
	_ = u.readRawUint32()

	stringOffsetsStart := u.readRawUint32()
	ou := &unmarshaller{
		buf:    data,
		offset: stringOffsetsStart,
		format: f,
	}
	stringOffsets, err := ou.readValue()
	if err != nil {
		return nil, fmt.Errorf("could not read string offsets: %w", err)
	}

	offs, ok := stringOffsets.(offsets)
	if !ok {
		return nil, errors.New("read value is not a string offsets table")
	}
	u.stringOffsets = offs
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
	buf           []byte
	offset        int
	stringOffsets offsets
	format        Format
}

func (u *unmarshaller) readValue() (any, error) {
	switch valueType := u.readTypeMarker(); valueType {
	case typeNull:
		return nil, nil
	case typeDictionary:
		return u.readDictionary()
	case typeArray:
		return u.readArray()
	case typeString, typeCoordinate, typeCoordinatePair, typeCoordinateList:
		return u.readString(), nil
	case typeInteger:
		return u.readInteger()
	case typeFloat:
		return u.readFloat()
	case typeStringOffsets:
		return u.readStringOffsets(), nil
	default:
		return nil, fmt.Errorf("unknown value type: %d", valueType)
	}
}

func (u *unmarshaller) readTypeMarker() valueType {
	return valueType(u.readRawByte())
}

func (u *unmarshaller) readDictionary() (map[string]any, error) {
	length := u.readRawUint32()
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
	length := u.readRawUint32()
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
	var strIndex int
	if u.format.ShortStringIndices {
		strIndex = u.readRawUint16()
	} else {
		strIndex = u.readRawUint32()
	}
	startOffset := u.stringOffsets[strIndex]
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

func (u *unmarshaller) readStringOffsets() offsets {
	var offs offsets
	for byteOrder.Uint32(u.buf[u.offset:]) != 0xFFFFFFFF {
		offs = append(offs, u.readRawUint32())
	}
	return offs
}

func (u *unmarshaller) readRawUint32() int {
	i := int(byteOrder.Uint32(u.buf[u.offset:]))
	u.offset += 4
	return i
}

func (u *unmarshaller) readRawUint16() int {
	i := int(byteOrder.Uint16(u.buf[u.offset:]))
	u.offset += 2
	return i
}

func (u *unmarshaller) readRawByte() byte {
	b := u.buf[u.offset]
	u.offset++
	return b
}
