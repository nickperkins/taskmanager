package idgen

import (
	"testing"
)

func TestGenerateTaskID_UniqueAndValid(t *testing.T) {
	ids := make(map[string]struct{})
	for i := 0; i < 1000; i++ {
		id := GenerateTaskID()
		if _, exists := ids[id]; exists {
			t.Fatalf("duplicate id generated: %s", id)
		}
		ids[id] = struct{}{}
		if len(id) != 36 {
			t.Errorf("id length is not 36: %s", id)
		}
	}
}
