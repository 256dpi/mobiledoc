package mobiledoc

// A is a short-hand for an array of interfaces.
type A = []interface{}

// M is a short-hand for map of interfaces.
type M = map[string]interface{}

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

// The available section identifiers.
const (
	MarkupSection = 1
	ImageSection  = 2
	ListSection   = 3
	CardSection   = 10
)

// DefaultMarkupSections defines the default markup sections.
var DefaultMarkupSections = []string{"aside", "blockquote", "h1", "h2", "h3", "h4", "h5", "h6", "p"}

// DefaultListSections defines the default list sections.
var DefaultListSections = []string{"ul", "ol"}

// DefaultImageSection defines the default image section validator.
var DefaultImageSection = func(string) bool { return true }

// The available marker identifiers.
const (
	TextMarker = 0
	AtomMarker = 1
)
