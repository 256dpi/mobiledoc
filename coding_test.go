package mobiledoc

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
)

type docStruct struct {
	Doc *Document `json:"doc" bson:"doc"`
}

func TestJSON(t *testing.T) {
	in, err := json.Marshal(sampleMap())
	require.NoError(t, err)

	var doc Document
	err = json.Unmarshal(in, &doc)
	require.NoError(t, err)

	err = formatValidator.Validate(doc)
	require.NoError(t, err)

	out, err := json.Marshal(doc)
	require.NoError(t, err)

	var res Map
	err = json.Unmarshal(out, &res)
	require.NoError(t, err)
	equalMaps(t, res, sampleMap())
}

func TestJSONNil(t *testing.T) {
	in, err := json.Marshal(docStruct{})
	require.NoError(t, err)

	var doc docStruct
	err = json.Unmarshal(in, &doc)
	require.NoError(t, err)
	assert.Nil(t, doc.Doc)

	out, err := json.Marshal(doc)
	require.NoError(t, err)

	var res Map
	err = json.Unmarshal(out, &res)
	require.NoError(t, err)
	equalMaps(t, res, Map{
		"doc": nil,
	})
}

func TestBSON(t *testing.T) {
	in, err := bson.Marshal(Map{"doc": sampleMap()})
	require.NoError(t, err)

	var doc docStruct
	err = bson.Unmarshal(in, &doc)
	require.NoError(t, err)

	err = formatValidator.Validate(*doc.Doc)
	require.NoError(t, err)

	out, err := bson.Marshal(doc)
	require.NoError(t, err)

	var res Map
	err = bson.Unmarshal(out, &res)
	require.NoError(t, err)
	equalMaps(t, res, Map{"doc": sampleMap()})
}

func TestBSONNil(t *testing.T) {
	in, err := bson.Marshal(docStruct{})
	require.NoError(t, err)

	var doc docStruct
	err = bson.Unmarshal(in, &doc)
	require.NoError(t, err)
	// assert.Nil(t, doc.Doc) // TODO: Should be nil.

	out, err := bson.Marshal(doc)
	require.NoError(t, err)

	var res Map
	err = bson.Unmarshal(out, &res)
	require.NoError(t, err)
	equalMaps(t, res, Map{
		"doc": Map{},
	})
}

func equalMaps(t *testing.T, m1, m2 Map) {
	o1, e1 := json.Marshal(m1)
	o2, e2 := json.Marshal(m2)
	assert.NoError(t, e1)
	assert.NoError(t, e2)
	assert.Equal(t, o1, o2)
}

func BenchmarkJSON(b *testing.B) {
	in := sampleDoc()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		bytes, err := json.Marshal(in)
		if err != nil {
			panic(err)
		}

		var out Document
		err = json.Unmarshal(bytes, &out)
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkBSON(b *testing.B) {
	in := sampleDoc()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		bytes, err := bson.Marshal(in)
		if err != nil {
			panic(err)
		}

		var out Document
		err = bson.Unmarshal(bytes, &out)
		if err != nil {
			panic(err)
		}
	}
}
