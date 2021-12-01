package mobiledoc

import "fmt"

const scratchSize = 128

var formatValidator = NewFormatValidator()

// Compile will compile provided document into its raw structure.
func Compile(doc Document) (Map, error) {
	// validate document
	err := formatValidator.Validate(doc)
	if err != nil {
		return nil, err
	}

	// prepare compiler
	compiler := compiler{doc: doc}

	// compile document
	result := compiler.compile()
	if len(compiler.errors) > 0 {
		return nil, compiler.errors[0]
	}

	return result, nil
}

type compiler struct {
	doc     Document
	scratch List
	errors  []error
}

func (c *compiler) compile() Map {
	return Map{
		"version":  c.doc.Version,
		"markups":  c.compileMarkups(c.doc.Markups),
		"atoms":    c.compileAtoms(c.doc.Atoms),
		"cards":    c.compileCards(c.doc.Cards),
		"sections": c.compileSections(c.doc.Sections),
	}
}

func (c *compiler) compileMarkups(markups []Markup) List {
	list := c.allocate(len(markups))
	for i, markup := range markups {
		if len(markup.Attributes) > 0 {
			l := c.allocate(2)
			l[0] = markup.Tag
			l[1] = c.mapToList(markup.Attributes)
			list[i] = l
		} else {
			l := c.allocate(1)
			l[0] = markup.Tag
			list[i] = l
		}
	}
	return list
}

func (c *compiler) compileAtoms(atoms []Atom) List {
	list := c.allocate(len(atoms))
	for i, atom := range atoms {
		l := c.allocate(3)
		l[0] = atom.Name
		l[1] = atom.Text
		l[2] = atom.Payload
		list[i] = l
	}
	return list
}

func (c *compiler) compileCards(cards []Card) List {
	list := c.allocate(len(cards))
	for i, card := range cards {
		l := c.allocate(2)
		l[0] = card.Name
		l[1] = card.Payload
		list[i] = l
	}
	return list
}

func (c *compiler) compileSections(sections []Section) List {
	list := c.allocate(len(sections))
	for i, section := range sections {
		switch section.Type {
		case MarkupSection:
			l := c.allocate(3)
			l[0] = section.Type
			l[1] = section.Tag
			l[2] = c.compileMarkers(section.Markers)
			list[i] = l
		case ImageSection:
			l := c.allocate(2)
			l[0] = section.Type
			l[1] = section.Source
			list[i] = l
		case ListSection:
			l := c.allocate(3)
			l[0] = section.Type
			l[1] = section.Tag
			l[2] = c.compileItems(section.Items)
			list[i] = l
		case CardSection:
			l := c.allocate(2)
			l[0] = section.Type
			l[1] = c.cardIndex(section.Card)
			list[i] = l
		}
	}
	return list
}

func (c *compiler) compileMarkers(markers []Marker) List {
	list := c.allocate(len(markers))
	for i, marker := range markers {
		switch marker.Type {
		case TextMarker:
			l := c.allocate(4)
			l[0] = marker.Type
			l[1] = c.markupIndexes(marker.OpenMarkups)
			l[2] = marker.ClosedMarkups
			l[3] = marker.Text
			list[i] = l
		case AtomMarker:
			l := c.allocate(4)
			l[0] = marker.Type
			l[1] = c.markupIndexes(marker.OpenMarkups)
			l[2] = marker.ClosedMarkups
			l[3] = c.atomIndex(marker.Atom)
			list[i] = l
		}
	}
	return list
}

func (c *compiler) compileItems(items [][]Marker) List {
	list := c.allocate(len(items))
	for i, item := range items {
		list[i] = c.compileMarkers(item)
	}
	return list
}

func (c *compiler) markupIndexes(markups []*Markup) List {
	list := c.allocate(len(markups))
	for i, markup := range markups {
		list[i] = c.markupIndex(markup)
	}
	return list
}

func (c *compiler) markupIndex(markup *Markup) int {
	for i := range c.doc.Markups {
		if &c.doc.Markups[i] == markup {
			return i
		}
	}
	c.errors = append(c.errors, fmt.Errorf("missing markup index"))
	return -1
}

func (c *compiler) cardIndex(card *Card) int {
	for i := range c.doc.Cards {
		if &c.doc.Cards[i] == card {
			return i
		}
	}
	c.errors = append(c.errors, fmt.Errorf("missing card index"))
	return -1
}

func (c *compiler) atomIndex(atom *Atom) int {
	for i := range c.doc.Atoms {
		if &c.doc.Atoms[i] == atom {
			return i
		}
	}
	c.errors = append(c.errors, fmt.Errorf("missing atom index"))
	return -1
}

func (c *compiler) mapToList(m Map) List {
	list := c.allocate(len(m) * 2)
	i := 0
	for key, value := range m {
		list[i] = key
		i++
		list[i] = value
		i++
	}
	return list
}

func (c *compiler) allocate(length int) List {
	// check big list
	if length > scratchSize {
		return make(List, length)
	}

	// otherwise, ensure scratch
	if len(c.scratch) < length {
		c.scratch = make(List, 1024)
	}

	// get list
	list := c.scratch[0:length]
	c.scratch = c.scratch[length:]

	return list
}
