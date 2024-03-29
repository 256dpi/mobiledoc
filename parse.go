package mobiledoc

import "fmt"

// Parse will parse the specified raw structure into a document.
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
	if value, ok := doc["markups"]; ok && value != nil {
		// coerce value
		markups, ok := toList(value)
		if !ok {
			return d, fmt.Errorf("invalid markups definition")
		}

		// allocate markups
		d.Markups = make([]Markup, 0, len(markups))

		// parse markups
		for _, value := range markups {
			// coerce value
			markup, ok := toList(value)
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
	if value, ok := doc["atoms"]; ok && value != nil {
		// coerce value
		atoms, ok := toList(value)
		if !ok {
			return d, fmt.Errorf("invalid atoms definition")
		}

		// allocate atoms
		d.Atoms = make([]Atom, 0, len(atoms))

		// parse atoms
		for _, value := range atoms {
			// coerce value
			atom, ok := toList(value)
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
	if value, ok := doc["cards"]; ok && value != nil {
		// coerce value
		cards, ok := toList(value)
		if !ok {
			return d, fmt.Errorf("invalid cards definition")
		}

		// allocate cards
		d.Cards = make([]Card, 0, len(cards))

		// parse cards
		for _, value := range cards {
			// coerce value
			card, ok := toList(value)
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
	if value, ok := doc["sections"]; ok && value != nil {
		// coerce value
		sections, ok := toList(value)
		if !ok {
			return d, fmt.Errorf("invalid sections definition")
		}

		// allocate sections
		d.Sections = make([]Section, 0, len(sections))

		// parse sections
		for _, value := range sections {
			// coerce value
			section, ok := toList(value)
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

	// check length
	if len(markup) == 0 || len(markup) > 2 {
		return m, fmt.Errorf("invalid markup definition")
	}

	// get tag
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
	attributes, ok := toList(markup[1])
	if !ok {
		return m, fmt.Errorf("invalid markup attributes")
	} else if len(attributes)%2 != 0 {
		return m, fmt.Errorf("invalid markup attributes")
	}

	// allocate attributes
	m.Attributes = make(Map, len(attributes)/2)

	// parse attributes
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

	// check length
	if len(atom) != 3 {
		return a, fmt.Errorf("invalid atom definition")
	}

	// get name
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
	payload, ok := toMap(atom[2])
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

	// check length
	if len(card) != 2 {
		return c, fmt.Errorf("invalid card definition")
	}

	// get name
	name, ok := card[0].(string)
	if !ok {
		return c, fmt.Errorf("invalid card name")
	}

	// set name
	c.Name = name

	// get payload
	payload, ok := toMap(card[1])
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

	// check length
	if len(section) == 0 {
		return s, fmt.Errorf("invalid section definition")
	}

	// get section type
	typ, ok := toInt(section[0])
	if !ok {
		return s, fmt.Errorf("invalid section type")
	}

	// parse section
	switch SectionType(typ) {
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

	// check length
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
	items, ok := toList(section[2])
	if !ok {
		return s, fmt.Errorf("invalid markup section items")
	}

	// prepare open markup counter
	openMarkups := 0

	// prepare vars
	var err error
	var m Marker

	// allocate markers
	s.Markers = make([]Marker, 0, len(items))

	// parse markers
	for _, value := range items {
		// coerce value
		marker, ok := toList(value)
		if !ok {
			return s, fmt.Errorf("invalid markup section marker definition")
		}

		// parse marker
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

	// check length
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

	// check length
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
	items, ok := toList(list[2])
	if !ok {
		return s, fmt.Errorf("invalid list section items")
	}

	// allocate items
	s.Items = make([][]Marker, 0, len(items))

	// parse items
	for _, value := range items {
		// coerce value
		item, ok := toList(value)
		if !ok {
			return s, fmt.Errorf("invalid list section item")
		}

		// prepare open markup counter
		openMarkups := 0

		// prepare vars
		var err error
		var m Marker

		// allocate markers
		list := make([]Marker, 0, len(item))

		// parse markers
		for _, value := range item {
			// coerce value
			marker, ok := toList(value)
			if !ok {
				return s, fmt.Errorf("invalid list section item marker")
			}

			// parse marker
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

	// check length
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

	// check length
	if len(marker) != 4 {
		return m, 0, fmt.Errorf("invalid marker definition")
	}

	// get marker type
	typ, ok := toInt(marker[0])
	if !ok {
		return m, 0, fmt.Errorf("invalid marker type")
	}

	// check marker type
	markerType := MarkerType(typ)
	if markerType != TextMarker && markerType != AtomMarker {
		return m, 0, fmt.Errorf("invalid marker type")
	}

	// set type
	m.Type = markerType

	// get opened markups
	openedMarkups, ok := toList(marker[1])
	if !ok {
		return m, 0, fmt.Errorf("invalid marker opened markups")
	}

	// allocate open markups
	if len(openedMarkups) > 0 {
		m.OpenMarkups = make([]*Markup, 0, len(openedMarkups))
	}

	// parse opened markups
	for _, value := range openedMarkups {
		// coerce value
		index, ok := toInt(value)
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

	// parse text marker
	if markerType == TextMarker {
		// get text
		text, ok := marker[3].(string)
		if !ok {
			return m, 0, fmt.Errorf("invalid marker text")
		}

		// set text
		m.Text = text
	}

	// parse atom marker
	if markerType == AtomMarker {
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
