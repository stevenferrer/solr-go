package index_test

import (
	"testing"

	"github.com/sf9v/solr-go/index"
	"github.com/stretchr/testify/assert"
)

func TestDocs(t *testing.T) {
	var names = []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}{
		{
			ID:   "1",
			Name: "Milana Vino",
		},
		{
			ID:   "2",
			Name: "Charly Jordan",
		},
		{
			ID:   "3",
			Name: "Daisy Keech",
		},
	}

	docs := index.NewDocs()
	for _, model := range names {
		docs.AddDoc(model)
	}

	assert.Equal(t, docs.Count(), 3)

	b, err := docs.Marshal()
	assert.NoError(t, err)

	expected := `[{"id":"1","name":"Milana Vino"},{"id":"2","name":"Charly Jordan"},{"id":"3","name":"Daisy Keech"}]`
	assert.Equal(t, expected, string(b))
}
