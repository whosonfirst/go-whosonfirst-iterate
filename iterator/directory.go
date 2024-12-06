package iterator

import (
	"context"
	"fmt"
	"io/fs"
	"iter"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/whosonfirst/go-ioutil"
	"github.com/whosonfirst/go-whosonfirst-iterate/v3/filters"
)

func init() {
	ctx := context.Background()
	RegisterIterator(ctx, "directory", NewDirectoryIterator)
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

	it := &DirectoryIterator{
		filters: f,
	}

	return it, nil
}

// WalkURI() walks (crawls) the directory named 'uri' and for each file (not excluded by any filters specified
// when `it` was created) invokes 'index_cb'.
func (it *DirectoryIterator) Iterate(ctx context.Context, uris ...string) iter.Seq2[Record, error] {

	return func(yield func(Record, error) bool) {

		for _, uri := range uris {

			for r, err := range it.iterate(ctx, uri) {
				yield(r, err)
			}
		}
	}
}

func (it *DirectoryIterator) iterate(ctx context.Context, uri string) iter.Seq2[Record, error] {

	logger := slog.Default()
	logger = logger.With("uri", uri)

	return func(yield func(Record, error) bool) {

		abs_path, err := filepath.Abs(uri)

		if err != nil {
			logger.Debug("Failed to derive absolute path", "error", err)
			yield(nil, fmt.Errorf("Failed to derive absolute path for '%s', %w", uri, err))
			return
		}

		dir_fs := os.DirFS(abs_path)

		err = fs.WalkDir(dir_fs, ".", func(path string, d fs.DirEntry, err error) error {

			logger := slog.Default()
			logger = logger.With("uri", uri)
			logger = logger.With("path", path)

			if err != nil {
				logger.Debug("WalkDir reported an error", "error", err)
				return fmt.Errorf("Walk error, %w", err)
			}

			if d.IsDir() {
				return nil
			}

			r, err := dir_fs.Open(path)

			if err != nil {
				logger.Debug("Failed to open path for reading", "error", err)
				return fmt.Errorf("Failed to open %s for reading, %w", path, err)
			}

			rsc, err := ioutil.NewReadSeekCloser(r)

			if err != nil {
				logger.Debug("Failed to create ReadSeekCloser", "error", err)
				return fmt.Errorf("Failed to create ReadSeekCloser for %s, %w", path, err)
			}

			defer rsc.Close()

			if it.filters != nil {

				ok, err := it.filters.Apply(ctx, rsc)

				if err != nil {
					logger.Debug("Failed to apply filters", "error", err)
					return fmt.Errorf("Failed to apply filters for '%s', %w", path, err)
				}

				if !ok {
					logger.Debug("No matches after applying filters, skipping")
					return nil
				}

				_, err = rsc.Seek(0, 0)

				if err != nil {
					logger.Debug("Failed to rewind reader", "error", err)
					return fmt.Errorf("Failed to seek(0, 0) on reader for '%s', %w", path, err)
				}
			}

			logger.Debug("Yield new record")

			iter_r := NewRecord(path, rsc)
			yield(iter_r, nil)

			return nil
		})

		if err != nil {
			logger.Debug("Failed to walk directory", "error", err)
			yield(nil, err)
			return
		}
	}
}
