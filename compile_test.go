package mobiledoc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompile(t *testing.T) {
	m, err := Compile(minimalDoc())
	assert.NoError(t, err)
	assert.Equal(t, minimalMap(), m)

	m, err = Compile(sampleDoc())
	assert.NoError(t, err)
	assert.Equal(t, sampleMap(), m)

	m, err = Compile(Document{Version: Version})
	assert.NoError(t, err)
	assert.Equal(t, Map{
		"version":  Version,
		"markups":  List{},
		"atoms":    List{},
		"cards":    List{},
		"sections": List{},
	}, m)
}

func BenchmarkCompile(b *testing.B) {
	in := sampleDoc()
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := Compile(in)
		if err != nil {
			panic(err)
		}
	}
}
