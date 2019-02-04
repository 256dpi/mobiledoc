package mobiledoc

import "fmt"

// TODO: Allow "target" attribute for URLs?

// TODO: Add parser that transforms the mobiledoc in a usable in memory layout.

// Validator validates a mobiledoc.
type Validator struct {
	// MarkupSections defines the expected markup sections.
	MarkupSections []string

	// ListSections defines the expected list sections.
	ListSections []string

	// ImageSection defines whether the image section is allowed.
	ImageSection func(string) bool

	// Markups defines the expected markups with the name as key and a map of
	// attributes and validations functions.
	Markups map[string]map[string]func(string) bool

	// Atoms defines the expected atoms with the name as the key and a validator
	// function.
	Atoms map[string]func(string, M) bool

	// Cards defines the expected cards with the name as the key and a validator
	// function.
	Cards map[string]func(M) bool
}

// NewValidator creates a validator that validates the mobiledoc standard.
func NewValidator() *Validator {
	return &Validator{
		MarkupSections: DefaultMarkupSections,
		ListSections:   DefaultListSections,
		ImageSection:   DefaultImageSection,
		Markups:        DefaultMarkups,
		Atoms:          make(map[string]func(string, M) bool),
		Cards:          make(map[string]func(M) bool),
	}
}

// Validate will walk the specified mobiledoc and check if it is valid.
func (v *Validator) Validate(doc M) error {
	// check version
	if version, ok := doc["version"]; !ok || version != "0.3.1" {
		return fmt.Errorf("invalid mobiledoc version")
	}

	// prepare num markups
	numMarkups := 0

	// check markups
	if _markups, ok := doc["markups"]; ok {
		// coerce value
		markups, ok := _markups.(A)
		if !ok {
			return fmt.Errorf("invalid markups definition")
		}

		// set num
		numMarkups = len(markups)

		// validate markups
		for _, _markup := range markups {
			// coerce value
			markup, ok := _markup.(A)
			if !ok {
				return fmt.Errorf("invalid markups definition")
			}

			// validate markup
			err := v.ValidateMarkup(markup)
			if err != nil {
				return err
			}
		}
	}

	// prepare num atoms
	numAtoms := 0

	// check atoms
	if value, ok := doc["atoms"]; ok {
		// coerce value
		atoms, ok := value.(A)
		if !ok {
			return fmt.Errorf("invalid atoms definition")
		}

		// set num
		numAtoms = len(atoms)

		// validate atom
		for _, _atom := range atoms {
			// coerce value
			atom, ok := _atom.(A)
			if !ok {
				return fmt.Errorf("invalid atoms definition")
			}

			// validate atom
			err := v.ValidateAtom(atom)
			if err != nil {
				return err
			}
		}
	}

	// prepare num cards
	numCards := 0

	// check cards
	if value, ok := doc["cards"]; ok {
		// coerce value
		cards, ok := value.(A)
		if !ok {
			return fmt.Errorf("invalid cards definition")
		}

		// set num
		numCards = len(cards)

		// validate cards
		for _, _card := range cards {
			// coerce value
			card, ok := _card.(A)
			if !ok {
				return fmt.Errorf("invalid cards definition")
			}

			// validate card
			err := v.ValidateCard(card)
			if err != nil {
				return err
			}
		}
	}

	// check sections
	if value, ok := doc["sections"]; ok {
		// coerce value
		sections, ok := value.(A)
		if !ok {
			return fmt.Errorf("invalid sections definition")
		}

		// validate sections
		for _, _section := range sections {
			// coerce value
			section, ok := _section.(A)
			if !ok {
				return fmt.Errorf("invalid sections definition")
			}

			// validate section
			err := v.ValidateSection(section, numMarkups, numAtoms, numCards)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// ValidateMarkup will validate a single markup.
func (v *Validator) ValidateMarkup(markup A) error {
	// validate length
	if len(markup) == 0 || len(markup) > 2 {
		return fmt.Errorf("invalid markup definition")
	}

	// check tag
	tag, ok := markup[0].(string)
	if !ok {
		return fmt.Errorf("invalid markup definition")
	}

	// check tag existence
	allowedAttributes, ok := v.Markups[tag]
	if !ok {
		return fmt.Errorf("invalid markup tag")
	}

	// return if attributes are missing
	if len(markup) == 1 {
		return nil
	}

	// get attributes
	attributes, ok := markup[1].(A)
	if !ok {
		return fmt.Errorf("invalid markup definition")
	}

	// check if attributes are even
	if len(attributes)%2 != 0 {
		return fmt.Errorf("invalid markup definition")
	}

	// check attributes
	for i := 0; i < len(attributes); i += 2 {
		// coerce name
		name, ok1 := attributes[i].(string)
		value, ok2 := attributes[i+1].(string)
		if !ok1 || !ok2 {
			return fmt.Errorf("invalid markup definition")
		}

		// get validator
		validator, ok := allowedAttributes[name]
		if !ok {
			return fmt.Errorf("invalid markup attribute")
		}

		// validate attribute
		if !validator(value) {
			return fmt.Errorf("invalid markup attribute")
		}
	}

	return nil
}

// ValidateAtom will validate a single atom.
func (v *Validator) ValidateAtom(atom A) error {
	// validate length
	if len(atom) == 0 || len(atom) > 3 {
		return fmt.Errorf("invalid atom definition")
	}

	// check name
	name, ok := atom[0].(string)
	if !ok {
		return fmt.Errorf("invalid atom definition")
	}

	// check atom existence
	validator, ok := v.Atoms[name]
	if !ok {
		return fmt.Errorf("invalid atom name")
	}

	// prepare text and payload
	var text string
	var payload M

	// get text
	if len(atom) > 1 {
		text, ok = atom[1].(string)
		if !ok {
			return fmt.Errorf("invalid atom definition")
		}
	}

	// get payload
	if len(atom) > 2 {
		payload, ok = atom[2].(M)
		if !ok {
			return fmt.Errorf("invalid atom definition")
		}
	}

	// validate atom
	if !validator(text, payload) {
		return fmt.Errorf("invalid atom text or payload")
	}

	return nil
}

// ValidateCard will validate a single card.
func (v *Validator) ValidateCard(card A) error {
	// validate length
	if len(card) == 0 || len(card) > 2 {
		return fmt.Errorf("invalid card definition")
	}

	// check name
	name, ok := card[0].(string)
	if !ok {
		return fmt.Errorf("invalid card definition")
	}

	// check card existence
	validator, ok := v.Cards[name]
	if !ok {
		return fmt.Errorf("invalid card name")
	}

	// prepare payload
	var payload M

	// get payload
	if len(card) > 1 {
		payload, ok = card[1].(M)
		if !ok {
			return fmt.Errorf("invalid card definition")
		}
	}

	// validate card
	if !validator(payload) {
		return fmt.Errorf("invalid card payload")
	}

	return nil
}

// ValidateSection will validate a single section.
func (v *Validator) ValidateSection(section A, numMarkups, numAtoms, numCards int) error {
	// validate length
	if len(section) == 0 {
		return fmt.Errorf("invalid section definition")
	}

	// get type
	typ, ok := toInt(section[0])
	if !ok {
		return fmt.Errorf("invalid section definition")
	}

	// run validators based on type
	switch typ {
	case MarkupSection:
		return v.ValidateMarkupSection(section, numMarkups, numAtoms)
	case ImageSection:
		return v.ValidateImageSection(section)
	case ListSection:
		return v.ValidateListSection(section, numMarkups, numAtoms)
	case CardSection:
		return v.ValidateCardSection(section, numCards)
	default:
		return fmt.Errorf("invalid section definition")
	}
}

// ValidateMarkupSection validates a single markup section.
func (v *Validator) ValidateMarkupSection(section A, numMarkups, numAtoms int) error {
	// validate length
	if len(section) != 3 {
		return fmt.Errorf("invalid markup section definition")
	}

	// get tag
	tag, ok := section[1].(string)
	if !ok {
		return fmt.Errorf("invalid markup section definition")
	}

	// validate tag
	if !contains(v.MarkupSections, tag) {
		return fmt.Errorf("invalid markup section tag")
	}

	// get items
	items, ok := section[2].(A)
	if !ok {
		return fmt.Errorf("invalid markup section definition")
	}

	// prepare open markup counter
	openMarkups := 0

	// validate markers
	for _, _marker := range items {
		// coerce value
		marker, ok := _marker.(A)
		if !ok {
			return fmt.Errorf("invalid markup section definition")
		}

		// validate marker
		var err error
		openMarkups, err = v.ValidateMarker(marker, numMarkups, numAtoms, openMarkups)
		if err != nil {
			return err
		}
	}

	return nil
}

// ValidateImageSection validates a single image section.
func (v *Validator) ValidateImageSection(image A) error {
	// check availability
	if v.ImageSection == nil {
		return fmt.Errorf("invalid image section")
	}

	// validate length
	if len(image) != 2 {
		return fmt.Errorf("invalid image section definition")
	}

	// get src
	src, ok := image[1].(string)
	if !ok {
		return fmt.Errorf("invalid image section definition")
	}

	// check src
	if !v.ImageSection(src) {
		return fmt.Errorf("invalid image section src")
	}

	return nil
}

// ValidateListSection validates a single list section.
func (v *Validator) ValidateListSection(list A, numMarkups, numAtoms int) error {
	// validate length
	if len(list) != 3 {
		return fmt.Errorf("invalid list section definition")
	}

	// get tag
	tag, ok := list[1].(string)
	if !ok {
		return fmt.Errorf("invalid list section definition")
	}

	// validate tag
	if !contains(v.ListSections, tag) {
		return fmt.Errorf("invalid list section tag")
	}

	// get items
	items, ok := list[2].(A)
	if !ok {
		return fmt.Errorf("invalid list section definition")
	}

	// validate items
	for _, _item := range items {
		// coerce value
		item, ok := _item.(A)
		if !ok {
			return fmt.Errorf("invalid list section definition")
		}

		// prepare open markup counter
		openMarkups := 0

		// validate markers
		for _, _marker := range item {
			// coerce value
			marker, ok := _marker.(A)
			if !ok {
				return fmt.Errorf("invalid list section definition")
			}

			// validate marker
			var err error
			openMarkups, err = v.ValidateMarker(marker, numMarkups, numAtoms, openMarkups)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// ValidateCardSection validates a single card section.
func (v *Validator) ValidateCardSection(card A, numCards int) error {
	// validate length
	if len(card) != 2 {
		return fmt.Errorf("invalid card section definition")
	}

	// get num
	num, ok := toInt(card[1])
	if !ok {
		return fmt.Errorf("invalid card section definition")
	}

	// check num
	if num >= numCards {
		return fmt.Errorf("invalid card section index")
	}

	return nil
}

// ValidateMarker validates a single marker.
func (v *Validator) ValidateMarker(marker A, numMarkups, numAtoms, openMarkups int) (int, error) {
	// validate length
	if len(marker) != 4 {
		return 0, fmt.Errorf("invalid marker definition")
	}

	// get marker type
	typ, ok := toInt(marker[0])
	if !ok {
		return 0, fmt.Errorf("invalid marker definition")
	}

	// check type
	if typ != TextMarker && typ != AtomMarker {
		return 0, fmt.Errorf("invalid marker definition")
	}

	// get opened markups
	openedMarkups, ok := marker[1].(A)
	if !ok {
		return 0, fmt.Errorf("invalid marker definition")
	}

	// validate opened markups
	for _, _markup := range openedMarkups {
		// coerce value
		markup, ok := toInt(_markup)
		if !ok {
			return 0, fmt.Errorf("invalid marker definition")
		}

		// check markup
		if markup >= numMarkups {
			return 0, fmt.Errorf("invalid marker markup index")
		}

		// increment counter
		openMarkups++
	}

	// get closed markups
	closedMarkups, ok := toInt(marker[2])
	if !ok {
		return 0, fmt.Errorf("invalid marker definition")
	}

	// decrement counter
	openMarkups -= closedMarkups
	if openMarkups < 0 {
		return 0, fmt.Errorf("invalid marker count")
	}

	// validate text marker
	if typ == TextMarker {
		if _, ok := marker[3].(string); !ok {
			return 0, fmt.Errorf("invalid marker definition")
		}
	}

	// validate atom marker
	if typ == AtomMarker {
		if atom, ok := toInt(marker[3]); !ok || atom >= numAtoms {
			return 0, fmt.Errorf("invalid marker atom index")
		}
	}

	return openMarkups, nil
}
