package idgen

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGenerateTaskID_UniqueAndValid(t *testing.T) {
	ids := make(map[string]struct{})
	for i := 0; i < 1000; i++ {
		id := GenerateTaskID()
		require.NotContains(t, ids, id, "duplicate id generated: %s", id)
		ids[id] = struct{}{}
		assert.Equal(t, 36, len(id), "id length is not 36: %s", id)
	}
}
