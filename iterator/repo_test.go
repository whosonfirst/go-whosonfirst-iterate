package iterator

import (
	"context"
	"sync/atomic"
	"testing"
)

func TestRepoIterator(t *testing.T) {

	ctx := context.Background()

	e, err := NewIterator(ctx, "repo://")

	if err != nil {
		t.Fatalf("Failed to create repo iterator, %v", err)
	}

	expected := int32(37)
	count := int32(0)

	for _, err := range e.Iterate(ctx, "../fixtures") {

		if err != nil {
			t.Fatalf("Failed to walk repo, %v", err)
		}

		atomic.AddInt32(&count, 1)
	}

	if count != expected {
		t.Fatalf("Unexpected count: %d", count)
	}
}
