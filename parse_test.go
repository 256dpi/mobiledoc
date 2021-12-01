package mobiledoc

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseJSON(t *testing.T) {
	var in Map
	err := json.Unmarshal([]byte(sampleJSON), &in)
	assert.NoError(t, err)

	doc, err := Parse(in)
	assert.NoError(t, err)
	assert.Equal(t, sampleDoc(), doc)
}

func TestParseMap(t *testing.T) {
	doc, err := Parse(sampleMap())
	assert.NoError(t, err)
	assert.Equal(t, sampleDoc(), doc)
}

func TestParseInvalidDocument(t *testing.T) {
	_, err := Parse(Map{
		"version": 1,
	})
	assert.Error(t, err)

	_, err = Parse(Map{
		"version": Version,
		"markups": 1,
	})
	assert.Error(t, err)

	_, err = Parse(Map{
		"version": Version,
		"markups": List{1},
	})
	assert.Error(t, err)

	_, err = Parse(Map{
		"version": Version,
		"atoms":   1,
	})
	assert.Error(t, err)

	_, err = Parse(Map{
		"version": Version,
		"atoms":   List{1},
	})
	assert.Error(t, err)

	_, err = Parse(Map{
		"version": Version,
		"cards":   1,
	})
	assert.Error(t, err)

	_, err = Parse(Map{
		"version": Version,
		"cards":   List{1},
	})
	assert.Error(t, err)

	_, err = Parse(Map{
		"version":  Version,
		"sections": 1,
	})
	assert.Error(t, err)

	_, err = Parse(Map{
		"version":  Version,
		"sections": List{1},
	})
	assert.Error(t, err)
}

func TestParseInvalidMarkups(t *testing.T) {
	_, err := Parse(Map{
		"version": Version,
		"markups": List{
			List{},
		},
	})
	assert.Error(t, err)

	_, err = Parse(Map{
		"version": Version,
		"markups": List{
			List{1},
		},
	})
	assert.Error(t, err)

	_, err = Parse(Map{
		"version": Version,
		"markups": List{
			List{"b", 1},
		},
	})
	assert.Error(t, err)

	_, err = Parse(Map{
		"version": Version,
		"markups": List{
			List{"b", List{1}},
		},
	})
	assert.Error(t, err)

	_, err = Parse(Map{
		"version": Version,
		"markups": List{
			List{"b", List{1, 1}},
		},
	})
	assert.Error(t, err)
}

func TestParseInvalidAtom(t *testing.T) {
	_, err := Parse(Map{
		"version": Version,
		"atoms": List{
			List{1},
		},
	})
	assert.Error(t, err)

	_, err = Parse(Map{
		"version": Version,
		"atoms": List{
			List{1, 1, 1},
		},
	})
	assert.Error(t, err)

	_, err = Parse(Map{
		"version": Version,
		"atoms": List{
			List{"atom", 1, 1},
		},
	})
	assert.Error(t, err)

	_, err = Parse(Map{
		"version": Version,
		"atoms": List{
			List{"atom", "foo", 1},
		},
	})
	assert.Error(t, err)
}

func TestParseInvalidCard(t *testing.T) {
	_, err := Parse(Map{
		"version": Version,
		"cards": List{
			List{1},
		},
	})
	assert.Error(t, err)

	_, err = Parse(Map{
		"version": Version,
		"cards": List{
			List{1, 1},
		},
	})
	assert.Error(t, err)

	_, err = Parse(Map{
		"version": Version,
		"cards": List{
			List{"foo", 1},
		},
	})
	assert.Error(t, err)
}

func TestParseInvalidSection(t *testing.T) {
	_, err := Parse(Map{
		"version": Version,
		"sections": List{
			List{},
		},
	})
	assert.Error(t, err)

	_, err = Parse(Map{
		"version": Version,
		"sections": List{
			List{false},
		},
	})
	assert.Error(t, err)

	_, err = Parse(Map{
		"version": Version,
		"sections": List{
			List{-1},
		},
	})
	assert.Error(t, err)
}

func TestParseInvalidMarkupSection(t *testing.T) {
	_, err := Parse(Map{
		"version": Version,
		"sections": List{
			List{MarkupSection},
		},
	})
	assert.Error(t, err)

	_, err = Parse(Map{
		"version": Version,
		"sections": List{
			List{MarkupSection, 1, 1},
		},
	})
	assert.Error(t, err)

	_, err = Parse(Map{
		"version": Version,
		"sections": List{
			List{MarkupSection, "p", 1},
		},
	})
	assert.Error(t, err)

	_, err = Parse(Map{
		"version": Version,
		"sections": List{
			List{MarkupSection, "p", List{
				1,
			}},
		},
	})
	assert.Error(t, err)

	_, err = Parse(Map{
		"version": Version,
		"sections": List{
			List{MarkupSection, "p", List{
				List{},
			}},
		},
	})
	assert.Error(t, err)
}

