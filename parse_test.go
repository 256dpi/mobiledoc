package mobiledoc

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateJSON(t *testing.T) {
	var in Map
	err := json.Unmarshal([]byte(`{
		"version":"0.3.1",
		"markups":[
			["b"],
			["i"],
			["a",["href","http://example.com"]]
		],
		"atoms":[
			["atom1","foo",{"bar":42}],
			["atom2","foo",{"bar":24}]
		],
		"cards":[
			["card1",{"foo":42}],
			["card2",{"foo":42}]
		],
		"sections":[
			[10,0],
			[1,"p",[
				[0,[],0,"foo"],
				[0,[0],1,"foo"],
				[0,[1],0,"foo"],
				[0,[],1,"foo"],
				[0,[1,2],1,"foo"],
				[0,[],1,"foo"]
			]],
			[1,"p",[
				[1,[],0,0],
				[1,[0],0,1],
				[1,[],1,0]
			]],
			[2,"http://example.com/foo.png"],
			[3,"ul",[
				[
					[0,[],0,"foo"],
					[0,[0],1,"foo"]
				],[
					[0,[0],0,"foo"],
					[0,[],1,"foo"]
				]
			]],
			[10,1]
		]
	}`), &in)
	assert.NoError(t, err)

	out := Document{
		Version: Version,
		Markups: []Markup{
			{Tag: "b"},
			{Tag: "i"},
			{Tag: "a", Attributes: Map{"href": "http://example.com"}},
		},
		Atoms: []Atom{
			{Name: "atom1", Text: "foo", Payload: Map{"bar": float64(42)}},
			{Name: "atom2", Text: "foo", Payload: Map{"bar": float64(24)}},
		},
		Cards: []Card{
			{Name: "card1", Payload: Map{"foo": float64(42)}},
			{Name: "card2", Payload: Map{"foo": float64(42)}},
		},
	}
	out.Sections = []Section{
		{Type: CardSection, Card: &out.Cards[0]},
		{Type: MarkupSection, Tag: "p", Markers: []Marker{
			{Type: TextMarker, Text: "foo"},
			{Type: TextMarker, OpenMarkups: []*Markup{&out.Markups[0]}, ClosedMarkups: 1, Text: "foo"},
			{Type: TextMarker, OpenMarkups: []*Markup{&out.Markups[1]}, Text: "foo"},
			{Type: TextMarker, ClosedMarkups: 1, Text: "foo"},
			{Type: TextMarker, OpenMarkups: []*Markup{&out.Markups[1], &out.Markups[2]}, ClosedMarkups: 1, Text: "foo"},
			{Type: TextMarker, ClosedMarkups: 1, Text: "foo"},
		}},
		{Type: MarkupSection, Tag: "p", Markers: []Marker{
			{Type: AtomMarker, Atom: &out.Atoms[0]},
			{Type: AtomMarker, OpenMarkups: []*Markup{&out.Markups[0]}, Atom: &out.Atoms[1]},
			{Type: AtomMarker, ClosedMarkups: 1, Atom: &out.Atoms[0]},
		}},
		{Type: ImageSection, Source: "http://example.com/foo.png"},
		{Type: ListSection, Tag: "ul", Items: [][]Marker{
			{
				{Type: TextMarker, ClosedMarkups: 0, Text: "foo"},
				{Type: TextMarker, OpenMarkups: []*Markup{&out.Markups[0]}, ClosedMarkups: 1, Text: "foo"},
			},
			{
				{Type: TextMarker, OpenMarkups: []*Markup{&out.Markups[0]}, Text: "foo"},
				{Type: TextMarker, ClosedMarkups: 1, Text: "foo"},
			},
		}},
		{Type: 10, Card: &out.Cards[1]},
	}

	doc, err := Parse(in)
	assert.NoError(t, err)
	assert.Equal(t, out, doc)
}

