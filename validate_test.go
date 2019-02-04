package mobiledoc

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateJSON(t *testing.T) {
	var doc Map
	err := json.Unmarshal([]byte(`{
		"version": "0.3.1",
		"markups": [
			["b"],
			["i"]
		],
		"atoms": [
			["mention", "@bob", { "id": 42 }],
	    	["mention", "@tom", { "id": 12 }]
		],
	  	"sections": [
	    	[1, "p", [
	      	[0, [], 0, "Example"],
	      	[0, [0], 1, "Example"],
	      	[0, [1], 0, "Example"],
	      	[0, [], 1, "Example"],
	      	[0, [1, 0], 1, "Example"],
	      	[0, [], 1, "Example"]
	    ]],
	    [1, "p", [
			[1, [], 0, 0],
	      	[1, [0], 1, 1]
	    ]]
	  ]
	}`), &doc)
	assert.NoError(t, err)

	v := NewValidator()
	v.Atoms["mention"] = func(s string, ms Map) bool {
		return true
	}

	assert.NoError(t, v.Validate(doc))
}

func TestValidateDocument(t *testing.T) {
	v := NewValidator()

	assert.Error(t, v.Validate(Map{}))
	assert.Error(t, v.Validate(Map{"version": ""}))
	assert.Error(t, v.Validate(Map{"version": "3.1.2"}))
	assert.NoError(t, v.Validate(Map{"version": Version}))

	assert.Error(t, v.Validate(Map{"version": Version, "markups": "foo"}))
	assert.Error(t, v.Validate(Map{"version": Version, "atoms": "foo"}))
	assert.Error(t, v.Validate(Map{"version": Version, "sections": "foo"}))
	assert.Error(t, v.Validate(Map{"version": Version, "cards": "foo"}))
	assert.NoError(t, v.Validate(Map{"version": Version, "markups": List{}}))
	assert.NoError(t, v.Validate(Map{"version": Version, "atoms": List{}}))
	assert.NoError(t, v.Validate(Map{"version": Version, "sections": List{}}))
	assert.NoError(t, v.Validate(Map{"version": Version, "cards": List{}}))

	v.Atoms["mention"] = func(s string, ms Map) bool {
		return true
	}

	doc := Map{
		"version": Version,
		"markups": List{
			List{"b"},
			List{"i"},
		},
		"atoms": List{
			List{"mention", "@bob", Map{"id": 42}},
			List{"mention", "@tom", Map{"id": 24}},
		},
		"sections": List{
			List{1, "p", List{
				List{TextMarker, List{}, 0, "Example"},
				List{TextMarker, List{0}, 1, "Example"},
				List{TextMarker, List{1}, 0, "Example"},
				List{TextMarker, List{}, 1, "Example"},
				List{TextMarker, List{1, 0}, 1, "Example"},
				List{TextMarker, List{}, 1, "Example"},
			}},
			List{1, "p", List{
				List{AtomMarker, List{}, 0, 0},
				List{AtomMarker, List{0}, 1, 1},
			}},
		},
	}
	assert.NoError(t, v.Validate(doc))
}

func TestValidateMarkups(t *testing.T) {
	v := NewValidator()

	assert.NoError(t, v.ValidateMarkup(List{"em"}))
	assert.NoError(t, v.ValidateMarkup(List{"strong", List{}}))
	assert.NoError(t, v.ValidateMarkup(List{"a", List{"href", "http://example.com"}}))
	assert.NoError(t, v.ValidateMarkup(List{"a", List{"href", "mailto:example@example.com"}}))

	assert.Error(t, v.ValidateMarkup(List{1}))
	assert.Error(t, v.ValidateMarkup(List{"foo"}))
	assert.Error(t, v.ValidateMarkup(List{"a", List{"href"}}))
	assert.Error(t, v.ValidateMarkup(List{"a", List{"foo", "bar"}}))
	assert.Error(t, v.ValidateMarkup(List{"strong", List{"href", "http://example.com"}}))
	assert.Error(t, v.ValidateMarkup(List{"strong", List{"href", "foo"}}))
}

