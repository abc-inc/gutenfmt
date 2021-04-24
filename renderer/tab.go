package renderer

import (
	"strings"
	"text/tabwriter"
)

func AsTab(r Renderer) Renderer {
	return RendererFunc(func(i interface{}) (string, error) {
		b := &strings.Builder{}
		tw := tabwriter.NewWriter(b, 1, 4, 3, ' ', 0)
		_, err := RenderTab(tw, r, i)
		return b.String(), err
	})
}

func RenderTab(tw *tabwriter.Writer, r Renderer, i interface{}) (int, error) {
	s, err := r.Render(i)
	if err != nil {
		return 0, err
	}

	if n, err := tw.Write([]byte(s)); err != nil {
		return n, err
	}
	return 0, tw.Flush()
}
