package emitter

import (
	"context"
	"io"
	"testing"
)

func TestNullEmitter(t *testing.T) {

	ctx := context.Background()

	e, err := NewEmitter(ctx, "null://")

	if err != nil {
		t.Fatalf("Failed to create null emitter, %v", err)
	}

	expected := int32(0)
	count := int32(0)

	cb := func(ctx context.Context, path string, r io.ReadSeeker, args ...interface{}) error {
		return nil
	}

	err = e.WalkURI(ctx, cb, "/dev/null")

	if err != nil {
		t.Fatalf("Failed to walk file, %v", err)
	}

	if count != expected {
		t.Fatalf("Unexpected count: %d", count)
	}
}
