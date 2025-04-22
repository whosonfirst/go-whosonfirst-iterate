package iterator

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"sync/atomic"
	"testing"
)

func TestIterator(t *testing.T) {

	ctx := context.Background()

	count := int32(0)
	expected := int32(37)

	iter_cb := func(ctx context.Context, path string, r io.ReadSeeker, args ...interface{}) error {
		atomic.AddInt32(&count, 1)
		return nil
	}

	iter, err := NewIterator(ctx, "repo://", iter_cb)

	if err != nil {
		t.Fatalf("Failed to create new iterator, %v", err)
	}

	rel_path := "../fixtures"
	abs_path, err := filepath.Abs(rel_path)

	if err != nil {
		t.Fatalf("Failed to derive absolute path for %s, %v", rel_path, err)
	}

	err = iter.IterateURIs(ctx, abs_path)

	if err != nil {
		t.Fatalf("Failed to iterate %s, %v", abs_path, err)
	}

	if count != expected {
		t.Fatalf("Unexpected count: %d", count)
	}
}

func TestIteratorWithQuery(t *testing.T) {

	ctx := context.Background()

	count := int32(0)
	expected := int32(36)

	iter_cb := func(ctx context.Context, path string, r io.ReadSeeker, args ...interface{}) error {
		atomic.AddInt32(&count, 1)
		return nil
	}

	iter, err := NewIterator(ctx, "repo://?include=properties.sfomuseum:placetype=map&exclude=properties.sfomuseum:uri=2019", iter_cb)

	if err != nil {
		t.Fatalf("Failed to create new iterator, %v", err)
	}

	rel_path := "../fixtures"
	abs_path, err := filepath.Abs(rel_path)

	if err != nil {
		t.Fatalf("Failed to derive absolute path for %s, %v", rel_path, err)
	}

	err = iter.IterateURIs(ctx, abs_path)

	if err != nil {
		t.Fatalf("Failed to iterate %s, %v", abs_path, err)
	}

	if count != expected {
		t.Fatalf("Unexpected count: %d", count)
	}
}

func TestIteratorWithExcludePath(t *testing.T) {

	ctx := context.Background()

	count := int32(0)
	expected := int32(35)

	iter_cb := func(ctx context.Context, path string, r io.ReadSeeker, args ...interface{}) error {
		atomic.AddInt32(&count, 1)
		return nil
	}

	iter, err := NewIterator(ctx, "repo://?_exclude=1(5|7)1.*.geojson$", iter_cb)

	if err != nil {
		t.Fatalf("Failed to create new iterator, %v", err)
	}

	rel_path := "../fixtures"
	abs_path, err := filepath.Abs(rel_path)

	if err != nil {
		t.Fatalf("Failed to derive absolute path for %s, %v", rel_path, err)
	}

	err = iter.IterateURIs(ctx, abs_path)

	if err != nil {
		t.Fatalf("Failed to iterate %s, %v", abs_path, err)
	}

	if count != expected {
		t.Fatalf("Unexpected count: %d", count)
	}
}

func TestIteratorWithError(t *testing.T) {

	ctx := context.Background()

	iter_cb := func(ctx context.Context, path string, r io.ReadSeeker, args ...interface{}) error {
		return fmt.Errorf("Nope")
	}

	iter, err := NewIterator(ctx, "repo://?_exclude=1(5|7)1.*.geojson$", iter_cb)

	if err != nil {
		t.Fatalf("Failed to create new iterator, %v", err)
	}

	rel_path := "../fixtures"
	abs_path, err := filepath.Abs(rel_path)

	if err != nil {
		t.Fatalf("Failed to derive absolute path for %s, %v", rel_path, err)
	}

	err = iter.IterateURIs(ctx, abs_path)

	if err == nil {
		t.Fatalf("Expected an error iterating %s", abs_path)
	}

}

func TestIteratorExcludeAlt(t *testing.T) {

	ctx := context.Background()

	iter_cb := func(ctx context.Context, path string, r io.ReadSeeker, args ...interface{}) error {
		return fmt.Errorf("Nope")
	}

	iter, err := NewIterator(ctx, "repo://?_exclude_alt=true", iter_cb)

	if err != nil {
		t.Fatalf("Failed to create new iterator, %v", err)
	}

	rel_path := "../fixtures"
	abs_path, err := filepath.Abs(rel_path)

	if err != nil {
		t.Fatalf("Failed to derive absolute path for %s, %v", rel_path, err)
	}

	err = iter.IterateURIs(ctx, abs_path)

	if err == nil {
		t.Fatalf("Expected an error iterating %s", abs_path)
	}

}
