package mobiledoc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertText(t *testing.T) {
	doc := ConvertText("")
	assert.NoError(t, formatValidator.Validate(doc))
	assert.Equal(t, Document{
		Version: Version,
		Markups: []Markup{},
		Atoms:   []Atom{},
		Cards:   []Card{},
		Sections: []Section{
			{
				Type: MarkupSection,
				Tag:  "p",
				Markers: []Marker{
					{
						Type: TextMarker,
						Text: "",
					},
				},
			},
		},
	}, doc)

	doc = ConvertText("Hello")
	assert.NoError(t, formatValidator.Validate(doc))
	assert.Equal(t, Document{
		Version: Version,
		Markups: []Markup{},
		Atoms:   []Atom{},
		Cards:   []Card{},
		Sections: []Section{
			{
				Type: MarkupSection,
				Tag:  "p",
				Markers: []Marker{
					{
						Type: TextMarker,
						Text: "Hello",
					},
				},
			},
		},
	}, doc)

	doc = ConvertText("Hello\nWorld!\n\nAwesome!\n")
	assert.NoError(t, formatValidator.Validate(doc))
	assert.Equal(t, Document{
		Version: Version,
		Markups: []Markup{},
		Atoms:   []Atom{},
		Cards:   []Card{},
		Sections: []Section{
			{
				Type: MarkupSection,
				Tag:  "p",
				Markers: []Marker{
					{
						Type: TextMarker,
						Text: "Hello",
					},
				},
			},
			{
				Type: MarkupSection,
				Tag:  "p",
				Markers: []Marker{
					{
						Type: TextMarker,
						Text: "World!",
					},
				},
			},
			{
				Type: MarkupSection,
				Tag:  "p",
				Markers: []Marker{
					{
						Type: TextMarker,
						Text: "Awesome!",
					},
				},
			},
		},
	}, doc)
}
