package iterate

import (
	"context"
	"io"
	"io/fs"
	"log/slog"
	"os"
	"testing"
)

func TestCwdIterator(t *testing.T) {

	slog.SetLogLoggerLevel(slog.LevelDebug)
	slog.Debug("Verbose logging enabled")

	expected := int64(0)

	cwd, err := os.Getwd()

	if err != nil {
		t.Fatalf("Failed to derive working directory, %v", err)
	}

	root, err := os.OpenRoot(cwd)

	if err != nil {
		t.Fatalf("Failed to open root, %v", err)
	}

	root_fs := root.FS()

	fs.WalkDir(root_fs, ".", func(path string, d fs.DirEntry, err error) error {

		if err != nil {
			t.Fatalf("Failed to walk root dir, %v", err)
		}

		if !d.IsDir() {
			expected += 1
		}

		return nil
	})

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

	if seen != expected {
		t.Fatalf("Unexpected record count. Got %d but expected %d", seen, expected)
	}

	slog.Info("Records seen", "count", seen)
}
