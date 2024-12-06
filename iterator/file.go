package iterator

import (
	"context"
	"fmt"
	"iter"
	"log/slog"

	"github.com/whosonfirst/go-whosonfirst-iterate/v3/filters"
)

func init() {
	ctx := context.Background()
	RegisterIterator(ctx, "file", NewFileIterator)
}

// FileIterator implements the `Iterator` interface for crawling individual file records.
type FileIterator struct {
	Iterator
	// filters is a `filters.Filters` instance used to include or exclude specific records from being crawled.
	filters filters.Filters
}

// NewFileIterator() returns a new `FileIterator` instance configured by 'uri' in the form of:
//
//	file://?{PARAMETERS}
//
// Where {PARAMETERS} may be:
// * `?include=` Zero or more `aaronland/go-json-query` query strings containing rules that must match for a document to be considered for further processing.
// * `?exclude=` Zero or more `aaronland/go-json-query`	query strings containing rules that if matched will prevent a document from being considered for further processing.
// * `?include_mode=` A valid `aaronland/go-json-query` query mode string for testing inclusion rules.
// * `?exclude_mode=` A valid `aaronland/go-json-query` query mode string for testing exclusion rules.
func NewFileIterator(ctx context.Context, uri string) (Iterator, error) {

	f, err := filters.NewQueryFiltersFromURI(ctx, uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to create filters from query, %w", err)
	}

	it := &FileIterator{
		filters: f,
	}

	return it, nil
}

func (it *FileIterator) Iterate(ctx context.Context, uris ...string) iter.Seq2[Record, error] {

	return func(yield func(Record, error) bool) {

		for _, uri := range uris {
			for r, err := range it.iterate(ctx, uri) {
				yield(r, err)
			}
		}
	}
}

func (it *FileIterator) iterate(ctx context.Context, uri string) iter.Seq2[Record, error] {

	logger := slog.Default()
	logger = logger.With("uri", uri)

	return func(yield func(Record, error) bool) {

		r, err := ReaderWithPath(ctx, uri)

		if err != nil {
			logger.Debug("Failed to create reader", "error", err)
			yield(nil, fmt.Errorf("Failed to create reader for '%s', %w", uri, err))
			return
		}

		defer r.Close()

		if it.filters != nil {

			ok, err := it.filters.Apply(ctx, r)

			if err != nil {
				logger.Debug("Failed to apply filters", "error", err)
				yield(nil, fmt.Errorf("Failed to apply filters for '%s', %w", uri, err))
				return
			}

			if !ok {
				logger.Debug("No matches after applying filters, skipping")
				return
			}

			_, err = r.Seek(0, 0)

			if err != nil {
				logger.Debug("Failed to rewind reader", "error", err)
				yield(nil, fmt.Errorf("Failed to seek(0,) for '%s', %w", uri, err))
				return
			}
		}

		logger.Debug("Yield new record")
		iter_r := NewRecord(uri, r)
		yield(iter_r, nil)
	}

}
