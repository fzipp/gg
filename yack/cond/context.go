// Copyright 2020 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cond

type Context interface {
	IsActor(actor string) bool
	IsCodeTrue(code string) bool
	IsOnce() bool
	IsShowOnce() bool
	IsOnceEver() bool
	IsShowOnceEver() bool
	IsTempOnce() bool
}
