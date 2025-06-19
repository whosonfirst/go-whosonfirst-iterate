package iterate

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

	it := &NullIterator{}
	return it, nil
}

// Iterate() does nothing.
func (it *NullIterator) Iterate(ctx context.Context, uris ...string) iter.Seq2[*Record, error] {
	return func(yield func(rec *Record, err error) bool) {}
}

func (it *NullIterator) Seen() int64 {
	return int64(0)
}

func (it *NullIterator) IsIterator() bool {
	return false
}
