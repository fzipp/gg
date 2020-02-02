// Copyright 2020 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ggdict

import (
	"errors"
	"fmt"
	"strconv"
)

func Unmarshal(data []byte) (map[string]interface{}, error) {
	u := &unmarshaller{buf: data}

	signature := u.readRawInt()
	if signature != formatSignature {
		return nil, fmt.Errorf("invalid format signature: %x", signature)
	}
	_ = u.readRawInt() // Unused, as far as known

	offsetIndexStart := u.readRawInt()
	ou := &unmarshaller{
		buf: data,
		offset: offsetIndexStart,
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
	dict, ok := root.(map[string]interface{})
	if !ok {
		return nil, errors.New("root is not a dictionary")
	}
	return dict, nil
}

type unmarshaller struct {
	buf         []byte
	offset      int
	offsetIndex offsets
}

func (u *unmarshaller) readValue() (interface{}, error) {
	switch valueType := u.readTypeMarker(); valueType {
	case typeNull:
		return nil, nil
	case typeDictionary:
		return u.readDictionary()
	case typeArray:
		return u.readArray()
	case typeString:
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

func (u *unmarshaller) readDictionary() (map[string]interface{}, error) {
	length := u.readRawInt()
	dictionary := make(map[string]interface{}, length)
	for i := 0; i < length; i++ {
		key := u.readString()
		value, err := u.readValue()
		if err != nil {
			return nil, fmt.Errorf("could not read dictionary value for key \"%s\": %w", key, err)
		}
		dictionary[key] = value
	}
	if u.readTypeMarker() != typeDictionary {
		return nil, fmt.Errorf("unterminated dictionary")
	}
	return dictionary, nil
}

func (u *unmarshaller) readArray() ([]interface{}, error) {
	length := u.readRawInt()
	array := make([]interface{}, length)
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
	startOffset := u.offsetIndex[u.readRawInt()]
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
		offsets = append(offsets, u.readRawInt())
	}
	return offsets
}

func (u *unmarshaller) readRawInt() int {
	i := int(byteOrder.Uint32(u.buf[u.offset:]))
	u.offset += 4
	return i
}

func (u *unmarshaller) readByte() byte {
	b := u.buf[u.offset]
	u.offset++
	return b
}
