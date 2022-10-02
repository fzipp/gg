// Copyright 2022 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ggdict

type Format struct {
	ShortStringIndices bool
	CoordinateTypes    bool
}

// FormatThimbleweed is the format found in Thimbleweed Park and Delores.
var FormatThimbleweed = Format{
	ShortStringIndices: false,
	CoordinateTypes:    false,
}

// FormatMonkey is the format found in Return to Monkey Island.
var FormatMonkey = Format{
	ShortStringIndices: true,
	CoordinateTypes:    true,
}
