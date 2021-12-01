package mobiledoc

// List is a general purpose list.
type List = []interface{}

// Map is a general purpose map.
type Map = map[string]interface{}

// SectionType defines a section type.
type SectionType int

// The available section identifiers.
const (
	MarkupSection SectionType = 1
	ImageSection  SectionType = 2
	ListSection   SectionType = 3
	CardSection   SectionType = 10
)

// MarkerType defines a marker type.
type MarkerType int

// The available marker identifiers.
const (
	TextMarker MarkerType = 0
	AtomMarker MarkerType = 1
)

// Document is a mobiledoc.
type Document struct {
	Version  string
	Markups  []Markup
	Atoms    []Atom
	Cards    []Card
	Sections []Section
}

// Markup is a single markup.
type Markup struct {
	Tag        string
	Attributes Map
}

// Atom is a single atom.
type Atom struct {
	Name    string
	Text    string
	Payload Map
}

// Card is a single card.
type Card struct {
	Name    string
	Payload Map
}

// Section is a single section.
type Section struct {
	Type    SectionType
	Tag     string
	Markers []Marker
	Source  string
	Items   [][]Marker
	Card    *Card
}

// Marker is a single marker.
type Marker struct {
	Type          MarkerType
	OpenMarkups   []*Markup
	ClosedMarkups int
	Text          string
	Atom          *Atom
}
