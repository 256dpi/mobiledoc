package mobiledoc

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateJSON(t *testing.T) {
	var doc M
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
	v.Atoms["mention"] = func(s string, ms M) bool {
		return true
	}

	assert.NoError(t, v.Validate(doc))
}

func TestValidateDocument(t *testing.T) {
	v := NewValidator()

	assert.Error(t, v.Validate(M{}))
	assert.Error(t, v.Validate(M{"version": ""}))
	assert.Error(t, v.Validate(M{"version": "3.1.2"}))
	assert.NoError(t, v.Validate(M{"version": "0.3.1"}))

	assert.Error(t, v.Validate(M{"version": "0.3.1", "markups": "foo"}))
	assert.Error(t, v.Validate(M{"version": "0.3.1", "atoms": "foo"}))
	assert.Error(t, v.Validate(M{"version": "0.3.1", "sections": "foo"}))
	assert.Error(t, v.Validate(M{"version": "0.3.1", "cards": "foo"}))
	assert.NoError(t, v.Validate(M{"version": "0.3.1", "markups": A{}}))
	assert.NoError(t, v.Validate(M{"version": "0.3.1", "atoms": A{}}))
	assert.NoError(t, v.Validate(M{"version": "0.3.1", "sections": A{}}))
	assert.NoError(t, v.Validate(M{"version": "0.3.1", "cards": A{}}))

	v.Atoms["mention"] = func(s string, ms M) bool {
		return true
	}

	doc := M{
		"version": "0.3.1",
		"markups": A{
			A{"b"},
			A{"i"},
		},
		"atoms": A{
			A{"mention", "@bob", M{"id": 42}},
			A{"mention", "@tom", M{"id": 24}},
		},
		"sections": A{
			A{1, "p", A{
				A{TextMarker, A{}, 0, "Example"},
				A{TextMarker, A{0}, 1, "Example"},
				A{TextMarker, A{1}, 0, "Example"},
				A{TextMarker, A{}, 1, "Example"},
				A{TextMarker, A{1, 0}, 1, "Example"},
				A{TextMarker, A{}, 1, "Example"},
			}},
			A{1, "p", A{
				A{AtomMarker, A{}, 0, 0},
				A{AtomMarker, A{0}, 1, 1},
			}},
		},
	}
	assert.NoError(t, v.Validate(doc))
}

func TestValidateMarkups(t *testing.T) {
	v := NewValidator()

	assert.NoError(t, v.ValidateMarkup(A{"em"}))
	assert.NoError(t, v.ValidateMarkup(A{"strong", A{}}))
	assert.NoError(t, v.ValidateMarkup(A{"a", A{"href", "http://example.com"}}))
	assert.NoError(t, v.ValidateMarkup(A{"a", A{"href", "mailto:example@example.com"}}))

	assert.Error(t, v.ValidateMarkup(A{"foo"}))
	assert.Error(t, v.ValidateMarkup(A{"a", A{"href"}}))
	assert.Error(t, v.ValidateMarkup(A{"a", A{"foo", "bar"}}))
	assert.Error(t, v.ValidateMarkup(A{"strong", A{"href", "http://example.com"}}))
	assert.Error(t, v.ValidateMarkup(A{"strong", A{"href", "foo"}}))
}

func TestValidateAtom(t *testing.T) {
	v := NewValidator()

	var lastText string
	var lastPayload M
	v.Atoms["foo"] = func(text string, payload M) bool {
		lastText = text
		lastPayload = payload
		return true
	}

	assert.Error(t, v.ValidateAtom(A{"bar", "bar"}))

	assert.NoError(t, v.ValidateAtom(A{"foo", "bar"}))
	assert.Equal(t, "bar", lastText)
	assert.Equal(t, M(nil), lastPayload)

	assert.NoError(t, v.ValidateAtom(A{"foo", "bar", M{"baz": "qux"}}))
	assert.Equal(t, "bar", lastText)
	assert.Equal(t, M{"baz": "qux"}, lastPayload)
}

func TestValidateCard(t *testing.T) {
	v := NewValidator()

	var lastPayload M
	v.Cards["foo"] = func(payload M) bool {
		lastPayload = payload
		return true
	}

	assert.Error(t, v.ValidateCard(A{"bar"}))

	assert.NoError(t, v.ValidateCard(A{"foo"}))
	assert.Equal(t, M(nil), lastPayload)

	assert.NoError(t, v.ValidateCard(A{"foo", M{"bar": "baz"}}))
	assert.Equal(t, M{"bar": "baz"}, lastPayload)
}

func TestValidateSection(t *testing.T) {
	v := NewValidator()

	assert.Error(t, v.ValidateMarkupSection(A{9}, 0, 0))
}

func TestValidateMarkupSection(t *testing.T) {
	v := NewValidator()

	assert.Error(t, v.ValidateMarkupSection(A{MarkupSection, "p", 0}, 0, 0))
	assert.Error(t, v.ValidateMarkupSection(A{MarkupSection, 0, A{}}, 0, 0))

	assert.NoError(t, v.ValidateMarkupSection(A{MarkupSection, "h1", A{
		A{TextMarker, A{0}, 0, "foo"},
		A{AtomMarker, A{1}, 0, 0},
		A{TextMarker, A{0}, 2, "bar"},
	}}, 2, 1))
}

func TestValidateImageSection(t *testing.T) {
	v := NewValidator()

	assert.NoError(t, v.ValidateImageSection(A{ImageSection, "foo"}))
	assert.NoError(t, v.ValidateImageSection(A{ImageSection, "http://example.com/foo.png"}))
}

func TestValidateListSection(t *testing.T) {
	v := NewValidator()

	assert.Error(t, v.ValidateListSection(A{ListSection, "ol", 0}, 0, 0))
	assert.Error(t, v.ValidateListSection(A{ListSection, "ul", A{
		A{
			A{TextMarker, A{0}, 0, "foo"},
		},
	}}, 0, 0))

	assert.NoError(t, v.ValidateListSection(A{ListSection, "ol", A{
		A{
			A{TextMarker, A{}, 0, "foo"},
		},
	}}, 0, 0))
	assert.NoError(t, v.ValidateListSection(A{ListSection, "ol", A{
		A{
			A{TextMarker, A{0}, 0, "foo"},
			A{AtomMarker, A{1}, 0, 0},
			A{TextMarker, A{0}, 2, "bar"},
		},
	}}, 2, 1))
}

func TestValidateCardSection(t *testing.T) {
	v := NewValidator()

	assert.Error(t, v.ValidateCardSection(A{CardSection, 0}, 0))
	assert.NoError(t, v.ValidateCardSection(A{CardSection, 0}, 1))
}

func TestValidateMarker(t *testing.T) {
	v := NewValidator()

	openMarkups, err := v.ValidateMarker(A{TextMarker, A{0}, 0, "foo"}, 0, 0, 0)
	assert.Error(t, err)
	assert.Equal(t, 0, openMarkups)

	openMarkups, err = v.ValidateMarker(A{TextMarker, A{0}, 0, "foo"}, 1, 0, 0)
	assert.NoError(t, err)
	assert.Equal(t, 1, openMarkups)

	openMarkups, err = v.ValidateMarker(A{AtomMarker, A{}, 1, 0}, 0, 1, 1)
	assert.NoError(t, err)
	assert.Equal(t, 0, openMarkups)
}
