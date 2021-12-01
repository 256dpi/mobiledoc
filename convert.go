package mobiledoc

import (
	"regexp"
	"strings"
)

var lineBreaks = regexp.MustCompile("\n+")

// ConvertText will convert a basic text with no formatting but newlines to a
// document.
func ConvertText(str string) Document {
	// trim space
	str = strings.TrimSpace(str)

	// prepare doc
	doc := Document{
		Version: Version,
	}

	// add sections
	for _, section := range lineBreaks.Split(str, -1) {
		doc.Sections = append(doc.Sections, Section{
			Type: MarkupSection,
			Tag:  "p",
			Markers: []Marker{
				{
					Type: TextMarker,
					Text: section,
				},
			},
		})
	}

	return doc
}
