package iterate

import (
	"context"
	"fmt"
	"path/filepath"
	"testing"
)

func TestDirectoryIterator(t *testing.T) {

	ctx := context.Background()

	abs_path, err := filepath.Abs("fixtures")

	if err != nil {
		t.Fatalf("Failed to derive absolute path for fixtures, %v", err)
	}

	it, err := NewIterator(ctx, "directory://")

	if err != nil {
		t.Fatalf("Failed to create new directory source, %v", err)
	}

	for rec, err := range it.Iterate(ctx, abs_path) {

		if err != nil {
			//t.Fatalf("Failed to walk '%s', %v", abs_path, err)
			//break
		}

		fmt.Println(rec.Path)
	}
}
