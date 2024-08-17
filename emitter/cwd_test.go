package emitter

import (
	"context"
	"io"
	"sync/atomic"
	"testing"
)

func TestCwdEmitter(t *testing.T) {

	ctx := context.Background()

	e, err := NewEmitter(ctx, "cwd://")

	if err != nil {
		t.Fatalf("Failed to create directory emitter, %v", err)
	}

	expected := int32(19)
	count := int32(0)

	cb := func(ctx context.Context, path string, r io.ReadSeeker, args ...interface{}) error {
		atomic.AddInt32(&count, 1)
		return nil
	}

	err = e.WalkURI(ctx, cb, ".")

	if err != nil {
		t.Fatalf("Failed to walk directory, %v", err)
	}

	if count != expected {
		t.Fatalf("Unexpected count: %d", count)
	}
}
