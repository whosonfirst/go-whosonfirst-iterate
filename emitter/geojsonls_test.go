package emitter

import (
	"context"
	"io"
	"sync/atomic"
	"testing"
)

func TestGeojsonLEmitter(t *testing.T) {

	ctx := context.Background()

	e, err := NewEmitter(ctx, "geojsonl://")

	if err != nil {
		t.Fatalf("Failed to create directory emitter, %v", err)
	}

	expected := int32(2)
	count := int32(0)

	cb := func(ctx context.Context, path string, r io.ReadSeeker, args ...interface{}) error {
		atomic.AddInt32(&count, 1)
		return nil
	}

	err = e.WalkURI(ctx, cb, "../fixtures/collection.geojsonl")

	if err != nil {
		t.Fatalf("Failed to walk list, %v", err)
	}

	if count != expected {
		t.Fatalf("Unexpected count: %d", count)
	}
}
