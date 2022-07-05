package emitter

import (
	"context"
	"testing"
	"io"
	"sync/atomic"
)

func TestDirectoryEmitter(t *testing.T) {

	ctx := context.Background()

	e, err := NewEmitter(ctx, "directory://")

	if err != nil {
		t.Fatalf("Failed to create directory emitter, %v", err)
	}

	expected := int32(37)
	count := int32(0)
	
	cb := func(ctx context.Context, path string, r io.ReadSeeker, args ...interface{}) error {
		atomic.AddInt32(&count, 1)
		return nil
	}
	
	err = e.WalkURI(ctx, cb, "../fixtures/data")

	if err != nil {
		t.Fatalf("Failed to walk directory, %v", err)
	}

	if count != expected {
		t.Fatalf("Unexpected count: %d", count)
	}
}
