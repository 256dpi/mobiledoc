package mobiledoc

import "fmt"

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

// Parse will parse the specified raw mobiledoc.
func Parse(doc Map) (Document, error) {
	// prepare document
	d := Document{}

	// get version
	version, ok := doc["version"].(string)
	if !ok {
		return d, fmt.Errorf("invalid version")
	}

	// set version
	d.Version = version

	// check markups
	if _markups, ok := doc["markups"]; ok {
		// coerce value
		markups, ok := _markups.(List)
		if !ok {
			return d, fmt.Errorf("invalid markups definition")
		}

		// parse markups
		for _, _markup := range markups {
			// coerce value
			markup, ok := _markup.(List)
			if !ok {
				return d, fmt.Errorf("invalid markups definition")
			}

			// parse markup
			m, err := parseMarkup(markup)
			if err != nil {
				return d, err
			}

			// add markup
			d.Markups = append(d.Markups, m)
		}
	}

	// check atoms
	if value, ok := doc["atoms"]; ok {
		// coerce value
		atoms, ok := value.(List)
		if !ok {
			return d, fmt.Errorf("invalid atoms definition")
		}

		// parse atoms
		for _, _atom := range atoms {
			// coerce value
			atom, ok := _atom.(List)
			if !ok {
				return d, fmt.Errorf("invalid atoms definition")
			}

			// parse atom
			a, err := parseAtom(atom)
			if err != nil {
				return d, err
			}

			// add atom
			d.Atoms = append(d.Atoms, a)
		}
	}

	// check cards
	if value, ok := doc["cards"]; ok {
		// coerce value
		cards, ok := value.(List)
		if !ok {
			return d, fmt.Errorf("invalid cards definition")
		}

		// parse cards
		for _, _card := range cards {
			// coerce value
			card, ok := _card.(List)
			if !ok {
				return d, fmt.Errorf("invalid cards definition")
			}

			// parse card
			c, err := parseCard(card)
			if err != nil {
				return d, err
			}

			// add card
			d.Cards = append(d.Cards, c)
		}
	}

	// check sections
	if value, ok := doc["sections"]; ok {
		// coerce value
		sections, ok := value.(List)
		if !ok {
			return d, fmt.Errorf("invalid sections definition")
		}

		// parse sections
		for _, _section := range sections {
			// coerce value
			section, ok := _section.(List)
			if !ok {
				return d, fmt.Errorf("invalid sections definition")
			}

			// parse section
			s, err := parseSection(section, d.Markups, d.Atoms, d.Cards)
			if err != nil {
				return d, err
			}

			// add section
			d.Sections = append(d.Sections, s)
		}
	}

	return d, nil
}

func parseMarkup(markup List) (Markup, error) {
	// prepare markup
	m := Markup{}

	// validate length
	if len(markup) == 0 || len(markup) > 2 {
		return m, fmt.Errorf("invalid markup definition")
	}

	// check tag
	tag, ok := markup[0].(string)
	if !ok || len(tag) == 0 {
		return m, fmt.Errorf("invalid markup tag")
	}

	// set tag
	m.Tag = tag

	// return if attributes are missing
	if len(markup) == 1 {
		return m, nil
	}

	// get attributes
	attributes, ok := markup[1].(List)
	if !ok {
		return m, fmt.Errorf("invalid markup attributes")
	}

	// check if attributes are even
	if len(attributes)%2 != 0 {
		return m, fmt.Errorf("invalid markup attributes")
	}

	// initialize attributes
	m.Attributes = Map{}

	// check attributes
	for i := 0; i < len(attributes); i += 2 {
		// get name
		name, ok := attributes[i].(string)
		if !ok {
			return m, fmt.Errorf("invalid markup attributes key")
		}

		// set attribute
		m.Attributes[name] = attributes[i+1]
	}

	return m, nil
}

func parseAtom(atom List) (Atom, error) {
	// prepare atom
	a := Atom{}

	// validate length
	if len(atom) != 3 {
		return a, fmt.Errorf("invalid atom definition")
	}

	// check name
	name, ok := atom[0].(string)
	if !ok {
		return a, fmt.Errorf("invalid atom name")
	}

	// set name
	a.Name = name

	// get text
	text, ok := atom[1].(string)
	if !ok {
		return a, fmt.Errorf("invalid atom text")
	}

	// set text
	a.Text = text

	// get payload
	payload, ok := atom[2].(Map)
	if !ok {
		return a, fmt.Errorf("invalid atom payload")
	}

	// set payload
	a.Payload = payload

	return a, nil
}

func parseCard(card List) (Card, error) {
	// prepare card
	c := Card{}

	// validate length
	if len(card) != 2 {
		return c, fmt.Errorf("invalid card definition")
	}

	// check name
	name, ok := card[0].(string)
	if !ok {
		return c, fmt.Errorf("invalid card name")
	}

	// set name
	c.Name = name

	// get payload
	payload, ok := card[1].(Map)
	if !ok {
		return c, fmt.Errorf("invalid card payload")
	}

	// set payload
	c.Payload = payload

	return c, nil
}

func parseSection(section List, markups []Markup, atoms []Atom, cards []Card) (Section, error) {
	// prepare section
	s := Section{}

	// validate length
	if len(section) == 0 {
		return s, fmt.Errorf("invalid section definition")
	}

	// get section type
	_typ, ok := toInt(section[0])
	if !ok {
		return s, fmt.Errorf("invalid section type")
	}

	// run validators based on type
	switch SectionType(_typ) {
	case MarkupSection:
		return parseMarkupSection(section, markups, atoms)
	case ImageSection:
		return parseImageSection(section)
	case ListSection:
		return parseListSection(section, markups, atoms)
	case CardSection:
		return parseCardSection(section, cards)
	default:
		return s, fmt.Errorf("invalid section type")
	}
}

