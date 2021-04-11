package gfmt

import (
	"io"
	"os"

	"github.com/mattn/go-isatty"
)

func NewAutoJSON(w io.Writer) InterfaceWriter {
	if (w == os.Stdout && isatty.IsTerminal(os.Stdout.Fd())) ||
		(w == os.Stderr && isatty.IsTerminal(os.Stderr.Fd())) {
		return NewColJSON(w)
	}
	return NewJSON(w)
}