package bnut_test

import (
	"io"
	"strings"
	"testing"

	"github.com/fzipp/gg/crypt/bnut"
)

func TestEncodingWriter(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"", ""},
		{"secret", "\x63\xca\x67\x6f\x23\x4e"},
		{"This is a test.", "\x9f\xdd\x40\x75\x0e\x90\x38\xea\x25\x7e\xc7\xc9\xf2\xdb\xa9"},
	}
	for _, tt := range tests {
		var sb strings.Builder
		w := bnut.EncodingWriter(&sb, int64(len(tt.input)))
		_, err := io.Copy(w, strings.NewReader(tt.input))
		if err != nil {
			t.Errorf("writing %q with bnut encryption returned error: %s", tt.input, err)
			continue
		}
		output := sb.String()
		if output != tt.want {
			t.Errorf("writing %q with bnut encryption resulted in '%x', want: '%x'", tt.input, output, tt.want)
		}
	}
}
