package mobiledoc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateDocument(t *testing.T) {
	assert.Error(t, Validate(M{}))
	assert.Error(t, Validate(M{"version": ""}))
	assert.Error(t, Validate(M{"version": "3.1.2"}))
	assert.NoError(t, Validate(M{"version": "0.3.1"}))

	assert.Error(t, Validate(M{"version": "0.3.1", "markups": "foo"}))
	assert.Error(t, Validate(M{"version": "0.3.1", "atoms": "foo"}))
	assert.Error(t, Validate(M{"version": "0.3.1", "sections": "foo"}))
	assert.Error(t, Validate(M{"version": "0.3.1", "cards": "foo"}))
	assert.NoError(t, Validate(M{"version": "0.3.1", "markups": A{}}))
	assert.NoError(t, Validate(M{"version": "0.3.1", "atoms": A{}}))
	assert.NoError(t, Validate(M{"version": "0.3.1", "sections": A{}}))
	assert.NoError(t, Validate(M{"version": "0.3.1", "cards": A{}}))

	// TODO: Test full document.
}

func TestValidateMarkups(t *testing.T) {
	assert.NoError(t, ValidateMarkup(A{"em"}))
	assert.NoError(t, ValidateMarkup(A{"strong", A{}}))
	assert.NoError(t, ValidateMarkup(A{"a", A{"href", "http://example.com"}}))
	assert.NoError(t, ValidateMarkup(A{"a", A{"href", "mailto:example@example.com"}}))

	assert.Error(t, ValidateMarkup(A{"foo"}))
	assert.Error(t, ValidateMarkup(A{"a", A{"href"}}))
	assert.Error(t, ValidateMarkup(A{"a", A{"foo", "bar"}}))
	assert.Error(t, ValidateMarkup(A{"strong", A{"href", "http://example.com"}}))
	assert.Error(t, ValidateMarkup(A{"strong", A{"href", "foo"}}))
}

func TestValidateAtom(t *testing.T) {
	var lastText string
	var lastPayload M
	Atoms["foo"] = func(text string, payload M) bool {
		lastText = text
		lastPayload = payload
		return true
	}

	assert.Error(t, ValidateAtom(A{"bar", "bar"}))

	assert.NoError(t, ValidateAtom(A{"foo", "bar"}))
	assert.Equal(t, "bar", lastText)
	assert.Equal(t, M(nil), lastPayload)

	assert.NoError(t, ValidateAtom(A{"foo", "bar", M{"baz": "qux"}}))
	assert.Equal(t, "bar", lastText)
	assert.Equal(t, M{"baz": "qux"}, lastPayload)
}

func TestValidateCard(t *testing.T) {
	var lastPayload M
	Cards["foo"] = func(payload M) bool {
		lastPayload = payload
		return true
	}

	assert.Error(t, ValidateCard(A{"bar"}))

	assert.NoError(t, ValidateCard(A{"foo"}))
	assert.Equal(t, M(nil), lastPayload)

	assert.NoError(t, ValidateCard(A{"foo", M{"bar": "baz"}}))
	assert.Equal(t, M{"bar": "baz"}, lastPayload)
}

func TestValidateSection(t *testing.T) {
	assert.Error(t, ValidateMarkupSection(A{9}, 0, 0))
}

func TestValidateMarkupSection(t *testing.T) {
	assert.Error(t, ValidateMarkupSection(A{MarkupSection, "p", 0}, 0, 0))
	assert.Error(t, ValidateMarkupSection(A{MarkupSection, 0, A{}}, 0, 0))

	assert.NoError(t, ValidateMarkupSection(A{MarkupSection, "h1", A{
		A{TextMarker, A{0}, 0, "foo"},
		A{AtomMarker, A{1}, 0, 0},
		A{TextMarker, A{0}, 2, "bar"},
	}}, 2, 1))
}

func TestValidateImageSection(t *testing.T) {
	assert.Error(t, ValidateImageSection(A{ImageSection, "foo"}))
	assert.NoError(t, ValidateImageSection(A{ImageSection, "http://example.com/foo.png"}))
}

func TestValidateListSection(t *testing.T) {
	assert.Error(t, ValidateListSection(A{ListSection, "ol", 0}, 0, 0))
	assert.Error(t, ValidateListSection(A{ListSection, "ul", A{
		A{
			A{TextMarker, A{0}, 0, "foo"},
		},
	}}, 0, 0))

	assert.NoError(t, ValidateListSection(A{ListSection, "ol", A{
		A{
			A{TextMarker, A{}, 0, "foo"},
		},
	}}, 0, 0))
	assert.NoError(t, ValidateListSection(A{ListSection, "ol", A{
		A{
			A{TextMarker, A{0}, 0, "foo"},
			A{AtomMarker, A{1}, 0, 0},
			A{TextMarker, A{0}, 2, "bar"},
		},
	}}, 2, 1))
}

func TestValidateCardSection(t *testing.T) {
	assert.Error(t, ValidateCardSection(A{CardSection, 0}, 0))
	assert.NoError(t, ValidateCardSection(A{CardSection, 0}, 1))
}

func TestValidateMarker(t *testing.T) {
	openMarkups, err := ValidateMarker(A{TextMarker, A{0}, 0, "foo"}, 0, 0, 0)
	assert.Error(t, err)
	assert.Equal(t, 0, openMarkups)

	openMarkups, err = ValidateMarker(A{TextMarker, A{0}, 0, "foo"}, 1, 0, 0)
	assert.NoError(t, err)
	assert.Equal(t, 1, openMarkups)

	openMarkups, err = ValidateMarker(A{AtomMarker, A{}, 1, 0}, 0, 1, 1)
	assert.NoError(t, err)
	assert.Equal(t, 0, openMarkups)
}
