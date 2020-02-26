package ggdict_test

import (
	"github.com/fzipp/gg/ggdict"
	"reflect"
	"testing"
)

func TestRoundtrip(t *testing.T) {
	dict := map[string]interface{}{
		"name":    "Test",
		"count":   4,
		"numbers": []interface{}{0.5, 3, 2.6, 1.4}, // TODO: allow []float64?
		"subobject": map[string]interface{}{
			"title": "Test 2",
			"id":    0,
		},
	}
	data := ggdict.Marshal(dict)
	newDict, err := ggdict.Unmarshal(data)
	if err != nil {
		t.Errorf("Unmarshal returned an error: %s", err)
		return
	}
	if !reflect.DeepEqual(dict, newDict) {
		t.Errorf("Marshal/unmarshal roundtrip resulted in %#v, want: %#v", newDict, dict)
	}
}
