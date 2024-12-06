package iterator

import (
	"context"
	"fmt"
	"iter"
	"os"

	"github.com/whosonfirst/go-whosonfirst-iterate/v3/filters"
)

func init() {
	ctx := context.Background()
	RegisterIterator(ctx, "cwd", NewCwdIterator)
}

// CwdIterator implements the `Iterator` interface for crawling records in the current (working) directory.
type CwdIterator struct {
	Iterator
	// filters is a `filters.Filters` instance used to include or exclude specific records from being crawled.
	filters  filters.Filters
	iterator Iterator
}

// NewCwdIterator() returns a new `CwdIterator` instance configured by 'uri' in the form of:
//
//	cwd://?{PARAMETERS}
//
// Where {PARAMETERS} may be:
// * `?include=` Zero or more `aaronland/go-json-query` query strings containing rules that must match for a document to be considered for further processing.
// * `?exclude=` Zero or more `aaronland/go-json-query`	query strings containing rules that if matched will prevent a document from being considered for further processing.
// * `?include_mode=` A valid `aaronland/go-json-query` query mode string for testing inclusion rules.
// * `?exclude_mode=` A valid `aaronland/go-json-query` query mode string for testing exclusion rules.
func NewCwdIterator(ctx context.Context, uri string) (Iterator, error) {

	directory_idx, err := NewDirectoryIterator(ctx, uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to create new directory iterator, %w", err)
	}

	f, err := filters.NewQueryFiltersFromURI(ctx, uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to create filters from query, %w", err)
	}

	it := &CwdIterator{
		filters:  f,
		iterator: directory_idx,
	}

	return it, nil
}

func (it *CwdIterator) Iterate(ctx context.Context, uris ...string) iter.Seq2[Record, error] {

	return func(yield func(Record, error) bool) {

		cwd, err := os.Getwd()

		if err != nil {
			yield(nil, fmt.Errorf("Failed to derive current working directory, %w", err))
			return
		}

		for r, err := range it.iterator.Iterate(ctx, cwd) {
			yield(r, err)
		}
	}
}
