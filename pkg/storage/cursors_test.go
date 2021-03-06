package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCursorString(t *testing.T) {
	cursor := itemsCursor{
		TS:    "c0ac559a-1296-4610-92a2-00b422fc26f4",
		Limit: 12,
	}

	expected := "eyJ1dWlkIjoiYzBhYzU1OWEtMTI5Ni00NjEwLTkyYTItMDBiNDIyZmMyNmY0IiwibGltaXQiOjEyfQ"

	actual := cursor.String()

	assert.Equal(t, expected, actual)
}

func TestCursorParse(t *testing.T) {
	c := "eyJ1dWlkIjoiYzBhYzU1OWEtMTI5Ni00NjEwLTkyYTItMDBiNDIyZmMyNmY0IiwibGltaXQiOjEyfQ"

	expected := itemsCursor{
		TS:    "c0ac559a-1296-4610-92a2-00b422fc26f4",
		Limit: 12,
	}

	var actual itemsCursor
	err := actual.Parse(c)
	require.NoError(t, err)

	assert.Equal(t, expected, actual)

}
