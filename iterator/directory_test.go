package iterator

import (
	"context"
	"log/slog"
	"sync/atomic"
	"testing"
)

func TestDirectoryEmitter(t *testing.T) {

	slog.SetLogLoggerLevel(slog.LevelDebug)

	ctx := context.Background()

	it, err := NewIterator(ctx, "directory://")

	if err != nil {
		t.Fatalf("Failed to create directory emitter, %v", err)
	}

	expected := int32(37)
	count := int32(0)

	for _, err := range it.Iterate(ctx, "../fixtures/data") {

		if err != nil {
			t.Fatalf("Failed to walk directory, %v", err)
		}

		atomic.AddInt32(&count, 1)
	}

	if count != expected {
		t.Fatalf("Unexpected count: %d", count)
	}
}
