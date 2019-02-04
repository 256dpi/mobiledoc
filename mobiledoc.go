package mobiledoc

import "github.com/asaskevich/govalidator"

// A is a short-hand for an array of interfaces.
type A = []interface{}

// M is a short-hand for map of interfaces.
type M = map[string]interface{}

// Markups defines the expected markups with the tag as the key and a map of
// attributes and validator functions.
var Markups = map[string]map[string]func(string) bool{
	"a":      {"href": govalidator.IsURL},
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

// Atoms defines the expected atoms with the name as the key and a validator
// function.
var Atoms = map[string]func(string, M) bool{}

// Cards defines the expected cards with the name as the key and a validator
// function.
var Cards = map[string]func(M) bool{}

// The available section identifiers.
const (
	MarkupSection = 1
	ImageSection  = 2
	ListSection   = 3
	CardSection   = 10
)

// MarkupSections defines the available markup sections.
var MarkupSections = []string{"aside", "blockquote", "h1", "h2", "h3", "h4", "h5", "h6", "p"}

// ListSections defines the available list sections.
var ListSections = []string{"ul", "ol"}

// The available marker identifiers.
const (
	TextMarker = 0
	AtomMarker = 1
)
