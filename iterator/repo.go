package iterator

import (
	"context"
	"fmt"
	"iter"
	"path/filepath"
)

func init() {
	ctx := context.Background()
	RegisterIterator(ctx, "repo", NewRepoIterator)
}

// RepoIterator implements the `Iterator` interface for crawling records in a Who's On First style data directory.
type RepoIterator struct {
	Iterator
	// iterator is the underlying `DirectoryIterator` instance for crawling records.
	iterator Iterator
}

// NewDirectoryIterator() returns a new `RepoIterator` instance configured by 'uri' in the form of:
//
//	repo://?{PARAMETERS}
//
// Where {PARAMETERS} may be:
// * `?include=` Zero or more `aaronland/go-json-query` query strings containing rules that must match for a document to be considered for further processing.
// * `?exclude=` Zero or more `aaronland/go-json-query`	query strings containing rules that if matched will prevent a document from being considered for further processing.
// * `?include_mode=` A valid `aaronland/go-json-query` query mode string for testing inclusion rules.
// * `?exclude_mode=` A valid `aaronland/go-json-query` query mode string for testing exclusion rules.
func NewRepoIterator(ctx context.Context, uri string) (Iterator, error) {

	directory_idx, err := NewDirectoryIterator(ctx, uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to create new directory iterator, %w", err)
	}

	idx := &RepoIterator{
		iterator: directory_idx,
	}

	return idx, nil
}

// WalkURI() appends 'uri' with "data" and then walks that directory and for each file (not excluded by any
// filters specified when `idx` was created) invokes 'index_cb'.
func (idx *RepoIterator) Iterate(ctx context.Context, uris ...string) iter.Seq2[*Record, error] {

	return func(yield func(*Record, error) bool) {

		for _, uri := range uris {
			for r, err := range idx.iterate(ctx, uri) {
				yield(r, err)
			}
		}
	}
}

func (idx *RepoIterator) iterate(ctx context.Context, uri string) iter.Seq2[*Record, error] {

	return func(yield func(*Record, error) bool) {

		abs_path, err := filepath.Abs(uri)

		if err != nil {
			yield(nil, fmt.Errorf("Failed to derive absolute path for '%s', %w", uri, err))
			return
		}

		data_path := filepath.Join(abs_path, "data")

		for r, err := range idx.iterator.Iterate(ctx, data_path) {
			yield(r, err)
		}
	}
}
