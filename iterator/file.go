package iterator

import (
	"context"
	"fmt"
	"iter"

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

	idx := &FileIterator{
		filters: f,
	}

	return idx, nil
}

func (idx *FileIterator) Iterate(ctx context.Context, uris ...string) iter.Seq2[Record, error] {

	return func(yield func(Record, error) bool) {

		for _, uri := range uris {

			fh, err := ReaderWithPath(ctx, uri)

			if err != nil {
				yield(nil, fmt.Errorf("Failed to create reader for '%s', %w", uri, err))
				continue
			}

			defer fh.Close()

			if idx.filters != nil {

				ok, err := idx.filters.Apply(ctx, fh)

				if err != nil {
					yield(nil, fmt.Errorf("Failed to apply filters for '%s', %w", uri, err))
					continue
				}

				if !ok {
					continue
				}

				_, err = fh.Seek(0, 0)

				if err != nil {
					yield(nil, fmt.Errorf("Failed to seek(0,) for '%s', %w", uri, err))
					continue
				}
			}

			iter_r := NewRecord(uri, fh)
			yield(iter_r, nil)
		}
	}

}
