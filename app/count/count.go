package count

import (
	"context"
	"flag"
	"log/slog"
	"sync/atomic"
	"time"

	"github.com/sfomuseum/go-flags/flagset"
	"github.com/whosonfirst/go-whosonfirst-iterate/v3"
)

// Run will execute a command line application to count records with a `go-whosonfirst-iterate/v3.Iterator`
// instance using a default flagset.
func Run(ctx context.Context) error {
	fs := DefaultFlagSet()
	return RunWithFlagSet(ctx, fs)
}

// RunWithFlagSet will execute a command line application to count records with a `go-whosonfirst-iterate/v3.Iterator`
// instance using 'fs'
func RunWithFlagSet(ctx context.Context, fs *flag.FlagSet) error {

	flagset.Parse(fs)

	if verbose {
		slog.SetLogLoggerLevel(slog.LevelDebug)
		slog.Debug("Verbose logging enabled")
	}

	paths := fs.Args()

	count := int64(0)

	iter, err := iterate.NewIterator(ctx, iterator_uri)

	if err != nil {
		return err
	}

	t1 := time.Now()

	for _, err := range iter.Iterate(ctx, paths...) {

		if err != nil {
			return err
		}

		atomic.AddInt64(&count, 1)
	}

	slog.Info("Counted records", "count", count, "time", time.Since(t1))
	return nil
}
