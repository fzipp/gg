// Copyright 2020 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package wimpy

import (
	"fmt"
	"image"
	"strconv"
	"strings"
)

// "{213,118}"
func parsePoint(s string) (image.Point, error) {
	var x, y int
	_, err := fmt.Sscanf(s, "{%d,%d}", &x, &y)
	if err != nil {
		return image.Point{}, err
	}
	return image.Pt(x, y), nil
}

// "{0.75,0.5}"
func parsePointF(s string) (PointF, error) {
	var x, y float64
	_, err := fmt.Sscanf(s, "{%g,%g}", &x, &y)
	if err != nil {
		return PointF{}, err
	}
	return PointF{X: x, Y: y}, nil
}

// "{{-23,-20},{17,20}}"
func parseRectangle(s string) (image.Rectangle, error) {
	var x1, y1, x2, y2 int
	_, err := fmt.Sscanf(s, "{{%d,%d},{%d,%d}}", &x1, &y1, &x2, &y2)
	if err != nil {
		return image.Rectangle{}, err
	}
	return image.Rect(x1, y1, x2, y2), nil
}

// "DIR_LEFT"
func parseDirection(s string) (Direction, error) {
	switch s {
	case "DIR_FRONT":
		return DirFront, nil
	case "DIR_BACK":
		return DirBack, nil
	case "DIR_LEFT":
		return DirLeft, nil
	case "DIR_RIGHT":
		return DirRight, nil
	}
	return Direction(0), fmt.Errorf("unknown direction: %q'", s)
}

// "{82,94};{134,94};{142,91};{174,91};{183,94}"
func parsePolygon(s string) ([]image.Point, error) {
	elements := strings.Split(s, ";")
	polygon := make([]image.Point, len(elements))
	for i, element := range elements {
		v, err := parsePoint(element)
		if err != nil {
			return nil, err
		}
		polygon[i] = v
	}
	return polygon, nil
}

// "1.2@68"
func parseScaling(s string) (Scaling, error) {
	parts := strings.Split(s, "@")
	if len(parts) != 2 {
		return Scaling{}, fmt.Errorf("unknown scaling format: %q", s)
	}
	scaleFactor, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return Scaling{}, err
	}
	at, err := strconv.Atoi(parts[1])
	if err != nil {
		return Scaling{}, err
	}
	return Scaling{Factor: scaleFactor, At: at}, nil
}
