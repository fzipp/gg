// Copyright 2020 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package wimpy

import "image"

type Room struct {
	Name       string
	Sheet      string
	Background []string
	Fullscreen int
	Height     int
	Layers     []Layer
	Objects    []Object
	RoomSize   image.Point
	Scalings   Scalings
	WalkBoxes  []WalkBox
}

type Layer struct {
	Name     []string
	Parallax PointF
	ZSort    int
}

type Object struct {
	Name       string
	Parent     string
	Animations []Animation
	HotSpot    image.Rectangle
	Pos        image.Point
	UseDir     Direction
	UsePos     image.Point
	ZSort      int
	Prop       bool
	Spot       bool
	Trigger    bool
}

type Animation struct {
	Name     string
	FPS      float64
	Triggers []string
	Frames   []string
}

type WalkBox struct {
	Name    string
	Polygon []image.Point
}

type Scalings struct {
	Scaling []Scaling
	Trigger string
}

type Scaling struct {
	Factor float64
	At     int
}

type PointF struct {
	X, Y float64
}

type Direction int

const (
	DirRight = Direction(1 << iota)
	DirLeft
	DirFront
	DirBack
)
