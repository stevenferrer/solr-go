package index_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	. "github.com/stevenferrer/solr-go"
	"github.com/stevenferrer/solr-go/index"
)

func TestCommands(t *testing.T) {
	t.Run("add command", func(t *testing.T) {
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

		cmdStrs := []string{}
		for _, cmd := range addCmds {
			cmdStr, err := cmd.Command()
			require.NoError(t, err)
			cmdStrs = append(cmdStrs, cmdStr)
		}

		got := strings.Join(cmdStrs, ",")
		expected := `"add":{"commitWithin":5000,"doc":{"id":"1","name":"Milana Vino"},"overwrite":true},"add":{"doc":{"id":"2","name":"Daisy Keech"}},"add":{"doc":{"id":"3","name":"Charley Jordan"}}`

		assert.Equal(t, expected, got)
	})

	t.Run("delete by query command", func(t *testing.T) {
		delByQryCmds := []index.DelByQryCommand{
			{
				Query: "*:*",
			},
			{
				Query: "delete me",
			},
		}
		cmdStrs := []string{}
		for _, cmd := range delByQryCmds {
			cmdStr, err := cmd.Command()
			require.NoError(t, err)
			cmdStrs = append(cmdStrs, cmdStr)
		}

		got := strings.Join(cmdStrs, ",")
		expected := `"delete":{"query":"*:*"},"delete":{"query":"delete me"}`

		assert.Equal(t, expected, got)
	})

	t.Run("delete by ids command", func(t *testing.T) {
		delByQryCmds := []index.DelByIDsCommand{
			{
				IDs: []string{"ID1", "ID2"},
			},
			{
				IDs: []string{"ID3", "ID4"},
			},
		}
		cmdStrs := []string{}
		for _, cmd := range delByQryCmds {
			cmdStr, err := cmd.Command()
			require.NoError(t, err)
			cmdStrs = append(cmdStrs, cmdStr)
		}

		got := strings.Join(cmdStrs, ",")
		expected := `"delete":["ID1","ID2"],"delete":["ID3","ID4"]`

		assert.Equal(t, expected, got)
	})
}
