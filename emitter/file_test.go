package emitter

import (
	"context"
	"io"
	"sync/atomic"
	"testing"
)

func TestFileEmitter(t *testing.T) {

	ctx := context.Background()

	e, err := NewEmitter(ctx, "file://")

	if err != nil {
		t.Fatalf("Failed to create directory emitter, %v", err)
	}

	expected := int32(1)
	count := int32(0)

	cb := func(ctx context.Context, path string, r io.ReadSeeker, args ...interface{}) error {
		atomic.AddInt32(&count, 1)
		return nil
	}

	err = e.WalkURI(ctx, cb, "../fixtures/data/151/183/838/5/1511838385.geojson")

	if err != nil {
		t.Fatalf("Failed to walk file, %v", err)
	}

	if count != expected {
		t.Fatalf("Unexpected count: %d", count)
	}
}
