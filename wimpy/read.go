// Copyright 2020 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package wimpy reads and writes wimpy files.
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
	return dictToRoom(dict)
}

func dictToRoom(dict map[string]interface{}) (r *Room, err error) {
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
	r.Sheet = dict["sheet"].(string)
	switch bg := dict["background"].(type) {
	case string:
		r.Background = []string{bg}
	case []interface{}:
		r.Background = optionalStrings(bg)
	}
	r.Fullscreen = optionalInt(dict["fullscreen"])
	r.Height = optionalInt(dict["height"])
	layers := optionalDicts(dict["layers"])
	r.Layers = make([]Layer, len(layers))
	for i, layerDict := range layers {
		layer := Layer{}
		switch name := layerDict["name"].(type) {
		case string:
			layer.Name = []string{name}
		case []interface{}:
			layer.Name = optionalStrings(name)
		}
		switch parallax := layerDict["parallax"].(type) {
		case string:
			p, err := parsePointF(parallax)
			if err != nil {
				return nil, fmt.Errorf("room %q, layer %q [%d]: invalid parallax string", r.Name, layer.Name, i)
			}
			layer.Parallax = p
		case float64:
			layer.Parallax = PointF{X: parallax, Y: 1}
		}
		layer.ZSort = layerDict["zsort"].(int)
		r.Layers[i] = layer
	}
	objects := optionalDicts(dict["objects"])
	r.Objects = make([]Object, len(objects))
	for i, objDict := range objects {
		obj := Object{}
		obj.Name = objDict["name"].(string)
		obj.Parent = optionalString(objDict["parent"])
		obj.Animations = readAnimations(objDict["animations"])
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
		r.Objects[i] = obj
	}
	r.RoomSize, err = parsePoint(dict["roomsize"].(string))
	if err != nil {
		return nil, fmt.Errorf("room %q: invalid roomsize", r.Name)
	}
	simpleScalings := Scalings{}
	for _, sc := range optionalSlice(dict["scaling"]) {
		switch sc := sc.(type) {
		case map[string]interface{}:
			scalings := Scalings{}
			for i, sx := range optionalStrings(sc["scaling"]) {
				s, err := parseScaling(sx)
				if err != nil {
					return nil, fmt.Errorf("room %q, scaling [%d]: invalid scaling", r.Name, i)
				}
				scalings.Scaling = append(scalings.Scaling, s)
			}
			scalings.Trigger = optionalString(sc["trigger"])
			r.Scaling = append(r.Scaling, scalings)
		case string:
			s, err := parseScaling(sc)
			if err != nil {
				return nil, fmt.Errorf("room %q: invalid scaling", r.Name)
			}
			simpleScalings.Scaling = append(simpleScalings.Scaling, s)
		}
	}
	if len(simpleScalings.Scaling) > 0 {
		r.Scaling = append(r.Scaling, simpleScalings)
	}
	walkboxes := optionalDicts(dict["walkboxes"])
	r.WalkBoxes = make([]WalkBox, len(walkboxes))
	for i, boxDict := range walkboxes {
		box := WalkBox{}
		box.Name = optionalString(boxDict["name"])
		box.Polygon, err = parsePolygon(boxDict["polygon"].(string))
		if err != nil {
			return nil, fmt.Errorf("room %q, walkbox %q [%d]: invalid polygon", r.Name, box.Name, i)
		}
		r.WalkBoxes[i] = box
	}
	return r, nil
}

func readAnimations(x interface{}) []Animation {
	dicts := optionalDicts(x)
	animations := make([]Animation, len(dicts))
	for i, dict := range dicts {
		anim := Animation{}
		anim.Name = dict["name"].(string)
		anim.FPS = optionalFloat(dict["fps"])
		anim.Triggers = optionalStrings(dict["triggers"])
		anim.Frames = optionalStrings(dict["frames"])
		anim.Loop = optionalBool(dict["loop"])
		anim.Flags = optionalInt(dict["flags"])
		anim.Layers = readAnimations(dict["layers"])
		animations[i] = anim
	}
	return animations
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
