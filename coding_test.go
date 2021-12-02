package mobiledoc

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
)

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

func TestBSON(t *testing.T) {
	in, err := bson.Marshal(sampleMap())
	require.NoError(t, err)

	var doc Document
	err = bson.Unmarshal(in, &doc)
	require.NoError(t, err)

	err = formatValidator.Validate(doc)
	require.NoError(t, err)

	out, err := bson.Marshal(doc)
	require.NoError(t, err)

	var res Map
	err = bson.Unmarshal(out, &res)
	require.NoError(t, err)
	equalMaps(t, res, sampleMap())
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
