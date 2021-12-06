package mobiledoc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatValidator(t *testing.T) {
	v := NewFormatValidator()

	err := v.Validate(minimalDoc())
	assert.NoError(t, err)

	err = v.Validate(sampleDoc())
	assert.NoError(t, err)
}

func TestDefaultValidator(t *testing.T) {
	v := NewDefaultValidator()
	v.UnknownCards = true

	i := 0
	v.Atoms["atom1"] = func(name string, payload Map) bool {
		assert.Equal(t, "foo", name)
		assert.Equal(t, Map{"bar": 42.0}, payload)
		i++
		return true
	}
	v.Atoms["atom2"] = func(name string, payload Map) bool {
		assert.Equal(t, "foo", name)
		assert.Equal(t, Map{"bar": 24.0}, payload)
		i++
		return true
	}
	v.Cards["card1"] = func(payload Map) bool {
		assert.Equal(t, Map{"foo": 42.0}, payload)
		i++
		return true
	}
	v.Cards["card2"] = func(payload Map) bool {
		assert.Equal(t, Map{"foo": 24.0}, payload)
		i++
		return true
	}

	err := v.Validate(sampleDoc())
	assert.NoError(t, err)
	assert.Equal(t, 4, i)
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
