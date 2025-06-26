package iterate

import (
	"context"
	"log/slog"
	"path/filepath"
	"testing"
)

func TestNullIterator(t *testing.T) {

	if *tests_verbose {
		slog.SetLogLoggerLevel(slog.LevelDebug)
		slog.Debug("Verbose logging enabled")
	}

	ctx := context.Background()

	abs_path, err := filepath.Abs("fixtures")

	if err != nil {
		t.Fatalf("Failed to derive absolute path for fixtures, %v", err)
	}

	it, err := NewIterator(ctx, "null://")

	if err != nil {
		t.Fatalf("Failed to create new null source, %v", err)
	}

	defer it.Close()

	for _, err := range it.Iterate(ctx, abs_path) {

		if err != nil {
			t.Fatalf("Failed to walk '%s', %v", abs_path, err)
			break
		}

		// slog.Debug("Record", "source", abs_path, "path", rec.Path)
	}

	seen := it.Seen()
	expected := int64(0)

	if seen != expected {
		t.Fatalf("Unexpected record count. Got %d but expected %d", seen, expected)
	}

	slog.Info("Records seen", "count", seen)
}
