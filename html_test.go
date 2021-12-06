package mobiledoc

import (
	"bufio"
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHTMLRenderer(t *testing.T) {
	r := NewHTMLRenderer()
	r.Atoms["atom1"] = func(w *bufio.Writer, text string, payload Map) error {
		_, err := w.WriteString(fmt.Sprintf("<span class=\"atom1\">%s</span>", text))
		return err
	}
	r.Atoms["atom2"] = func(w *bufio.Writer, text string, payload Map) error {
		_, err := w.WriteString(fmt.Sprintf("<span class=\"atom2\">%s</span>", text))
		return err
	}
	r.Cards["card1"] = func(w *bufio.Writer, payload Map) error {
		_, err := w.WriteString("<div>card1</div>")
		return err
	}
	r.Cards["card2"] = func(w *bufio.Writer, payload Map) error {
		_, err := w.WriteString("<div>card2</div>")
		return err
	}

	out := `<div>card1</div><p>foo<b>foo</b><i>foofoo</i><i><a href="https://example.com">foo</a>foo</i></p><p><span class="atom1">foo</span><b><span class="atom2">foo</span><span class="atom1">foo</span></b></p><img src="https://example.com/foo.png"><ul><li>foo<b>foo</b></li><li><b>foo&lt;foo&gt;</b></li></ul><ol><li>bar<i>bar</i></li><li><i>bar&lt;bar&gt;</i></li></ol><div>card2</div>`

	buf := &bytes.Buffer{}
	err := r.Render(buf, sampleDoc())
	assert.NoError(t, err)
	assert.Equal(t, out, buf.String())
}
