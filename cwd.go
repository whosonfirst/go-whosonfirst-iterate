package iterate

import (
	"context"
	"fmt"
	"iter"
	"os"
)

func init() {
	ctx := context.Background()
	RegisterIterator(ctx, "cwd", NewCwdIterator)
}

// CwdIterator implements the `Iterator` interface for crawling records in a Who's On First style data directory.
type CwdIterator struct {
	Iterator
	// iterator is the underlying `DirectoryIterator` instance for crawling records.
	iterator Iterator
}

// NewDirectoryIterator() returns a new `CwdIterator` instance configured by 'uri' in the form of:
//
//	cwd://?{PARAMETERS}
//
// Where {PARAMETERS} may be:
// * `?include=` Zero or more `aaronland/go-json-query` query strings containing rules that must match for a document to be considered for further processing.
// * `?exclude=` Zero or more `aaronland/go-json-query`	query strings containing rules that if matched will prevent a document from being considered for further processing.
// * `?include_mode=` A valid `aaronland/go-json-query` query mode string for testing inclusion rules.
// * `?exclude_mode=` A valid `aaronland/go-json-query` query mode string for testing exclusion rules.
func NewCwdIterator(ctx context.Context, uri string) (Iterator, error) {

	directory_it, err := NewDirectoryIterator(ctx, uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to create new directory iterator, %w", err)
	}

	it := &CwdIterator{
		iterator: directory_it,
	}

	return it, nil
}

func (it *CwdIterator) Iterate(ctx context.Context, uris ...string) iter.Seq2[*Record, error] {

	cwd, err := os.Getwd()

	if err != nil {
		return func(yield func(rec *Record, err error) bool) {
			yield(nil, err)
		}
	}

	return it.iterator.Iterate(ctx, cwd)
}

func (it *CwdIterator) Seen() int64 {
	return it.iterator.Seen()
}

func (it *CwdIterator) IsIterating() bool {
	return it.iterator.IsIterating()
}
