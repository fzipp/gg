// Copyright 2020 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package wimpy

import (
	"fmt"
	"image"
	"io"
	"strconv"
	"strings"

	"github.com/fzipp/gg/ggdict"
)

func Write(w io.Writer, r *Room) (n int, err error) {
	return w.Write(ggdict.Marshal(roomToDict(r)))
}

func roomToDict(r *Room) map[string]interface{} {
	dict := make(map[string]interface{})
	dict["name"] = r.Name
	dict["sheet"] = r.Sheet
	if len(r.Background) == 1 {
		dict["background"] = r.Background[0]
	} else {
		dict["background"] = stringsToSlice(r.Background)
	}
	if r.Fullscreen != 0 {
		dict["fullscreen"] = r.Fullscreen
	}
	if r.Height != 0 {
		dict["height"] = r.Height
	}
	layers := make([]interface{}, len(r.Layers))
	for i, l := range r.Layers {
		layer := make(map[string]interface{})
		if len(l.Name) == 1 {
			layer["name"] = l.Name[0]
		} else {
			layer["name"] = stringsToSlice(l.Name)
		}
		if l.Parallax.Y == 1 {
			layer["parallax"] = l.Parallax.X
		} else {
			layer["parallax"] = formatPointF(l.Parallax)
		}
		layer["zsort"] = l.ZSort
		layers[i] = layer
	}
	if len(layers) > 0 {
		dict["layers"] = layers
	}
	objects := make([]interface{}, len(r.Objects))
	for i, o := range r.Objects {
		obj := make(map[string]interface{})
		obj["name"] = o.Name
		if o.Parent != "" {
			obj["parent"] = o.Parent
		}
		if len(o.Animations) > 0 {
			obj["animations"] = animationsToSlice(o.Animations)
		}
		obj["hotspot"] = fmt.Sprintf("{%s,%s}",
			formatPoint(o.HotSpot.Min),
			formatPoint(o.HotSpot.Max))
		obj["pos"] = formatPoint(o.Pos)
		obj["usedir"] = o.UseDir.String()
		obj["usepos"] = formatPoint(o.UsePos)
		obj["zsort"] = o.ZSort
		if o.Prop {
			obj["prop"] = boolToInt(o.Prop)
		}
		if o.Spot {
			obj["spot"] = boolToInt(o.Spot)
		}
		if o.Trigger {
			obj["trigger"] = boolToInt(o.Trigger)
		}
		objects[i] = obj
	}
	dict["objects"] = objects
	dict["roomsize"] = formatPoint(r.RoomSize)
	if len(r.Scaling) == 1 && r.Scaling[0].Trigger == "" {
		dict["scaling"] = scalingsToSlice(r.Scaling[0].Scaling)
	} else if len(r.Scaling) > 0 {
		scaling := make([]interface{}, len(r.Scaling))
		for i, sc := range r.Scaling {
			s := make(map[string]interface{})
			s["scaling"] = scalingsToSlice(sc.Scaling)
			if sc.Trigger != "" {
				s["trigger"] = sc.Trigger
			}
			scaling[i] = s
		}
		dict["scaling"] = scaling
	}
	walkboxes := make([]interface{}, len(r.WalkBoxes))
	for i, wb := range r.WalkBoxes {
		box := make(map[string]interface{})
		if wb.Name != "" {
			box["name"] = wb.Name
		}
		box["polygon"] = formatPolygon(wb.Polygon)
		walkboxes[i] = box
	}
	dict["walkboxes"] = walkboxes
	return dict
}

func animationsToSlice(anims []Animation) []interface{} {
	slice := make([]interface{}, len(anims))
	for i, a := range anims {
		anim := make(map[string]interface{})
		anim["name"] = a.Name
		if a.FPS != 0 {
			anim["fps"] = a.FPS
		}
		if len(a.Triggers) > 0 {
			anim["triggers"] = stringsToSlice(a.Triggers)
		}
		if len(a.Layers) == 0 || len(a.Frames) > 0 {
			anim["frames"] = stringsToSlice(a.Frames)
		}
		if a.Loop {
			anim["loop"] = boolToInt(a.Loop)
		}
		if a.Flags != 0 {
			anim["flags"] = a.Flags
		}
		if len(a.Layers) > 0 {
			anim["layers"] = animationsToSlice(a.Layers)
		}
		slice[i] = anim
	}
	return slice
}

func stringsToSlice(a []string) []interface{} {
	slice := make([]interface{}, len(a))
	for i, s := range a {
		slice[i] = s
	}
	return slice
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func scalingsToSlice(s []Scaling) []interface{} {
	slice := make([]interface{}, len(s))
	for i, sc := range s {
		slice[i] = formatScaling(sc)
	}
	return slice
}

func formatPoint(pt image.Point) string {
	return fmt.Sprintf("{%d,%d}", pt.X, pt.Y)
}

func formatPointF(pt PointF) string {
	return fmt.Sprintf("{%g,%g}", pt.X, pt.Y)
}

func formatScaling(s Scaling) string {
	return fmt.Sprintf("%s@%d", formatScalingFactor(s.Factor), s.At)
}

func formatScalingFactor(f float64) string {
	s := strconv.FormatFloat(f, 'g', -1, 64)
	if strings.HasPrefix(s, "0.") {
		return s[1:]
	}
	if !strings.Contains(s, ".") {
		return s + ".0"
	}
	return s
}

func formatPolygon(polygon []image.Point) string {
	var sb strings.Builder
	for i, pt := range polygon {
		if i > 0 {
			sb.WriteRune(';')
		}
		sb.WriteString(formatPoint(pt))
	}
	return sb.String()
}
