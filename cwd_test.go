package iterate

import (
	"context"
	"io"
	"log/slog"
	"testing"
)

func TestCwdIterator(t *testing.T) {

	slog.SetLogLoggerLevel(slog.LevelDebug)
	slog.Debug("Verbose logging enabled")

	ctx := context.Background()

	it, err := NewIterator(ctx, "cwd://")

	if err != nil {
		t.Fatalf("Failed to create new directory source, %v", err)
	}

	for rec, err := range it.Iterate(ctx, ".") {

		if err != nil {
			t.Fatalf("Failed to walk working directory, %v", err)
			break
		}

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
	expected := int64(311)

	if seen != expected {
		t.Fatalf("Unexpected record count. Got %d but expected %d", seen, expected)
	}

	slog.Info("Records seen", "count", seen)
}
