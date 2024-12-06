package count

import (
	"context"
	"flag"
	"fmt"
	"log"
	"sync/atomic"

	"github.com/whosonfirst/go-whosonfirst-iterate/v3"
)

func Run(ctx context.Context) error {
	fs := DefaultFlagSet()
	return RunWithFlagSet(ctx, fs)
}

func RunWithFlagSet(ctx context.Context, fs *flag.FlagSet) error {
	opts := RunOptionsFromFlagSet(fs)
	return RunWithOptions(ctx, opts)
}

func RunWithOptions(ctx context.Context, opts *RunOptions) error {

	var count int64
	count = 0

	it, err := iterate.NewIterator(ctx)

	if err != nil {
		return fmt.Errorf("Failed to create new iterator, %w", err)
	}

	for _, err := range it.Iterate(ctx, opts.IteratorURI, opts.IteratorSources...) {

		if err != nil {
			return fmt.Errorf("Iterator reported an error, %w", err)
		}

		atomic.AddInt64(&count, 1)
	}

	log.Printf("Counted %d records (saw %d records)\n", count, it.Seen)
	return nil
}
