package iterator

import (
	"context"
	"sync/atomic"
	"testing"
)

func TestFileIterator(t *testing.T) {

	ctx := context.Background()

	it, err := NewIterator(ctx, "file://")

	if err != nil {
		t.Fatalf("Failed to create directory iterator, %v", err)
	}

	expected := int32(1)
	count := int32(0)

	for _, err := range it.Iterate(ctx, "../fixtures/data/151/183/838/5/1511838385.geojson") {

		if err != nil {
			t.Fatalf("Failed to walk file, %v", err)
		}

		atomic.AddInt32(&count, 1)
	}

	if count != expected {
		t.Fatalf("Unexpected count: %d", count)
	}
}
