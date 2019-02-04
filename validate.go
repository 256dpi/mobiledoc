package mobiledoc

import (
	"errors"

	"github.com/asaskevich/govalidator"
)

// TODO: Is the URL validation too strict?

// TODO: Allow "target" attribute for URLs?

// TODO: Use a validator type.

// TODO: Add parser that transforms the mobiledoc in a usable in memory layout.

// ErrInvalidMobileDoc is returned if the specified mobiledoc is invalid.
var ErrInvalidMobileDoc = errors.New("invalid mobiledoc")

// Validate will walk the specified mobiledoc and check if it is valid.
func Validate(doc M) error {
	// check version
	if version, ok := doc["version"]; !ok || version != "0.3.1" {
		return ErrInvalidMobileDoc
	}

	// prepare num markups
	numMarkups := 0

	// check markups
	if _markups, ok := doc["markups"]; ok {
		// coerce value
		markups, ok := _markups.(A)
		if !ok {
			return ErrInvalidMobileDoc
		}

		// set num
		numMarkups = len(markups)

		// validate markups
		for _, _markup := range markups {
			// coerce value
			markup, ok := _markup.(A)
			if !ok {
				return ErrInvalidMobileDoc
			}

			// validate markup
			err := ValidateMarkup(markup)
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
			return ErrInvalidMobileDoc
		}

		// set num
		numAtoms = len(atoms)

		// validate atom
		for _, _atom := range atoms {
			// coerce value
			atom, ok := _atom.(A)
			if !ok {
				return ErrInvalidMobileDoc
			}

			// validate atom
			err := ValidateAtom(atom)
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
			return ErrInvalidMobileDoc
		}

		// set num
		numCards = len(cards)

		// validate cards
		for _, _card := range cards {
			// coerce value
			card, ok := _card.(A)
			if !ok {
				return ErrInvalidMobileDoc
			}

			// validate card
			err := ValidateCard(card)
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
			return ErrInvalidMobileDoc
		}

		// validate sections
		for _, _section := range sections {
			// coerce value
			section, ok := _section.(A)
			if !ok {
				return ErrInvalidMobileDoc
			}

			// validate section
			err := ValidateSection(section, numMarkups, numAtoms, numCards)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// ValidateMarkup will validate a single markup.
func ValidateMarkup(markup A) error {
	// validate length
	if len(markup) == 0 || len(markup) > 2 {
		return ErrInvalidMobileDoc
	}

	// check tag
	tag, ok := markup[0].(string)
	if !ok {
		return ErrInvalidMobileDoc
	}

	// check tag existence
	allowedAttributes, ok := Markups[tag]
	if !ok {
		return ErrInvalidMobileDoc
	}

	// return if attributes are missing
	if len(markup) == 1 {
		return nil
	}

	// get attributes
	attributes, ok := markup[1].(A)
	if !ok {
		return ErrInvalidMobileDoc
	}

	// check if attributes are even
	if len(attributes)%2 != 0 {
		return ErrInvalidMobileDoc
	}

	// check attributes
	for i := 0; i < len(attributes); i += 2 {
		// coerce name
		name, ok1 := attributes[i].(string)
		value, ok2 := attributes[i+1].(string)
		if !ok1 || !ok2 {
			return ErrInvalidMobileDoc
		}

		// get validator
		validator, ok := allowedAttributes[name]
		if !ok {
			return ErrInvalidMobileDoc
		}

		// validate attribute
		if !validator(value) {
			return ErrInvalidMobileDoc
		}
	}

	return nil
}

// ValidateAtom will validate a single atom.
func ValidateAtom(atom A) error {
	// validate length
	if len(atom) == 0 || len(atom) > 3 {
		return ErrInvalidMobileDoc
	}

	// check name
	name, ok := atom[0].(string)
	if !ok {
		return ErrInvalidMobileDoc
	}

	// check atom existence
	validator, ok := Atoms[name]
	if !ok {
		return ErrInvalidMobileDoc
	}

	// prepare text and payload
	var text string
	var payload M

	// get text
	if len(atom) > 1 {
		text, ok = atom[1].(string)
		if !ok {
			return ErrInvalidMobileDoc
		}
	}

	// get payload
	if len(atom) > 2 {
		payload, ok = atom[2].(M)
		if !ok {
			return ErrInvalidMobileDoc
		}
	}

	// validate atom
	if !validator(text, payload) {
		return ErrInvalidMobileDoc
	}

	return nil
}

// ValidateCard will validate a single card.
func ValidateCard(card A) error {
	// validate length
	if len(card) == 0 || len(card) > 2 {
		return ErrInvalidMobileDoc
	}

	// check name
	name, ok := card[0].(string)
	if !ok {
		return ErrInvalidMobileDoc
	}

	// check card existence
	validator, ok := Cards[name]
	if !ok {
		return ErrInvalidMobileDoc
	}

	// prepare payload
	var payload M

	// get payload
	if len(card) > 1 {
		payload, ok = card[1].(M)
		if !ok {
			return ErrInvalidMobileDoc
		}
	}

	// validate card
	if !validator(payload) {
		return ErrInvalidMobileDoc
	}

	return nil
}

/*
[sectionTypeIdentifier, <type dependent>], ──── Card
[1, "p", [
      [textTypeIdentifier, openMarkupsIndexes, numberOfClosedMarkups, value],
      [0, [], 0, "Example with no markup"],      ──── textTypeIdentifier for markup is always 0
      [0, [0], 1, "Example wrapped in b tag (opened markup #0), 1 closed markup"],
      [0, [1], 0, "Example opening i tag (opened markup with #1, 0 closed markups)"],
      [0, [], 1, "Example closing i tag (no opened markups, 1 closed markup)"],
      [0, [1, 0], 1, "Example opening i tag and b tag, closing b tag (opened markups #1 and #0, 1 closed markup [closes markup #0])"],
      [0, [], 1, "Example closing i tag, (no opened markups, 1 closed markup [closes markup #1])"],
    ]],
    [1, "p", [
      [textTypeIdentifier, openMarkupsIndexes, numberOfClosedMarkups, atomIndex],
      [1, [], 0, 0],             ──── mention atom at index 0 (@bob), textTypeIdentifier for atom is always 1
      [1, [0], 1, 1]             ──── mention atom at index 1 (@tom) wrapped in b tag (markup index 0)
    ]]
*/

// ValidateSection will validate a single section.
func ValidateSection(section A, numMarkups, numAtoms, numCards int) error {
	// validate length
	if len(section) == 0 {
		return ErrInvalidMobileDoc
	}

	// get type
	typ, ok := section[0].(int)
	if !ok {
		return ErrInvalidMobileDoc
	}

	// run validators based on type
	switch typ {
	case MarkupSection:
		return ValidateMarkupSection(section, numMarkups, numAtoms)
	case ImageSection:
		return ValidateImageSection(section)
	case ListSection:
		return ValidateListSection(section, numMarkups, numAtoms)
	case CardSection:
		return ValidateCardSection(section, numCards)
	default:
		return ErrInvalidMobileDoc
	}
}

// ValidateMarkupSection validates a single markup section.
func ValidateMarkupSection(section A, numMarkups, numAtoms int) error {
	// validate length
	if len(section) != 3 {
		return ErrInvalidMobileDoc
	}

	// get tag
	tag, ok := section[1].(string)
	if !ok {
		return ErrInvalidMobileDoc
	}

	// validate tag
	if !contains(MarkupSections, tag) {
		return ErrInvalidMobileDoc
	}

	// get items
	items, ok := section[2].(A)
	if !ok {
		return ErrInvalidMobileDoc
	}

	// prepare open markup counter
	openMarkups := 0

	// validate markers
	for _, _marker := range items {
		// coerce value
		marker, ok := _marker.(A)
		if !ok {
			return ErrInvalidMobileDoc
		}

		// validate marker
		var err error
		openMarkups, err = ValidateMarker(marker, numMarkups, numAtoms, openMarkups)
		if err != nil {
			return err
		}
	}

	return nil
}

// ValidateImageSection validates a single image section.
func ValidateImageSection(image A) error {
	// TODO: Allow disable.

	// validate length
	if len(image) != 2 {
		return ErrInvalidMobileDoc
	}

	// get src
	src, ok := image[1].(string)
	if !ok {
		return ErrInvalidMobileDoc
	}

	// check src
	if !govalidator.IsURL(src) {
		return ErrInvalidMobileDoc
	}

	return nil
}

// ValidateListSection validates a single list section.
func ValidateListSection(list A, numMarkups, numAtoms int) error {
	// TODO: Allow disable.

	// validate length
	if len(list) != 3 {
		return ErrInvalidMobileDoc
	}

	// get tag
	tag, ok := list[1].(string)
	if !ok {
		return ErrInvalidMobileDoc
	}

	// validate tag
	if !contains(ListSections, tag) {
		return ErrInvalidMobileDoc
	}

	// get items
	items, ok := list[2].(A)
	if !ok {
		return ErrInvalidMobileDoc
	}

	// validate items
	for _, _item := range items {
		// coerce value
		item, ok := _item.(A)
		if !ok {
			return ErrInvalidMobileDoc
		}

		// prepare open markup counter
		openMarkups := 0

		// validate markers
		for _, _marker := range item {
			// coerce value
			marker, ok := _marker.(A)
			if !ok {
				return ErrInvalidMobileDoc
			}

			// validate marker
			var err error
			openMarkups, err = ValidateMarker(marker, numMarkups, numAtoms, openMarkups)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// ValidateCardSection validates a single card section.
func ValidateCardSection(card A, numCards int) error {
	// validate length
	if len(card) != 2 {
		return ErrInvalidMobileDoc
	}

	// get num
	num, ok := card[1].(int)
	if !ok {
		return ErrInvalidMobileDoc
	}

	// check num
	if num >= numCards {
		return ErrInvalidMobileDoc
	}

	return nil
}

// ValidateMarker validates a single marker.
func ValidateMarker(marker A, numMarkups, numAtoms, openMarkups int) (int, error) {
	// validate length
	if len(marker) != 4 {
		return 0, ErrInvalidMobileDoc
	}

	// get marker type
	typ, ok := marker[0].(int)
	if !ok {
		return 0, ErrInvalidMobileDoc
	}

	// check type
	if typ != TextMarker && typ != AtomMarker {
		return 0, ErrInvalidMobileDoc
	}

	// get opened markups
	openedMarkups, ok := marker[1].(A)
	if !ok {
		return 0, ErrInvalidMobileDoc
	}

	// validate opened markups
	for _, _markup := range openedMarkups {
		// coerce value
		markup, ok := _markup.(int)
		if !ok {
			return 0, ErrInvalidMobileDoc
		}

		// check markup
		if markup >= numMarkups {
			return 0, ErrInvalidMobileDoc
		}

		// increment counter
		openMarkups++
	}

	// get closed markups
	closedMarkups, ok := marker[2].(int)
	if !ok {
		return 0, ErrInvalidMobileDoc
	}

	// decrement counter
	openMarkups -= closedMarkups
	if openMarkups < 0 {
		return 0, ErrInvalidMobileDoc
	}

	// validate text marker
	if typ == TextMarker {
		if _, ok := marker[3].(string); !ok {
			return 0, ErrInvalidMobileDoc
		}
	}

	// validate atom marker
	if typ == AtomMarker {
		if atom, ok := marker[3].(int); !ok || atom >= numAtoms {
			return 0, ErrInvalidMobileDoc
		}
	}

	return openMarkups, nil
}
