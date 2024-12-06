package iterator

import (
	"context"
	"iter"
)

func init() {
	ctx := context.Background()
	RegisterIterator(ctx, "null", NewNullIterator)
}

// NullIterator implements the `Iterator` interface for appearing to crawl records but not doing anything.
type NullIterator struct {
	Iterator
}

// NewNullIterator() returns a new `NullIterator` instance configured by 'uri' in the form of:
//
//	null://
func NewNullIterator(ctx context.Context, uri string) (Iterator, error) {

	idx := &NullIterator{}
	return idx, nil
}

// WalkURI() does nothing.
func (idx *NullIterator) Iterate(ctx context.Context, uris ...string) iter.Seq2[*Record, error] {

	return func(yield func(*Record, error) bool) {
		return
	}
}
