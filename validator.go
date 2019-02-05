package mobiledoc

import "fmt"

// Version specifies the mobiledoc version.
const Version = "0.3.1"

// LinkValidator validates the href attribute.
func LinkValidator(attributes Map) bool {
	for key, value := range attributes {
		switch key {
		case "href":
			// get href
			_, ok := value.(string)
			if !ok {
				return false
			}
		default:
			return false
		}
	}

	return true
}

// NoAttributesValidator returns true if the provided attributes are empty.
func NoAttributesValidator(attributes Map) bool {
	return len(attributes) == 0
}

// DefaultMarkups defines the default expected markups with the tag as the key
// and a map of attributes and validator functions.
var DefaultMarkups = map[string]func(Map) bool{
	"a":      LinkValidator,
	"b":      NoAttributesValidator,
	"code":   NoAttributesValidator,
	"em":     NoAttributesValidator,
	"i":      NoAttributesValidator,
	"s":      NoAttributesValidator,
	"strong": NoAttributesValidator,
	"sub":    NoAttributesValidator,
	"sup":    NoAttributesValidator,
	"u":      NoAttributesValidator,
}

// DefaultMarkupSections defines the default markup sections.
var DefaultMarkupSections = []string{"aside", "blockquote", "h1", "h2", "h3", "h4", "h5", "h6", "p"}

// DefaultListSections defines the default list sections.
var DefaultListSections = []string{"ul", "ol"}

// DefaultImageSection defines the default image section validator.
var DefaultImageSection = func(source string) bool {
	return len(source) > 0
}

// Validator validates a mobiledoc.
type Validator struct {
	// Markups defines the allowed markups with the name as key and a
	// attributes validator function.
	Markups map[string]func(attributes Map) bool

	// Atoms defines the allowed atoms with the name as the key and a validator
	// function.
	Atoms map[string]func(name string, payload Map) bool

	// Cards defines the allowed cards with the name as the key and a validator
	// function.
	Cards map[string]func(payload Map) bool

	// MarkupSections defines the allowed markup sections.
	MarkupSections []string

	// ListSections defines the allowed list sections.
	ListSections []string

	// ImageSection defines whether the image section is allowed when a source
	// validator is set.
	ImageSection func(source string) bool
}

// NewEmptyValidator creates an empty validator.
func NewEmptyValidator() *Validator {
	return &Validator{
		Markups:        make(map[string]func(Map) bool),
		Atoms:          make(map[string]func(string, Map) bool),
		Cards:          make(map[string]func(Map) bool),
		MarkupSections: nil,
		ListSections:   nil,
		ImageSection:   nil,
	}
}

// NewDefaultValidator creates a validator that validates the default mobiledoc
// standard.
func NewDefaultValidator() *Validator {
	return &Validator{
		Markups:        DefaultMarkups,
		Atoms:          make(map[string]func(string, Map) bool),
		Cards:          make(map[string]func(Map) bool),
		MarkupSections: DefaultMarkupSections,
		ListSections:   DefaultListSections,
		ImageSection:   DefaultImageSection,
	}
}

// Validate will walk the specified mobiledoc and check if it is valid.
func (v *Validator) Validate(doc Document) error {
	// check version
	if doc.Version != Version {
		return fmt.Errorf("invalid version")
	}

	// validate markups
	for _, markup := range doc.Markups {
		err := v.validateMarkup(markup)
		if err != nil {
			return err
		}
	}

	// validate markups
	for _, atom := range doc.Atoms {
		err := v.validateAtom(atom)
		if err != nil {
			return err
		}
	}

	// validate cards
	for _, card := range doc.Cards {
		err := v.validateCard(card)
		if err != nil {
			return err
		}
	}

	// validate sections
	for _, section := range doc.Sections {
		err := v.validateSection(section)
		if err != nil {
			return err
		}
	}

	return nil
}

func (v *Validator) validateMarkup(markup Markup) error {
	// check markup allowance
	validator, ok := v.Markups[markup.Tag]
	if !ok {
		return fmt.Errorf("invalid markup tag")
	}

	// return if validator is missing
	if validator == nil {
		return nil
	}

	// validate attributes
	if !validator(markup.Attributes) {
		return fmt.Errorf("invalid markup attributes")
	}

	return nil
}

func (v *Validator) validateAtom(atom Atom) error {
	// check atom existence
	validator, ok := v.Atoms[atom.Name]
	if !ok {
		return fmt.Errorf("invalid atom name")
	}

	// check validator
	if validator == nil {
		return nil
	}

	// validate atom
	if !validator(atom.Text, atom.Payload) {
		return fmt.Errorf("invalid atom text or payload")
	}

	return nil
}

func (v *Validator) validateCard(card Card) error {
	// check card existence
	validator, ok := v.Cards[card.Name]
	if !ok {
		return fmt.Errorf("invalid card name")
	}

	// check validator
	if validator == nil {
		return nil
	}

	// validate card
	if !validator(card.Payload) {
		return fmt.Errorf("invalid card payload")
	}

	return nil
}

func (v *Validator) validateSection(section Section) error {
	// run validators based on type
	switch section.Type {
	case MarkupSection:
		return v.validateMarkupSection(section)
	case ImageSection:
		return v.validateImageSection(section)
	case ListSection:
		return v.validateListSection(section)
	}

	return nil
}

func (v *Validator) validateMarkupSection(section Section) error {
	// validate tag
	if !contains(v.MarkupSections, section.Tag) {
		return fmt.Errorf("invalid markup section tag")
	}

	return nil
}

func (v *Validator) validateImageSection(image Section) error {
	// check availability
	if v.ImageSection == nil {
		return fmt.Errorf("invalid image section")
	}

	// check src
	if !v.ImageSection(image.Source) {
		return fmt.Errorf("invalid image section src")
	}

	return nil
}

func (v *Validator) validateListSection(list Section) error {
	// validate tag
	if !contains(v.ListSections, list.Tag) {
		return fmt.Errorf("invalid list section tag")
	}

	return nil
}
