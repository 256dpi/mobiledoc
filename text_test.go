package mobiledoc

import (
	"bufio"
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTextRenderer(t *testing.T) {
	r := NewTextRenderer()
	r.Atoms["atom1"] = func(w *bufio.Writer, text string, payload Map) error {
		_, err := w.WriteString(fmt.Sprintf("(atom1: %s)", text))
		return err
	}
	r.Atoms["atom2"] = func(w *bufio.Writer, text string, payload Map) error {
		_, err := w.WriteString(fmt.Sprintf("(atom2: %s)", text))
		return err
	}
	r.Cards["card1"] = func(w *bufio.Writer, payload Map) error {
		_, err := w.WriteString("(card1)")
		return err
	}
	r.Cards["card2"] = func(w *bufio.Writer, payload Map) error {
		_, err := w.WriteString("(card2)")
		return err
	}

	out := `(card1)
foo foo foo foo foo foo
(atom1: foo) (atom2: foo) (atom1: foo)
[https://example.com/foo.png]
- foo foo
- foo <foo>
1. bar bar
2. bar <bar>
(card2)`

	buf := &bytes.Buffer{}
	err := r.Render(buf, sampleDoc())
	assert.NoError(t, err)
	assert.Equal(t, out, buf.String())
}
