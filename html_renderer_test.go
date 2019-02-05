package mobiledoc

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHTMLRenderer(t *testing.T) {
	doc := Document{
		Version: Version,
		Markups: []Markup{
			{Tag: "b"},
			{Tag: "i"},
			{Tag: "a", Attributes: Map{"href": "http://example.com"}},
		},
		Atoms: []Atom{
			{Name: "atom1", Text: "foo", Payload: Map{"bar": 42}},
			{Name: "atom2", Text: "foo", Payload: Map{"bar": 24}},
		},
		Cards: []Card{
			{Name: "card1", Payload: Map{"foo": 42}},
			{Name: "card2", Payload: Map{"foo": 42}},
		},
	}
	doc.Sections = []Section{
		{Type: CardSection, Card: &doc.Cards[0]},
		{Type: MarkupSection, Tag: "p", Markers: []Marker{
			{Type: TextMarker, Text: "foo"},
			{Type: TextMarker, OpenMarkups: []*Markup{&doc.Markups[0]}, ClosedMarkups: 1, Text: "foo"},
			{Type: TextMarker, OpenMarkups: []*Markup{&doc.Markups[1]}, Text: "foo"},
			{Type: TextMarker, ClosedMarkups: 1, Text: "foo"},
			{Type: TextMarker, OpenMarkups: []*Markup{&doc.Markups[1], &doc.Markups[2]}, ClosedMarkups: 1, Text: "foo"},
			{Type: TextMarker, ClosedMarkups: 1, Text: "foo"},
		}},
		{Type: MarkupSection, Tag: "p", Markers: []Marker{
			{Type: AtomMarker, Atom: &doc.Atoms[0]},
			{Type: AtomMarker, OpenMarkups: []*Markup{&doc.Markups[0]}, Atom: &doc.Atoms[1]},
			{Type: AtomMarker, ClosedMarkups: 1, Atom: &doc.Atoms[0]},
		}},
		{Type: ImageSection, Source: "http://example.com/foo.png"},
		{Type: ListSection, Tag: "ul", Items: [][]Marker{
			{
				{Type: TextMarker, ClosedMarkups: 0, Text: "foo"},
				{Type: TextMarker, OpenMarkups: []*Markup{&doc.Markups[0]}, ClosedMarkups: 1, Text: "foo"},
			},
			{
				{Type: TextMarker, OpenMarkups: []*Markup{&doc.Markups[0]}, Text: "foo"},
				{Type: TextMarker, ClosedMarkups: 1, Text: "<foo>"},
			},
		}},
		{Type: 10, Card: &doc.Cards[1]},
	}

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

	out := `<div>card1</div><p>foo<b>foo</b><i>foofoo</i><i><a href="http://example.com">foo</a>foo</i></p><p><span class="atom1">foo</span><b><span class="atom2">foo</span><span class="atom1">foo</span></b></p><img src="http://example.com/foo.png"><ul><li>foo<b>foo</b></li><li><b>foo&lt;foo&gt;</b></li></ul><div>card2</div>`

	buf := &bytes.Buffer{}
	err := r.Render(buf, doc)
	assert.NoError(t, err)
	assert.Equal(t, out, buf.String())
}