func TestParse(t *testing.T) {
	in := Map{
		"version": Version,
		"markups": List{
			List{"b"},
			List{"i"},
			List{"a", List{"href", "http://example.com"}},
		},
		"atoms": List{
			List{"atom1", "foo", Map{"bar": 42}},
			List{"atom2", "foo", Map{"bar": 24}},
		},
		"cards": List{
			List{"card1", Map{"foo": 42}},
			List{"card2", Map{"foo": 42}},
		},
		"sections": List{
			List{CardSection, 0},
			List{MarkupSection, "p", List{
				List{TextMarker, List{}, 0, "foo"},
				List{TextMarker, List{0}, 1, "foo"},
				List{TextMarker, List{1}, 0, "foo"},
				List{TextMarker, List{}, 1, "foo"},
				List{TextMarker, List{1, 2}, 1, "foo"},
				List{TextMarker, List{}, 1, "foo"},
			}},
			List{MarkupSection, "p", List{
				List{AtomMarker, List{}, 0, 0},
				List{AtomMarker, List{0}, 0, 1},
				List{AtomMarker, List{}, 1, 0},
			}},
			List{ImageSection, "http://example.com/foo.png"},
			List{ListSection, "ul", List{
				List{
					List{TextMarker, List{}, 0, "foo"},
					List{TextMarker, List{0}, 1, "foo"},
				},
				List{
					List{TextMarker, List{0}, 0, "foo"},
					List{TextMarker, List{}, 1, "foo"},
				},
			}},
			List{CardSection, 1},
		},
	}

	out := Document{
		Version: Version,
		Markups: []Markup{
			{Tag: "b"},
			{Tag: "i"},
			{Tag: "a", Attributes: Map{"href": "http://example.com"}},
		},
		Atoms: []Atom{
			{Name: "atom1", Text: "foo", Payload: Map{"bar": 42}},
			{Name: "atom2", Text: "foo", Payload: Map{"bar": 24}},
		},
		Cards: []Card{
			{Name: "card1", Payload: Map{"foo": 42}},
			{Name: "card2", Payload: Map{"foo": 42}},
		},
	}
	out.Sections = []Section{
		{Type: CardSection, Card: &out.Cards[0]},
		{Type: MarkupSection, Tag: "p", Markers: []Marker{
			{Type: TextMarker, Text: "foo"},
			{Type: TextMarker, OpenMarkups: []*Markup{&out.Markups[0]}, ClosedMarkups: 1, Text: "foo"},
			{Type: TextMarker, OpenMarkups: []*Markup{&out.Markups[1]}, Text: "foo"},
			{Type: TextMarker, ClosedMarkups: 1, Text: "foo"},
			{Type: TextMarker, OpenMarkups: []*Markup{&out.Markups[1], &out.Markups[2]}, ClosedMarkups: 1, Text: "foo"},
			{Type: TextMarker, ClosedMarkups: 1, Text: "foo"},
		}},
		{Type: MarkupSection, Tag: "p", Markers: []Marker{
			{Type: AtomMarker, Atom: &out.Atoms[0]},
			{Type: AtomMarker, OpenMarkups: []*Markup{&out.Markups[0]}, Atom: &out.Atoms[1]},
			{Type: AtomMarker, ClosedMarkups: 1, Atom: &out.Atoms[0]},
		}},
		{Type: ImageSection, Source: "http://example.com/foo.png"},
		{Type: ListSection, Tag: "ul", Items: [][]Marker{
			{
				{Type: TextMarker, ClosedMarkups: 0, Text: "foo"},
				{Type: TextMarker, OpenMarkups: []*Markup{&out.Markups[0]}, ClosedMarkups: 1, Text: "foo"},
			},
			{
				{Type: TextMarker, OpenMarkups: []*Markup{&out.Markups[0]}, Text: "foo"},
				{Type: TextMarker, ClosedMarkups: 1, Text: "foo"},
			},
		}},
		{Type: 10, Card: &out.Cards[1]},
	}

	doc, err := Parse(in)
	assert.NoError(t, err)
	assert.Equal(t, out, doc)
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
