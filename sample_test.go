package mobiledoc

const sampleJSON = `{
	"version":"0.3.1",
	"markups":[
		["b"],
		["i"],
		["a",["href","https://example.com"]]
	],
	"atoms":[
		["atom1","foo",{"bar":42}],
		["atom2","foo",{"bar":24}]
	],
	"cards":[
		["card1",{"foo":42}],
		["card2",{"foo":24}]
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
		[2,"https://example.com/foo.png"],
		[3,"ul",[
			[
				[0,[],0,"foo"],
				[0,[0],1,"foo"]
			],[
				[0,[0],0,"foo"],
				[0,[],1,"<foo>"]
			]
		]],
		[3,"ol",[
			[
				[0,[],0,"bar"],
				[0,[1],1,"bar"]
			],[
				[0,[1],0,"bar"],
				[0,[],1,"<bar>"]
			]
		]],
		[10,1]
	]
}`

func sampleMap() Map {
	return Map{
		"version": Version,
		"markups": List{
			List{"b"},
			List{"i"},
			List{"a", List{"href", "https://example.com"}},
		},
		"atoms": List{
			List{"atom1", "foo", Map{"bar": 42.0}},
			List{"atom2", "foo", Map{"bar": 24.0}},
		},
		"cards": List{
			List{"card1", Map{"foo": 42.0}},
			List{"card2", Map{"foo": 24.0}},
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
			List{ImageSection, "https://example.com/foo.png"},
			List{ListSection, "ul", List{
				List{
					List{TextMarker, List{}, 0, "foo"},
					List{TextMarker, List{0}, 1, "foo"},
				},
				List{
					List{TextMarker, List{0}, 0, "foo"},
					List{TextMarker, List{}, 1, "<foo>"},
				},
			}},
			List{ListSection, "ol", List{
				List{
					List{TextMarker, List{}, 0, "bar"},
					List{TextMarker, List{1}, 1, "bar"},
				},
				List{
					List{TextMarker, List{1}, 0, "bar"},
					List{TextMarker, List{}, 1, "<bar>"},
				},
			}},
			List{CardSection, 1},
		},
	}
}

func sampleDoc() Document {
	doc := Document{
		Version: Version,
		Markups: []Markup{
			{Tag: "b"},
			{Tag: "i"},
			{Tag: "a", Attributes: Map{"href": "https://example.com"}},
		},
		Atoms: []Atom{
			{Name: "atom1", Text: "foo", Payload: Map{"bar": float64(42)}},
			{Name: "atom2", Text: "foo", Payload: Map{"bar": float64(24)}},
		},
		Cards: []Card{
			{Name: "card1", Payload: Map{"foo": float64(42)}},
			{Name: "card2", Payload: Map{"foo": float64(24)}},
		},
	}
	doc.Sections = []Section{
		{Type: CardSection, Card: &doc.Cards[0]},
		{Type: MarkupSection, Tag: "p", Markers: []Marker{
			{Type: TextMarker, Text: "foo"},
			{Type: TextMarker, OpenMarkups: []*Markup{&doc.Markups[0]}, ClosedMarkups: 1, Text: "foo"},
			{Type: TextMarker, OpenMarkups: []*Markup{&doc.Markups[1]}, Text: "foo"},
			{Type: TextMarker, ClosedMarkups: 1, Text: "foo"},
			{Type: TextMarker, OpenMarkups: []*Markup{&doc.Markups[1], &doc.Markups[2]}, ClosedMarkups: 1, Text: "foo"},
			{Type: TextMarker, ClosedMarkups: 1, Text: "foo"},
		}},
		{Type: MarkupSection, Tag: "p", Markers: []Marker{
			{Type: AtomMarker, Atom: &doc.Atoms[0]},
			{Type: AtomMarker, OpenMarkups: []*Markup{&doc.Markups[0]}, Atom: &doc.Atoms[1]},
			{Type: AtomMarker, ClosedMarkups: 1, Atom: &doc.Atoms[0]},
		}},
		{Type: ImageSection, Source: "https://example.com/foo.png"},
		{Type: ListSection, Tag: "ul", Items: [][]Marker{
			{
				{Type: TextMarker, ClosedMarkups: 0, Text: "foo"},
				{Type: TextMarker, OpenMarkups: []*Markup{&doc.Markups[0]}, ClosedMarkups: 1, Text: "foo"},
			},
			{
				{Type: TextMarker, OpenMarkups: []*Markup{&doc.Markups[0]}, Text: "foo"},
				{Type: TextMarker, ClosedMarkups: 1, Text: "<foo>"},
			},
		}},
		{Type: ListSection, Tag: "ol", Items: [][]Marker{
			{
				{Type: TextMarker, ClosedMarkups: 0, Text: "bar"},
				{Type: TextMarker, OpenMarkups: []*Markup{&doc.Markups[1]}, ClosedMarkups: 1, Text: "bar"},
			},
			{
				{Type: TextMarker, OpenMarkups: []*Markup{&doc.Markups[1]}, Text: "bar"},
				{Type: TextMarker, ClosedMarkups: 1, Text: "<bar>"},
			},
		}},
		{Type: 10, Card: &doc.Cards[1]},
	}
	return doc
}
