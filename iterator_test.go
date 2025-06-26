package iterate

import (
	"context"
	"log/slog"
	"testing"
)

func TestNewIterator(t *testing.T) {

	if *tests_verbose {
		slog.SetLogLoggerLevel(slog.LevelDebug)
		slog.Debug("Verbose logging enabled")
	}

	ctx := context.Background()

	for _, s := range IteratorSchemes() {

		it, err := NewIterator(ctx, s)

		if err != nil {
			t.Fatalf("Failed to create new iterator for '%s', %v", s, err)
		}

		if it.Seen() != 0 {
			t.Fatalf("Unexpected seen count for '%s', %d", s, it.Seen())
		}

		if it.IsIterating() {
			t.Fatalf("Why is '%s' iterating?", s)
		}

		err = it.Close()

		if err != nil {
			t.Fatalf("Failed to close iterator %s, %v", s, err)
		}
	}
}
