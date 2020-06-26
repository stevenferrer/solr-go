package config_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	solrconfig "github.com/stevenferrer/solr-go/config"
)

func TestCommands(t *testing.T) {
	t.Run("common properties", func(t *testing.T) {
		t.Run("set property", func(t *testing.T) {
			setPropCommand := solrconfig.NewSetPropCommand("updateHandler.autoCommit.maxTime", 15000)

			got, err := setPropCommand.Command()
			require.NoError(t, err)

			expected := `"set-property": {"updateHandler.autoCommit.maxTime":15000}`
			assert.Equal(t, expected, got)
		})

		t.Run("unset property", func(t *testing.T) {
			unsetPropCommand := solrconfig.NewUnsetPropCommand("updateHandler.autoCommit.maxTime")

			got, err := unsetPropCommand.Command()
			require.NoError(t, err)

			expected := `"unset-property": "updateHandler.autoCommit.maxTime"`
			assert.Equal(t, expected, got)
		})
	})
}
