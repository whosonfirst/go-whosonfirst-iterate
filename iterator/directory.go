package iterator

import (
	"context"
	"fmt"
	"iter"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/whosonfirst/go-whosonfirst-iterate/v3/filters"
)

func init() {
	ctx := context.Background()
	err := RegisterIterator(ctx, "directory", NewDirectoryIterator)

	if err != nil {
		panic(err)
	}
}

// DirectoryIterator implements the `Iterator` interface for crawling records in a directory.
type DirectoryIterator struct {
	Iterator
	// filters is a `filters.Filters` instance used to include or exclude specific records from being crawled.
	filters filters.Filters
}

// NewDirectoryIterator() returns a new `DirectoryIterator` instance configured by 'uri' in the form of:
//
//	directory://?{PARAMETERS}
//
// Where {PARAMETERS} may be:
// * `?include=` Zero or more `aaronland/go-json-query` query strings containing rules that must match for a document to be considered for further processing.
// * `?exclude=` Zero or more `aaronland/go-json-query`	query strings containing rules that if matched will prevent a document from being considered for further processing.
// * `?include_mode=` A valid `aaronland/go-json-query` query mode string for testing inclusion rules.
// * `?exclude_mode=` A valid `aaronland/go-json-query` query mode string for testing exclusion rules.
func NewDirectoryIterator(ctx context.Context, uri string) (Iterator, error) {

	f, err := filters.NewQueryFiltersFromURI(ctx, uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to create filters from query, %w", err)
	}

	idx := &DirectoryIterator{
		filters: f,
	}

	return idx, nil
}

func (idx *DirectoryIterator) Walk(ctx context.Context, uris ...string) iter.Seq2[*Candidate, error] {

	return func(yield func(*Candidate, error) bool) {

		for _, uri := range uris {

			logger := slog.Default()
			logger = logger.With("uri", uri)

			abs_path, err := filepath.Abs(uri)

			if err != nil {
				logger.Error("Failed to derive absolute path", "error", err)
				yield(nil, err)
				continue
			}

			logger = logger.With("path", abs_path)

			fs_opts := &FSIteratorOptions{
				Filters: idx.filters,
				FS:      os.DirFS(abs_path),
			}

			fs_iter, err := NewFSIteratorWithOptions(ctx, fs_opts)

			if err != nil {
				logger.Error("Failed to create new FS iterator", "error", err)
				yield(nil, err)
				continue
			}

			for c, err := range fs_iter.Walk(ctx) {
				yield(c, err)
			}
		}
	}
}
