package iterate

import (
	"context"
	"io"
	"log/slog"
	"path/filepath"
	"testing"
)

func TestRepoIterator(t *testing.T) {

	if *tests_verbose {
		slog.SetLogLoggerLevel(slog.LevelDebug)
		slog.Debug("Verbose logging enabled")
	}

	ctx := context.Background()

	abs_path, err := filepath.Abs("fixtures")

	if err != nil {
		t.Fatalf("Failed to derive absolute path for fixtures, %v", err)
	}

	it, err := NewIterator(ctx, "repo://")

	if err != nil {
		t.Fatalf("Failed to create new directory source, %v", err)
	}

	for rec, err := range it.Iterate(ctx, abs_path) {

		if err != nil {
			t.Fatalf("Failed to walk '%s', %v", abs_path, err)
			break
		}

		defer rec.Body.Close()
		_, err = io.ReadAll(rec.Body)

		if err != nil {
			t.Fatalf("Failed to read body for %s, %v", rec.Path, err)
		}

		_, err = rec.Body.Seek(0, 0)

		if err != nil {
			t.Fatalf("Failed to rewind body for %s, %v", rec.Path, err)
		}

		_, err = io.ReadAll(rec.Body)

		if err != nil {
			t.Fatalf("Failed second read body for %s, %v", rec.Path, err)
		}
	}

	seen := it.Seen()
	expected := int64(37)

	if seen != expected {
		t.Fatalf("Unexpected record count. Got %d but expected %d", seen, expected)
	}

	if err != nil {
		t.Fatalf("Failed to close iterator")
	}
}
