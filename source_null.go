package iterate

import (
	"context"
	"iter"
)

func init() {
	// ctx := context.Background()
	// RegisterSource(ctx, "null", NewNullSource)
}

// NullSource implements the `Source` interface for appearing to crawl records but not doing anything.
type NullSource struct {
	Source
}

// NewNullSource() returns a new `NullSource` instance configured by 'uri' in the form of:
//
//	null://
func NewNullSource(ctx context.Context, uri string) (Source, error) {

	idx := &NullSource{}
	return idx, nil
}

// WalkURI() does nothing.
func (idx *NullSource) Walk(ctx context.Context, uri string) iter.Seq2[*Record, error] {

	return func(yield func(rec *Record, err error) bool) {}
}
