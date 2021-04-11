package gfmt_test

import (
	. "github.com/abc-inc/gutenfmt/gfmt"
	. "github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestNewAutoJSON(t *testing.T) {
	f := NewAutoJSON(&strings.Builder{})
	IsType(t, &JSON{}, f)
}