func parseMarkupSection(section List, markups []Markup, atoms []Atom) (Section, error) {
	// prepare section
	s := Section{Type: MarkupSection}

	// validate length
	if len(section) != 3 {
		return s, fmt.Errorf("invalid markup section definition")
	}

	// get tag
	tag, ok := section[1].(string)
	if !ok {
		return s, fmt.Errorf("invalid markup section tag")
	}

	// set tag
	s.Tag = tag

	// get items
	items, ok := section[2].(List)
	if !ok {
		return s, fmt.Errorf("invalid markup section items")
	}

	// prepare open markup counter
	openMarkups := 0

	// prepare vars
	var err error
	var m Marker

	// validate markers
	for _, _marker := range items {
		// coerce value
		marker, ok := _marker.(List)
		if !ok {
			return s, fmt.Errorf("invalid markup section marker definition")
		}

		// validate marker
		m, openMarkups, err = parseMarker(marker, markups, atoms, openMarkups)
		if err != nil {
			return s, err
		}

		// add marker
		s.Markers = append(s.Markers, m)
	}

	return s, nil
}

func parseImageSection(image List) (Section, error) {
	// prepare section
	s := Section{Type: ImageSection}

	// validate length
	if len(image) != 2 {
		return s, fmt.Errorf("invalid image section definition")
	}

	// get source
	source, ok := image[1].(string)
	if !ok {
		return s, fmt.Errorf("invalid image section source")
	}

	// set source
	s.Source = source

	return s, nil
}

func parseListSection(list List, markups []Markup, atoms []Atom) (Section, error) {
	// prepare section
	s := Section{Type: ListSection}

	// validate length
	if len(list) != 3 {
		return s, fmt.Errorf("invalid list section definition")
	}

	// get tag
	tag, ok := list[1].(string)
	if !ok {
		return s, fmt.Errorf("invalid list section tag")
	}

	// set tag
	s.Tag = tag

	// get items
	items, ok := list[2].(List)
	if !ok {
		return s, fmt.Errorf("invalid list section items")
	}

	// validate items
	for _, _item := range items {
		// coerce value
		item, ok := _item.(List)
		if !ok {
			return s, fmt.Errorf("invalid list section item")
		}

		// prepare open markup counter
		openMarkups := 0

		// prepare list
		var list []Marker

		// prepare vars
		var err error
		var m Marker

		// validate markers
		for _, _marker := range item {
			// coerce value
			marker, ok := _marker.(List)
			if !ok {
				return s, fmt.Errorf("invalid list section item marker")
			}

			// validate marker
			m, openMarkups, err = parseMarker(marker, markups, atoms, openMarkups)
			if err != nil {
				return s, err
			}

			// add marker
			list = append(list, m)
		}

		// add list
		s.Items = append(s.Items, list)
	}

	return s, nil
}

func parseCardSection(card List, cards []Card) (Section, error) {
	// prepare card
	s := Section{Type: CardSection}

	// validate length
	if len(card) != 2 {
		return s, fmt.Errorf("invalid card section definition")
	}

	// get index
	index, ok := toInt(card[1])
	if !ok {
		return s, fmt.Errorf("invalid card section index")
	}

	// check index
	if index >= len(cards) {
		return s, fmt.Errorf("invalid card section index")
	}

	// set card
	s.Card = &cards[index]

	return s, nil
}

func parseMarker(marker List, markups []Markup, atoms []Atom, openMarkups int) (Marker, int, error) {
	// prepare marker
	m := Marker{}

	// validate length
	if len(marker) != 4 {
		return m, 0, fmt.Errorf("invalid marker definition")
	}

	// get marker type
	_typ, ok := toInt(marker[0])
	if !ok {
		return m, 0, fmt.Errorf("invalid marker type")
	}

	// check marker type
	typ := MarkerType(_typ)
	if typ != TextMarker && typ != AtomMarker {
		return m, 0, fmt.Errorf("invalid marker type")
	}

	// set type
	m.Type = typ

	// get opened markups
	openedMarkups, ok := marker[1].(List)
	if !ok {
		return m, 0, fmt.Errorf("invalid marker opened markups")
	}

	// validate opened markups
	for _, markup := range openedMarkups {
		// coerce value
		index, ok := toInt(markup)
		if !ok {
			return m, 0, fmt.Errorf("invalid marker markup index")
		}

		// check index
		if index >= len(markups) {
			return m, 0, fmt.Errorf("invalid marker markup index")
		}

		// add markup
		m.OpenMarkups = append(m.OpenMarkups, &markups[index])

		// increment counter
		openMarkups++
	}

	// get closed markups
	closedMarkups, ok := toInt(marker[2])
	if !ok {
		return m, 0, fmt.Errorf("invalid marker closed markup")
	}

	// decrement counter
	openMarkups -= closedMarkups
	if openMarkups < 0 {
		return m, 0, fmt.Errorf("invalid marker open markups count")
	}

	// set closed markups
	m.ClosedMarkups = closedMarkups

	// validate text marker
	if typ == TextMarker {
		// get text
		text, ok := marker[3].(string)
		if !ok {
			return m, 0, fmt.Errorf("invalid marker text")
		}

		// set text
		m.Text = text
	}

	// validate atom marker
	if typ == AtomMarker {
		// get index
		index, ok := toInt(marker[3])
		if !ok || index >= len(atoms) {
			return m, 0, fmt.Errorf("invalid marker atom index")
		}

		// set atom
		m.Atom = &atoms[index]
	}

	return m, openMarkups, nil
}
