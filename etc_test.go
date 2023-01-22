package solr

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_wrapErr(t *testing.T) {
	tests := []struct {
		name       string
		msg        string
		wrappedErr error
	}{
		{
			name: "nil error",
			msg:  "an error",
		},
		{
			name:       "non-nil error",
			msg:        "an error",
			wrappedErr: assert.AnError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := wrapErr(tc.wrappedErr, tc.msg)
			if tc.wrappedErr != nil {
				assert.ErrorIs(t, err, tc.wrappedErr)
				expectStr := fmt.Sprintf("%s: %s", tc.msg, tc.wrappedErr.Error())
				assert.Equal(t, expectStr, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