func TestValidateAtom(t *testing.T) {
	v := NewValidator()

	var lastText string
	var lastPayload Map
	v.Atoms["foo"] = func(text string, payload Map) bool {
		lastText = text
		lastPayload = payload
		return true
	}

	assert.Error(t, v.ValidateAtom(List{"bar", "bar"}))

	assert.NoError(t, v.ValidateAtom(List{"foo", "bar"}))
	assert.Equal(t, "bar", lastText)
	assert.Equal(t, Map(nil), lastPayload)

	assert.NoError(t, v.ValidateAtom(List{"foo", "bar", Map{"baz": "qux"}}))
	assert.Equal(t, "bar", lastText)
	assert.Equal(t, Map{"baz": "qux"}, lastPayload)
}

func TestValidateCard(t *testing.T) {
	v := NewValidator()

	var lastPayload Map
	v.Cards["foo"] = func(payload Map) bool {
		lastPayload = payload
		return true
	}

	assert.Error(t, v.ValidateCard(List{"bar"}))

	assert.NoError(t, v.ValidateCard(List{"foo"}))
	assert.Equal(t, Map(nil), lastPayload)

	assert.NoError(t, v.ValidateCard(List{"foo", Map{"bar": "baz"}}))
	assert.Equal(t, Map{"bar": "baz"}, lastPayload)
}

func TestValidateSection(t *testing.T) {
	v := NewValidator()

	assert.Error(t, v.ValidateMarkupSection(List{9}, 0, 0))
}

func TestValidateMarkupSection(t *testing.T) {
	v := NewValidator()

	assert.Error(t, v.ValidateMarkupSection(List{MarkupSection, "p", 0}, 0, 0))
	assert.Error(t, v.ValidateMarkupSection(List{MarkupSection, 0, List{}}, 0, 0))

	assert.NoError(t, v.ValidateMarkupSection(List{MarkupSection, "h1", List{
		List{TextMarker, List{0}, 0, "foo"},
		List{AtomMarker, List{1}, 0, 0},
		List{TextMarker, List{0}, 2, "bar"},
	}}, 2, 1))
}

func TestValidateImageSection(t *testing.T) {
	v := NewValidator()

	assert.NoError(t, v.ValidateImageSection(List{ImageSection, "foo"}))
	assert.NoError(t, v.ValidateImageSection(List{ImageSection, "http://example.com/foo.png"}))
}

func TestValidateListSection(t *testing.T) {
	v := NewValidator()

	assert.Error(t, v.ValidateListSection(List{ListSection, "ol", 0}, 0, 0))
	assert.Error(t, v.ValidateListSection(List{ListSection, "ul", List{
		List{
			List{TextMarker, List{0}, 0, "foo"},
		},
	}}, 0, 0))

	assert.NoError(t, v.ValidateListSection(List{ListSection, "ol", List{
		List{
			List{TextMarker, List{}, 0, "foo"},
		},
	}}, 0, 0))
	assert.NoError(t, v.ValidateListSection(List{ListSection, "ol", List{
		List{
			List{TextMarker, List{0}, 0, "foo"},
			List{AtomMarker, List{1}, 0, 0},
			List{TextMarker, List{0}, 2, "bar"},
		},
	}}, 2, 1))
}

func TestValidateCardSection(t *testing.T) {
	v := NewValidator()

	assert.Error(t, v.ValidateCardSection(List{CardSection, 0}, 0))
	assert.NoError(t, v.ValidateCardSection(List{CardSection, 0}, 1))
}

func TestValidateMarker(t *testing.T) {
	v := NewValidator()

	openMarkups, err := v.ValidateMarker(List{TextMarker, List{0}, 0, "foo"}, 0, 0, 0)
	assert.Error(t, err)
	assert.Equal(t, 0, openMarkups)

	openMarkups, err = v.ValidateMarker(List{TextMarker, List{0}, 0, "foo"}, 1, 0, 0)
	assert.NoError(t, err)
	assert.Equal(t, 1, openMarkups)

	openMarkups, err = v.ValidateMarker(List{AtomMarker, List{}, 1, 0}, 0, 1, 1)
	assert.NoError(t, err)
	assert.Equal(t, 0, openMarkups)
}
