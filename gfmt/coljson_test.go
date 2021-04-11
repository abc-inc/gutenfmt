package gfmt_test

import (
	"strings"
	"testing"

	. "github.com/abc-inc/gutenfmt/gfmt"
	. "github.com/stretchr/testify/require"
)

func TestColJSON_Write(t *testing.T) {
	tests := []struct {
		name string
		arg  interface{}
		want string
	}{
		{"nil", nil, ""},
		{"bool", false, "false"},
		{"int", -42, "-42"},
		{"string", "∮∯∰", "∮∯∰"},
		{"empty_array", [0]string{}, `^\[\]$`},
		{"int_slice", []int{1, 2, 3}, `\[\n.*1.*,\n.*2.*,\n.*3.*\n\]`},
		{"struct", NewUser("John", "Doe"), `"email".+: .+"john.doe@local"`},
		{"mixed_array", []interface{}{[0]string{}, true, -42, "a", NewUser("f", "l")}, `.*true.+,\n.*-42.+,(\n.+)+"f l"`},
		{"map", map[string]interface{}{"a a": 1, ":": ":"}, `":".+:.+":".+`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &strings.Builder{}
			_, err := NewColJSON(b).Write(tt.arg)
			NoError(t, err)
			Regexp(t, tt.want, b.String())
		})
	}
}

func TestColJSON_WriteJSONTypes(t *testing.T) {
	b := &strings.Builder{}
	_, err := NewColJSON(b).Write(jsonTypes)
	NoError(t, err)
	Contains(t, b.String(), "\"Ptr\"\x1b[0m:")
}

func TestColJSON_WriteStruct(t *testing.T) {
	b := &strings.Builder{}
	_, err := NewColJSON(b).Write(NewUser("John", "Doe"))
	NoError(t, err)
	Regexp(t, `\x1b\[\d+m.*\x1b\[\d+m"John Doe"`, b.String())
}
