package iterator

import (
	"context"
	"sync/atomic"
	"testing"
)

func TestCwdIterator(t *testing.T) {

	ctx := context.Background()

	e, err := NewIterator(ctx, "cwd://")

	if err != nil {
		t.Fatalf("Failed to create directory iterator, %v", err)
	}

	expected := int32(16)
	count := int32(0)

	for _, err := range e.Iterate(ctx, "cwd://") {

		if err != nil {
			t.Fatal(err)
		}

		atomic.AddInt32(&count, 1)
	}

	if count != expected {
		t.Fatalf("Unexpected count: %d", count)
	}
}
