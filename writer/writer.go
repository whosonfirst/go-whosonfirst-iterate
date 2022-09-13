package writer

import (
	"context"
	"fmt"
	"github.com/whosonfirst/go-whosonfirst-iterate/v2/iterator"
	wof_writer "github.com/whosonfirst/go-writer/v2"
	"io"
)

func IterateWithWriter(ctx context.Context, wr wof_writer.Writer, iterator_uri string, iterator_paths ...string) error {

	iter_cb := func(ctx context.Context, path string, r io.ReadSeeker, args ...interface{}) error {

		_, err := wr.Write(ctx, path, r)

		if err != nil {
			return fmt.Errorf("Failed to write %s, %v", path, err)
		}

		return nil
	}

	iter, err := iterator.NewIterator(ctx, iterator_uri, iter_cb)

	if err != nil {
		return fmt.Errorf("Failed to create new iterator, %w", err)
	}

	err = iter.IterateURIs(ctx, iterator_paths...)

	if err != nil {
		return fmt.Errorf("Failed to iterate paths, %w", err)
	}

	err = wr.Close(ctx)

	if err != nil {
		return fmt.Errorf("Failed to close ES writer, %w", err)
	}

	return nil
}
