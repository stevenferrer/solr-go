package index_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/dnaeon/go-vcr/recorder"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/stevenferrer/solr-go/index"
	"github.com/stevenferrer/solr-go/schema"
)

type M = map[string]interface{}

func TestJSONClient(t *testing.T) {
	ctx := context.Background()
	collection := "gettingstarted"
	host := "localhost"
	port := 8983
	timeout := time.Second * 60

	r, err := recorder.New("fixtures/init-schema")
	assert.NoError(t, err)

	// only for covering
	_ = index.NewJSONClient(host, port)

	schemaClient := schema.NewClientWithHTTPClient(host, port, &http.Client{
		Timeout:   timeout,
		Transport: r,
	})
	err = schemaClient.AddField(ctx, collection, schema.Field{
		Name:    "name",
		Type:    "text_general",
		Indexed: true,
		Stored:  true,
	})
	require.NoError(t, err)

	// add copy field
	err = schemaClient.AddCopyField(ctx, collection, schema.CopyField{
		Source: "*",
		Dest:   "_text_",
	})
	require.NoError(t, err)
	require.NoError(t, r.Stop())

	t.Run("add docs", func(t *testing.T) {
		t.Run("ok", func(t *testing.T) {
			a := assert.New(t)

			rec, err := recorder.New("fixtures/add-docs")
			require.NoError(t, err)
			defer rec.Stop()

			client := index.NewJSONClientWithHTTPClient(host, port, &http.Client{
				Timeout:   timeout,
				Transport: rec,
			})

			err = client.UpdateCommands(ctx, collection)
			a.NoError(err)

			var docs = []struct {
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

			err = client.AddDocs(ctx, collection, docs)
			a.NoError(err)
			err = client.Commit(ctx, collection)
			a.NoError(err)
		})

		t.Run("error", func(t *testing.T) {
			t.Run("parse url error", func(t *testing.T) {
				client := index.NewJSONClientWithHTTPClient("wtf:\\:\\", port, &http.Client{})
				err := client.AddDocs(ctx, "wtf:\\//\\::", nil)
				assert.Error(t, err)
			})
		})
	})

	t.Run("multiple update commands", func(t *testing.T) {
		t.Run("ok", func(t *testing.T) {
			a := assert.New(t)

			rec, err := recorder.New("fixtures/update-commands")
			require.NoError(t, err)
			defer rec.Stop()

			client := index.NewJSONClientWithHTTPClient(host, port, &http.Client{
				Timeout:   timeout,
				Transport: rec,
			})

			err = client.UpdateCommands(ctx, collection,
				index.AddCommand{
					CommitWithin: 5000,
					Overwrite:    true,
					Doc: M{
						"id":   "1",
						"name": "Milana Vino",
					},
				},
				index.AddCommand{
					Doc: M{
						"id":   "2",
						"name": "Daisy Keech",
					},
				},

				index.AddCommand{
					Doc: M{
						"id":   "3",
						"name": "Charley Jordan",
					},
				},
				index.DelByIDsCommand{
					IDs: []string{"2"},
				},
				index.DelByQryCommand{
					Query: "*:*",
				},
			)
			a.NoError(err)
			err = client.Commit(ctx, collection)
			a.NoError(err)
		})

		t.Run("error", func(t *testing.T) {
			t.Run("parse url error", func(t *testing.T) {
				client := index.NewJSONClientWithHTTPClient("wtf:\\:\\", port, &http.Client{})
				err := client.UpdateCommands(ctx, "wtf:\\//\\::", index.AddCommand{})
				assert.Error(t, err)
			})
		})
	})
}
