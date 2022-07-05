package emitter

import (
	"context"
	"io"
	"sync/atomic"
	"testing"
)

func TestRepoEmitter(t *testing.T) {

	ctx := context.Background()

	e, err := NewEmitter(ctx, "repo://")

	if err != nil {
		t.Fatalf("Failed to create repo emitter, %v", err)
	}

	expected := int32(37)
	count := int32(0)

	cb := func(ctx context.Context, path string, r io.ReadSeeker, args ...interface{}) error {
		atomic.AddInt32(&count, 1)
		return nil
	}

	err = e.WalkURI(ctx, cb, "../fixtures")

	if err != nil {
		t.Fatalf("Failed to walk repo, %v", err)
	}

	if count != expected {
		t.Fatalf("Unexpected count: %d", count)
	}
}
