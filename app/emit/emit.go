package emit

import (
	"context"
	"flag"
	"io"
	"os"

	"github.com/sfomuseum/go-flags/flagset"
	"github.com/whosonfirst/go-whosonfirst-iterate/v3"
)

// Run will execute a command line application for emittingrecords with a `go-whosonfirst-iterate/v3.Iterator`
// instance using a default flagset.
func Run(ctx context.Context) error {
	fs := DefaultFlagSet()
	return RunWithFlagSet(ctx, fs)
}

// RunWithFlagSet will execute a command line application for emitting records with a `go-whosonfirst-iterate/v3.Iterator`
// instance using 'fs'
func RunWithFlagSet(ctx context.Context, fs *flag.FlagSet) error {

	flagset.Parse(fs)

	paths := fs.Args()

	iter, err := iterate.NewIterator(ctx, iterator_uri)

	if err != nil {
		return err
	}

	for rec, err := range iter.Iterate(ctx, paths...) {

		if err != nil {
			return err
		}

		_, err = io.Copy(os.Stdout, rec.Body)

		if err != nil {
			return err
		}
	}

	return nil
}
