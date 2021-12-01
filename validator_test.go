package mobiledoc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidator(t *testing.T) {
	v := NewDefaultValidator()

	var atom1 Map
	v.Atoms["atom1"] = func(_ string, payload Map) bool {
		atom1 = payload
		return true
	}
	v.Atoms["atom2"] = func(string, Map) bool {
		return true
	}

	var card1 Map
	v.Cards["card1"] = func(payload Map) bool {
		card1 = payload
		return true
	}
	v.Cards["card2"] = func(payload Map) bool {
		return true
	}

	err := v.Validate(sampleDoc())
	assert.NoError(t, err)
	assert.Equal(t, Map{"bar": 42.0}, atom1)
	assert.Equal(t, Map{"foo": 42.0}, card1)
}

func TestNewEmptyValidator(t *testing.T) {
	NewEmptyValidator()
}

func TestValidatorInvalidVersion(t *testing.T) {
	v := NewDefaultValidator()

	err := v.Validate(Document{
		Version: "foo",
	})
	assert.Error(t, err)
}

func TestValidatorInvalidMarkup(t *testing.T) {
	v := NewDefaultValidator()

	err := v.Validate(Document{
		Version: Version,
		Markups: []Markup{
			{Tag: "x"},
		},
	})
	assert.Error(t, err)

	v.Markups["x"] = func(maps Map) bool {
		return false
	}

	err = v.Validate(Document{
		Version: Version,
		Markups: []Markup{
			{Tag: "x"},
		},
	})
	assert.Error(t, err)
}

func TestValidatorInvalidAtom(t *testing.T) {
	v := NewDefaultValidator()

	err := v.Validate(Document{
		Version: Version,
		Atoms: []Atom{
			{Name: "x"},
		},
	})
	assert.Error(t, err)

	v.Atoms["x"] = nil

	err = v.Validate(Document{
		Version: Version,
		Atoms: []Atom{
			{Name: "x"},
		},
	})
	assert.NoError(t, err)

	v.Atoms["x"] = func(string, Map) bool {
		return false
	}

	err = v.Validate(Document{
		Version: Version,
		Atoms: []Atom{
			{Name: "x"},
		},
	})
	assert.Error(t, err)
}

func TestValidatorInvalidCard(t *testing.T) {
	v := NewDefaultValidator()

	err := v.Validate(Document{
		Version: Version,
		Cards: []Card{
			{Name: "x"},
		},
	})
	assert.Error(t, err)

	v.Cards["x"] = nil

	err = v.Validate(Document{
		Version: Version,
		Cards: []Card{
			{Name: "x"},
		},
	})
	assert.NoError(t, err)

	v.Cards["x"] = func(Map) bool {
		return false
	}

	err = v.Validate(Document{
		Version: Version,
		Cards: []Card{
			{Name: "x"},
		},
	})
	assert.Error(t, err)
}

func TestValidatorInvalidMarkupSection(t *testing.T) {
	v := NewDefaultValidator()

	err := v.Validate(Document{
		Version: Version,
		Sections: []Section{
			{Type: MarkupSection, Tag: "x"},
		},
	})
	assert.Error(t, err)
}

func TestValidatorInvalidImageSection(t *testing.T) {
	v := NewDefaultValidator()

	v.ImageSection = nil

	err := v.Validate(Document{
		Version: Version,
		Sections: []Section{
			{Type: ImageSection},
		},
	})
	assert.Error(t, err)

	v.ImageSection = func(string) bool {
		return false
	}

	err = v.Validate(Document{
		Version: Version,
		Sections: []Section{
			{Type: ImageSection},
		},
	})
	assert.Error(t, err)
}

func TestValidatorInvalidListSection(t *testing.T) {
	v := NewDefaultValidator()

	err := v.Validate(Document{
		Version: Version,
		Sections: []Section{
			{Type: ListSection, Tag: "x"},
		},
	})
	assert.Error(t, err)
}
