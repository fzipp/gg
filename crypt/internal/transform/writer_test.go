package transform_test

import (
	"io"
	"strings"
	"testing"

	"github.com/fzipp/gg/crypt/internal/transform"
)

func TestWriter(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"", ""},
		{"abcdefg", "ABCDEFG"},
		{"This is a test.", "THIS IS A TEST."},
	}
	for _, tt := range tests {
		var sb strings.Builder
		w := transform.NewWriter(&sb, upperCaseTransformer{})
		_, err := io.Copy(w, strings.NewReader(tt.input))
		if err != nil {
			t.Errorf("writing %q with transformer returned error: %s", tt.input, err)
			continue
		}
		output := sb.String()
		if output != tt.want {
			t.Errorf("writing %q with transformer resulted in %q, want: %q", tt.input, output, tt.want)
		}
	}
}
