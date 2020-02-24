package bnut_test

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/fzipp/gg/crypt/bnut"
)

func TestDecodingReader(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"", ""},
		{"\x63\xca\x67\x6f\x23\x4e", "secret"},
		{"\x9f\xdd\x40\x75\x0e\x90\x38\xea\x25\x7e\xc7\xc9\xf2\xdb\xa9", "This is a test."},
	}
	for _, tt := range tests {
		r := bnut.DecodingReader(strings.NewReader(tt.input), int64(len(tt.input)))
		buf, err := ioutil.ReadAll(r)
		output := string(buf)
		if err != nil {
			t.Errorf("reading %q with from encryption returned error: %s", tt.input, err)
			continue
		}
		if output != tt.want {
			t.Errorf("reading %q with from encryption resulted in %q, want: %q", tt.input, output, tt.want)
		}
	}
}
