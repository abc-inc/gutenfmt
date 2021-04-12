package gfmt

import (
	. "github.com/stretchr/testify/require"
	"os"
	"testing"
)

func Test_wrapCountingWriter(t *testing.T) {
	w := wrapCountingWriter(os.Stdout)
	Same(t, w, wrapCountingWriter(w))
}
