// Copyright 2020 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ggdict

import "encoding/binary"

type offsets []int
type valueType byte

const (
	typeNull = valueType(iota + 1)
	typeDictionary
	typeArray
	typeString
	typeInteger
	typeFloat
	typeOffsets
)

const formatSignature = 0x04030201

var byteOrder = binary.LittleEndian