func TestParseInvalidImageSection(t *testing.T) {
	_, err := Parse(Map{
		"version": Version,
		"sections": List{
			List{ImageSection},
		},
	})
	assert.Error(t, err)

	_, err = Parse(Map{
		"version": Version,
		"sections": List{
			List{ImageSection, 1},
		},
	})
	assert.Error(t, err)
}

func TestParseInvalidListSection(t *testing.T) {
	_, err := Parse(Map{
		"version": Version,
		"sections": List{
			List{ListSection},
		},
	})
	assert.Error(t, err)

	_, err = Parse(Map{
		"version": Version,
		"sections": List{
			List{ListSection, 1, 1},
		},
	})
	assert.Error(t, err)

	_, err = Parse(Map{
		"version": Version,
		"sections": List{
			List{ListSection, "ol", 1},
		},
	})
	assert.Error(t, err)

	_, err = Parse(Map{
		"version": Version,
		"sections": List{
			List{ListSection, "ol", List{
				1,
			}},
		},
	})
	assert.Error(t, err)

	_, err = Parse(Map{
		"version": Version,
		"sections": List{
			List{ListSection, "ol", List{
				List{1},
			}},
		},
	})
	assert.Error(t, err)

	_, err = Parse(Map{
		"version": Version,
		"sections": List{
			List{ListSection, "ol", List{
				List{
					List{},
				},
			}},
		},
	})
	assert.Error(t, err)
}

func TestParseInvalidCardSection(t *testing.T) {
	_, err := Parse(Map{
		"version": Version,
		"sections": List{
			List{CardSection},
		},
	})
	assert.Error(t, err)

	_, err = Parse(Map{
		"version": Version,
		"sections": List{
			List{CardSection, false},
		},
	})
	assert.Error(t, err)

	_, err = Parse(Map{
		"version": Version,
		"sections": List{
			List{CardSection, 1},
		},
	})
	assert.Error(t, err)
}

func TestParseInvalidMarker(t *testing.T) {
	_, err := Parse(Map{
		"version": Version,
		"sections": List{
			List{MarkupSection, "p", List{
				List{false, 1, 1, 1},
			}},
		},
	})
	assert.Error(t, err)

	_, err = Parse(Map{
		"version": Version,
		"sections": List{
			List{MarkupSection, "p", List{
				List{-1, 1, 1, 1},
			}},
		},
	})
	assert.Error(t, err)

	_, err = Parse(Map{
		"version": Version,
		"sections": List{
			List{MarkupSection, "p", List{
				List{TextMarker, 1, 1, 1},
			}},
		},
	})
	assert.Error(t, err)

	_, err = Parse(Map{
		"version": Version,
		"sections": List{
			List{MarkupSection, "p", List{
				List{TextMarker, List{false}, false, 1},
			}},
		},
	})
	assert.Error(t, err)

	_, err = Parse(Map{
		"version": Version,
		"sections": List{
			List{MarkupSection, "p", List{
				List{TextMarker, List{0}, false, 1},
			}},
		},
	})
	assert.Error(t, err)

	_, err = Parse(Map{
		"version": Version,
		"sections": List{
			List{MarkupSection, "p", List{
				List{TextMarker, List{}, false, 1},
			}},
		},
	})
	assert.Error(t, err)

	_, err = Parse(Map{
		"version": Version,
		"sections": List{
			List{MarkupSection, "p", List{
				List{TextMarker, List{}, 1, 1},
			}},
		},
	})
	assert.Error(t, err)

	_, err = Parse(Map{
		"version": Version,
		"sections": List{
			List{MarkupSection, "p", List{
				List{TextMarker, List{}, 0, 1},
			}},
		},
	})
	assert.Error(t, err)

	_, err = Parse(Map{
		"version": Version,
		"sections": List{
			List{MarkupSection, "p", List{
				List{AtomMarker, List{}, 0, 1},
			}},
		},
	})
	assert.Error(t, err)
}

func BenchmarkParse(b *testing.B) {
	in := sampleMap()
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := Parse(in)
		if err != nil {
			panic(err)
		}
	}
}
