package mobiledoc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var defaultValidator = NewDefaultValidator()

func TestConvertText(t *testing.T) {
	doc := ConvertText("")
	assert.NoError(t, defaultValidator.Validate(doc))
	assert.Equal(t, Document{
		Version: Version,
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
	assert.NoError(t, defaultValidator.Validate(doc))
	assert.Equal(t, Document{
		Version: Version,
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
	assert.NoError(t, defaultValidator.Validate(doc))
	assert.Equal(t, Document{
		Version: Version,
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
