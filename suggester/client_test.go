package suggester

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/dnaeon/go-vcr/recorder"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient(t *testing.T) {
	ctx := context.Background()
	host := "localhost"
	port := 8983

	suggestEndpoint := "suggest"

	t.Run("suggest ok", func(t *testing.T) {
		a := assert.New(t)

		rec, err := recorder.New("fixtures/suggest-ok")
		require.NoError(t, err)
		defer rec.Stop()

		client := NewCustomClient(host, port, suggestEndpoint,
			&http.Client{
				Timeout:   time.Second * 60,
				Transport: rec,
			},
		)

		resp, err := client.Suggest(ctx, Request{
			Collection: "techproducts",
			Params: Params{
				Query:        "elec",
				Count:        10,
				Dictionaries: []string{"mySuggester"},
				Build:        true,
			},
		})
		a.NoError(err)
		a.NotNil(resp)
	})

	t.Run("suggest error", func(t *testing.T) {
		t.Run("endpoint not found", func(t *testing.T) {
			a := assert.New(t)

			rec, err := recorder.New("fixtures/suggest-error-404")
			require.NoError(t, err)
			defer rec.Stop()

			client := NewCustomClient(host, port, "/not-exists",
				&http.Client{
					Timeout:   time.Second * 60,
					Transport: rec,
				},
			)

			_, err = client.Suggest(ctx, Request{
				Collection: "techproducts",
				Params: Params{
					Query:        "elec",
					Count:        10,
					Dictionaries: []string{"mySuggester"},
					Build:        true,
				},
			})
			a.Error(err)
		})

		t.Run("parse url error", func(t *testing.T) {
			a := assert.New(t)

			client := NewCustomClient("http\\\\\\::whatever:://\\::", 1234, "/not-exists",
				&http.Client{
					Timeout: time.Second * 60,
				},
			)

			_, err := client.Suggest(ctx, Request{
				Collection: "techproducts///\\4343::343",
				Params: Params{
					Query:        "elec??&&",
					Count:        10,
					Dictionaries: []string{"mySuggester"},
					Build:        true,
				},
			})
			a.Error(err)
		})

		t.Run("other errors", func(t *testing.T) {
			a := assert.New(t)

			rec, err := recorder.New("fixtures/suggest-error")
			require.NoError(t, err)
			defer rec.Stop()

			client := NewCustomClient(host, port, suggestEndpoint,
				&http.Client{
					Timeout:   time.Second * 60,
					Transport: rec,
				},
			)

			_, err = client.Suggest(ctx, Request{})
			a.Error(err)

			_, err = client.Suggest(ctx, Request{
				Collection: "techproducts",
				Params:     Params{},
			})
			a.Error(err)

			_, err = client.Suggest(ctx, Request{
				Collection: "techproducts",
				Params: Params{
					Query: "elec",
				},
			})
			a.Error(err)
		})
	})
}

func Test_buildParams(t *testing.T) {
	params := buildParams(Params{
		Query:        "elec",
		Dictionaries: []string{"mySuggester"},
		Count:        10,
		Cfq:          "memory",
		Build:        true,
		Reload:       true,
		BuildAll:     true,
		ReloadAll:    true,
	})

	assert.Len(t, params, 9)
}
