package mobiledoc

// List is a general purpose list.
type List = []interface{}

// Map is a general purpose map.
type Map = map[string]interface{}

// DefaultMarkups defines the default expected markups with the tag as the key
// and a map of attributes and validator functions.
var DefaultMarkups = map[string]map[string]func(string) bool{
	"a":      {"href": func(string) bool { return true }},
	"b":      nil,
	"code":   nil,
	"em":     nil,
	"i":      nil,
	"s":      nil,
	"strong": nil,
	"sub":    nil,
	"sup":    nil,
	"u":      nil,
}

// SectionType defines a section type.
type SectionType int

// The available section identifiers.
const (
	MarkupSection SectionType = 1
	ImageSection  SectionType = 2
	ListSection   SectionType = 3
	CardSection   SectionType = 10
)

// Valid returns whether the section type is valid.
func (t SectionType) Valid() bool {
	return t == MarkupSection || t == ImageSection || t == ListSection || t == CardSection
}

// DefaultMarkupSections defines the default markup sections.
var DefaultMarkupSections = []string{"aside", "blockquote", "h1", "h2", "h3", "h4", "h5", "h6", "p"}

// DefaultListSections defines the default list sections.
var DefaultListSections = []string{"ul", "ol"}

// DefaultImageSection defines the default image section validator.
var DefaultImageSection = func(string) bool { return true }

// MarkerType defines a marker type.
type MarkerType int

// The available marker identifiers.
const (
	TextMarker MarkerType = 0
	AtomMarker MarkerType = 1
)

// Valid returns whether the marker type is valid.
func (t MarkerType) Valid() bool {
	return t == TextMarker || t == AtomMarker
}
