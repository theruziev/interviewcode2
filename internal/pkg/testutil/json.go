package testutil

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMarshalJSON(t *testing.T, data any) []byte {
	t.Helper()
	body, err := json.Marshal(data)
	require.NoError(t, err)
	return body
}
