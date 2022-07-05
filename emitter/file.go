package emitter

import (
	"context"
	"fmt"
	"github.com/whosonfirst/go-whosonfirst-iterate/v2/filters"
)

func init() {
	ctx := context.Background()
	RegisterEmitter(ctx, "file", NewFileEmitter)
}

// FileEmitter implements the `Emitter` interface for crawling individual file records.
type FileEmitter struct {
	Emitter
	// filters is a `filters.Filters` instance used to include or exclude specific records from being crawled.
	filters filters.Filters
}

// NewFileCollectionEmitter() returns a new `FileEmitter` instance configured by 'uri' in the form of:
//
//	file://?{PARAMETERS}
func NewFileEmitter(ctx context.Context, uri string) (Emitter, error) {

	f, err := filters.NewQueryFiltersFromURI(ctx, uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to create filters from query, %w", err)
	}

	idx := &FileEmitter{
		filters: f,
	}

	return idx, nil
}

// WalkURI() applies 'index_cb' to the file named 'uri'.
func (idx *FileEmitter) WalkURI(ctx context.Context, index_cb EmitterCallbackFunc, uri string) error {

	fh, err := ReaderWithPath(ctx, uri)

	if err != nil {
		return fmt.Errorf("Failed to create reader for '%s', %w", uri, err)
	}

	defer fh.Close()

	if idx.filters != nil {

		ok, err := idx.filters.Apply(ctx, fh)

		if err != nil {
			return fmt.Errorf("Failed to apply filters for '%s', %w", uri, err)
		}

		if !ok {
			return nil
		}

		_, err = fh.Seek(0, 0)

		if err != nil {
			return fmt.Errorf("Failed to seek(0,) for '%s', %w", uri, err)
		}
	}

	err = index_cb(ctx, uri, fh)

	if err != nil {
		return fmt.Errorf("Index callback failed for '%s', %w", uri, err)
	}

	return nil
}
