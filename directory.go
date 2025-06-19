package iterate

import (
	"context"
	"fmt"
	"io"
	"iter"
	_ "log/slog"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"

	"github.com/whosonfirst/go-whosonfirst-crawl"
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
	filters   filters.Filters
	seen      int64
	iterating *atomic.Bool
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
		filters:   f,
		seen:      int64(0),
		iterating: new(atomic.Bool),
	}

	return it, nil
}

func (it *DirectoryIterator) Iterate(ctx context.Context, uris ...string) iter.Seq2[*Record, error] {

	return func(yield func(rec *Record, err error) bool) {

		it.iterating.Swap(true)
		defer it.iterating.Swap(false)

		for _, uri := range uris {

			abs_path, err := filepath.Abs(uri)

			if err != nil {
				yield(nil, fmt.Errorf("Failed to derive absolute path for '%s', %w", uri, err))
			}

			mu := new(sync.RWMutex)

			crawl_cb := func(path string, info os.FileInfo) error {

				select {
				case <-ctx.Done():
					return nil
				default:
					// pass
				}

				if info.IsDir() {
					return nil
				}

				atomic.AddInt64(&it.seen, 1)

				r, err := ReaderWithPath(ctx, path)

				if err != nil {
					return fmt.Errorf("Failed to create reader for '%s', %w", abs_path, err)
				}

				if it.filters != nil {

					ok, err := it.filters.Apply(ctx, r)

					if err != nil {
						r.Close()
						return fmt.Errorf("Failed to apply filters for '%s', %w", abs_path, err)
					}

					if !ok {
						r.Close()
						return nil
					}

					_, err = r.Seek(0, 0)

					if err != nil {
						r.Close()
						return fmt.Errorf("Failed to seek(0, 0) on reader for '%s', %w", abs_path, err)
					}
				}

				rec := NewRecord(path, r)

				mu.Lock()
				defer mu.Unlock()

				if !yield(rec, nil) {
					return io.EOF
				}

				return nil
			}

			c := crawl.NewCrawler(abs_path)
			err = c.Crawl(crawl_cb)

			if err != nil && err != io.EOF {
				yield(nil, err)
			}
		}
	}
}

func (it *DirectoryIterator) Seen() int64 {
	return atomic.LoadInt64(&it.seen)
}

func (it *DirectoryIterator) IsIterating() bool {
	return it.iterating.Load()
}
