package index_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/dnaeon/go-vcr/recorder"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	. "github.com/stevenferrer/helios"
	"github.com/stevenferrer/helios/index"
	"github.com/stevenferrer/helios/schema"
)

func TestJSONClient(t *testing.T) {
	ctx := context.Background()
	collection := "gettingstarted"
	host := "localhost"
	port := 8983
	timeout := time.Second * 60

	r, err := recorder.New("fixtures/add-name-field")
	assert.NoError(t, err)

	schemaClient := schema.NewClient(host, port, &http.Client{
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

	t.Run("add single doc", func(t *testing.T) {
		t.Run("ok", func(t *testing.T) {
			a := assert.New(t)

			rec, err := recorder.New("fixtures/add-single-doc-ok")
			require.NoError(t, err)
			defer rec.Stop()

			client := index.NewJSONClient(host, port, &http.Client{
				Timeout:   timeout,
				Transport: rec,
			})

			var doc = struct {
				ID   string `json:"id"`
				Name string `json:"name"`
			}{
				ID:   "1",
				Name: "Milana Vino",
			}

			err = client.AddSingle(ctx, collection, doc)
			a.NoError(err)
		})
		t.Run("error", func(t *testing.T) {
			a := assert.New(t)

			rec, err := recorder.New("fixtures/add-single-doc-error")
			require.NoError(t, err)
			defer rec.Stop()

			client := index.NewJSONClient(host, port, &http.Client{
				Timeout:   timeout,
				Transport: rec,
			})

			// empty document
			var doc = struct {
				ID   string `json:"id"`
				Name string `json:"name"`
			}{}

			err = client.AddSingle(ctx, collection, doc)
			a.Error(err)
		})
	})

	t.Run("add multiple docs", func(t *testing.T) {
		a := assert.New(t)

		rec, err := recorder.New("fixtures/add-multiple-docs")
		require.NoError(t, err)
		defer rec.Stop()

		client := index.NewJSONClient(host, port, &http.Client{
			Timeout:   timeout,
			Transport: rec,
		})

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
				Name: "Charlie Jordan",
			},
			{
				ID:   "3",
				Name: "Daisy Keech",
			},
		}

		err = client.AddMultiple(ctx, collection, docs)
		a.NoError(err)
	})

	t.Run("multiple update commands", func(t *testing.T) {
		a := assert.New(t)

		rec, err := recorder.New("fixtures/multiple-update-commands")
		require.NoError(t, err)
		defer rec.Stop()

		client := index.NewJSONClient(host, port, &http.Client{
			Timeout:   timeout,
			Transport: rec,
		})

		addCmds := []index.AddCommand{
			{
				CommitWithin: 5000,
				Overwrite:    true,
				Doc: M{
					"id":   "1",
					"name": "Milana Vino",
				},
			},
			{
				Doc: M{
					"id":   "2",
					"name": "Daisy Keech",
				},
			},

			{
				Doc: M{
					"id":   "3",
					"name": "Charley Jordan",
				},
			},
		}

		delByIDsCmd := []index.DelByIDsCommand{
			{
				IDs: []string{"2"},
			},
		}

		delByQryCmd := []index.DelByQryCommand{
			{
				Query: "*:*",
			},
		}

		cmds := []index.Commander{}
		for _, ac := range addCmds {
			cmds = append(cmds, ac)
		}

		for _, dc := range delByIDsCmd {
			cmds = append(cmds, dc)
		}

		for _, dc := range delByQryCmd {
			cmds = append(cmds, dc)
		}

		err = client.UpdateCommands(ctx, collection, cmds...)
		a.NoError(err)
	})
}
