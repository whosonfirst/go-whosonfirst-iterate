package emit

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"

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

	it, err := iterate.NewIterator(ctx)

	if err != nil {
		return fmt.Errorf("Failed to create new iterator, %w", err)
	}

	for r, err := range it.Iterate(ctx, opts.IteratorURI, opts.IteratorSources...) {

		if err != nil {
			return fmt.Errorf("Iterator reported an error, %w", err)
		}

		_, err = io.Copy(os.Stdout, r.Body())

		if err != nil {
			return fmt.Errorf("Failed to copy %s to STDOUT, %w", r.URI(), err)
		}
	}

	return nil
}
