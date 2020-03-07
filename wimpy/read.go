// Copyright 2020 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package wimpy

import (
	"bytes"
	"fmt"
	"io"
	"runtime"

	"github.com/fzipp/gg/ggdict"
)

func Read(r io.Reader) (*Room, error) {
	var buf bytes.Buffer
	_, err := io.Copy(&buf, r)
	if err != nil {
		return nil, fmt.Errorf("could not read wimpy data: %w", err)
	}
	dict, err := ggdict.Unmarshal(buf.Bytes())
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal wimpy dictionary: %w", err)
	}
	return fromDict(dict)
}

func fromDict(dict map[string]interface{}) (r *Room, err error) {
	defer func() {
		msg := recover()
		if msg == nil {
			return
		}
		var ok bool
		if err, ok = msg.(*runtime.TypeAssertionError); !ok {
			panic(msg)
		}
	}()

	r = &Room{}
	r.Name = dict["name"].(string)
	switch bg := dict["background"].(type) {
	case string:
		r.Background = []string{bg}
	case []string:
		r.Background = bg
	}
	r.Fullscreen = optionalInt(dict["fullscreen"])
	r.Height = optionalInt(dict["height"])
	r.Layers = make([]Layer, 0)
	for i, layerDict := range optionalDicts(dict["layers"]) {
		layer := Layer{}
		switch name := layerDict["name"].(type) {
		case string:
			layer.Name = []string{name}
		case []string:
			layer.Name = name
		}
		switch parallax := layerDict["parallax"].(type) {
		case string:
			p, err := parsePointFloat(parallax)
			if err != nil {
				return nil, fmt.Errorf("room %q, layer %q [%d]: invalid parallax string", r.Name, layer.Name, i)
			}
			layer.Parallax = p
		case float64:
			layer.Parallax = PointFloat{X: parallax, Y: 1}
		}
		layer.ZSort = layerDict["zsort"].(int)
		r.Layers = append(r.Layers, layer)
	}
	r.Objects = make([]Object, 0)
	for i, o := range dict["objects"].([]interface{}) {
		obj := Object{}
		objDict := o.(map[string]interface{})
		obj.Name = objDict["name"].(string)
		obj.Parent = optionalString(objDict["parent"])
		obj.Animations = make([]Animation, 0)
		for _, animDict := range optionalDicts(objDict["animations"]) {
			anim := Animation{}
			anim.FPS = optionalFloat(animDict["fps"])
			anim.Triggers = optionalStrings(animDict["triggers"])
			anim.Frames = optionalStrings(animDict["triggers"])
			anim.Name = animDict["name"].(string)
			obj.Animations = append(obj.Animations, anim)
		}
		obj.HotSpot, err = parseRectangle(objDict["hotspot"].(string))
		if err != nil {
			return nil, fmt.Errorf("room %q, object %q [%d]: invalid hotspot rectangle", r.Name, obj.Name, i)
		}
		obj.Pos, err = parsePoint(objDict["pos"].(string))
		if err != nil {
			return nil, fmt.Errorf("room %q, object %q [%d]: invalid pos", r.Name, obj.Name, i)
		}
		obj.UseDir, err = parseDirection(objDict["usedir"].(string))
		if err != nil {
			return nil, fmt.Errorf("room %q, object %q [%d]: invalid usedir", r.Name, obj.Name, i)
		}
		obj.UsePos, err = parsePoint(objDict["usepos"].(string))
		if err != nil {
			return nil, fmt.Errorf("room %q, object %q [%d]: invalid usepos", r.Name, obj.Name, i)
		}
		obj.ZSort = objDict["zsort"].(int)
		obj.Prop = optionalBool(objDict["prop"])
		obj.Spot = optionalBool(objDict["spot"])
		obj.Trigger = optionalBool(objDict["trigger"])
		r.Objects = append(r.Objects, obj)
	}
	r.RoomSize, err = parsePoint(dict["roomsize"].(string))
	if err != nil {
		return nil, fmt.Errorf("room %q: invalid roomsize", r.Name)
	}
	scalings := Scalings{}
	for _, sc := range optionalSlice(dict["scaling"]) {
		var scaling Scaling
		switch sc := sc.(type) {
		case map[string]interface{}:
			for i, s := range optionalStrings(sc["scaling"]) {
				scaling, err = parseScaling(s)
				if err != nil {
					return nil, fmt.Errorf("room %q, scaling [%d]: invalid scaling", r.Name, i)
				}
				scalings.Scaling = append(scalings.Scaling, scaling)
			}
			scalings.Trigger = optionalString(sc["trigger"])
		case string:
			scaling, err = parseScaling(sc)
			if err != nil {
				return nil, fmt.Errorf("room %q: invalid scaling", r.Name)
			}
		}
		scalings.Scaling = append(scalings.Scaling, scaling)
	}
	r.Scalings = scalings
	r.Sheet = dict["sheet"].(string)
	r.WalkBoxes = make([]WalkBox, 0)
	for i, boxDict := range optionalDicts(dict["walkboxes"]) {
		box := WalkBox{}
		box.Name = optionalString(boxDict["name"])
		box.Polygon, err = parsePolygon(boxDict["polygon"].(string))
		if err != nil {
			return nil, fmt.Errorf("room %q, walkbox %q [%d]: invalid polygon", r.Name, box.Name, i)
		}
		r.WalkBoxes = append(r.WalkBoxes, box)
	}
	return r, nil
}

func optionalBool(x interface{}) bool {
	return optionalInt(x) != 0
}

func optionalInt(x interface{}) int {
	if x == nil {
		return 0
	}
	return x.(int)
}

func optionalFloat(x interface{}) float64 {
	if x == nil {
		return 0
	}
	if _, ok := x.(int); ok {
		return float64(x.(int))
	}
	return x.(float64)
}

func optionalString(x interface{}) string {
	if x == nil {
		return ""
	}
	return x.(string)
}

func optionalDicts(x interface{}) []map[string]interface{} {
	xs := optionalSlice(x)
	slice := make([]map[string]interface{}, len(xs))
	for i, elem := range xs {
		slice[i] = elem.(map[string]interface{})
	}
	return slice
}

func optionalStrings(x interface{}) []string {
	xs := optionalSlice(x)
	slice := make([]string, len(xs))
	for i, elem := range xs {
		if elem == nil {
			slice[i] = ""
			continue
		}
		slice[i] = elem.(string)
	}
	return slice
}

func optionalSlice(x interface{}) []interface{} {
	if x == nil {
		return nil
	}
	return x.([]interface{})
}
