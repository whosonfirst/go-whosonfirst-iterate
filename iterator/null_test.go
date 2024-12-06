package iterator

import (
	"context"
	"testing"
)

func TestNullIterator(t *testing.T) {

	ctx := context.Background()

	e, err := NewIterator(ctx, "null://")

	if err != nil {
		t.Fatalf("Failed to create null iterator, %v", err)
	}

	expected := int32(0)
	count := int32(0)

	for _, err := range e.Iterate(ctx, "null://", "/dev/null") {

		if err != nil {
			t.Fatalf("Failed to walk file, %v", err)
		}
	}

	if count != expected {
		t.Fatalf("Unexpected count: %d", count)
	}
}
