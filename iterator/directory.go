package iterator

import (
	"context"
	"fmt"
	"io/fs"
	"iter"
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

	idx := &DirectoryIterator{
		filters: f,
	}

	return idx, nil
}

// WalkURI() walks (crawls) the directory named 'uri' and for each file (not excluded by any filters specified
// when `idx` was created) invokes 'index_cb'.
func (idx *DirectoryIterator) Iterate(ctx context.Context, uris ...string) iter.Seq2[Record, error] {

	return func(yield func(Record, error) bool) {

		for _, uri := range uris {

			for r, err := range idx.iterate(ctx, uri) {
				yield(r, err)
			}
		}
	}
}

func (idx *DirectoryIterator) iterate(ctx context.Context, uri string) iter.Seq2[Record, error] {

	return func(yield func(Record, error) bool) {

		abs_path, err := filepath.Abs(uri)

		if err != nil {
			yield(nil, fmt.Errorf("Failed to derive absolute path for '%s', %w", uri, err))
			return
		}

		dir_fs := os.DirFS(abs_path)

		err = fs.WalkDir(dir_fs, ".", func(path string, d fs.DirEntry, err error) error {

			if err != nil {
				return err
			}

			if d.IsDir() {
				return nil
			}

			r, err := dir_fs.Open(path)

			if err != nil {
				return fmt.Errorf("Failed to open %s for reading, %w", path, err)
			}

			rsc, err := ioutil.NewReadSeekCloser(r)

			if err != nil {
				return fmt.Errorf("Failed to create ReadSeekCloser for %s, %w", path, err)
			}

			defer rsc.Close()

			if idx.filters != nil {

				ok, err := idx.filters.Apply(ctx, rsc)

				if err != nil {
					return fmt.Errorf("Failed to apply filters for '%s', %w", path, err)
				}

				if !ok {
					return nil
				}

				_, err = rsc.Seek(0, 0)

				if err != nil {
					return fmt.Errorf("Failed to seek(0, 0) on reader for '%s', %w", path, err)
				}
			}

			iter_r := NewRecord(path, rsc)

			yield(iter_r, nil)
			return nil
		})

		if err != nil {
			yield(nil, err)
			return
		}
	}
}
