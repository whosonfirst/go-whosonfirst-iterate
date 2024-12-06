package iterator

import (
	"context"
	"log/slog"
	"sync/atomic"
	"testing"
)

func TestGeojsonLIterator(t *testing.T) {

	slog.SetLogLoggerLevel(slog.LevelDebug)

	ctx := context.Background()

	it, err := NewIterator(ctx, "geojsonl://")

	if err != nil {
		t.Fatalf("Failed to create directory iterator, %v", err)
	}

	expected := int32(2)
	count := int32(0)

	for _, err := range it.Iterate(ctx, "../fixtures/collection.geojsonl") {

		if err != nil {
			t.Fatalf("Failed to walk list, %v", err)
		}

		atomic.AddInt32(&count, 1)
	}

	if count != expected {
		t.Fatalf("Unexpected count: %d", count)
	}
}
