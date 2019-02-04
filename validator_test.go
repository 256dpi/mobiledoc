package mobiledoc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidator(t *testing.T) {
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
				{Type: TextMarker, OpenMarkups: []*Markup(nil), ClosedMarkups: 0, Text: "foo"},
				{Type: TextMarker, OpenMarkups: []*Markup{&doc.Markups[0]}, ClosedMarkups: 1, Text: "foo"},
			},
			{
				{Type: TextMarker, OpenMarkups: []*Markup{&doc.Markups[0]}, Text: "foo"},
				{Type: TextMarker, OpenMarkups: []*Markup(nil), ClosedMarkups: 1, Text: "foo"},
			},
		}},
		{Type: 10, Card: &doc.Cards[1]},
	}

	v := NewValidator()

	var atom1 interface{}
	v.Atoms["atom1"] = func(_ string, payload interface{}) bool {
		atom1 = payload
		return true
	}
	v.Atoms["atom2"] = func(string, interface{}) bool {
		return true
	}

	var card1 interface{}
	v.Cards["card1"] = func(payload interface{}) bool {
		card1 = payload
		return true
	}
	v.Cards["card2"] = func(payload interface{}) bool {
		return true
	}

	err := v.Validate(doc)
	assert.NoError(t, err)
	assert.Equal(t, Map{"bar": 42}, atom1)
	assert.Equal(t, Map{"foo": 42}, card1)
}
