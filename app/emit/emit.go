package emit

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/sfomuseum/go-flags/flagset"
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

	if verbose {
		slog.SetLogLoggerLevel(slog.LevelDebug)
		slog.Debug("Verbose logging enabled")
	}

	writers := make([]io.Writer, 0)

	if to_stdout {
		writers = append(writers, os.Stdout)
	}

	if to_devnull {
		writers = append(writers, io.Discard)
	}

	wr := io.MultiWriter(writers...)

	em := &FeatureEmitter{
		AsJSON:    as_json,
		AsGeoJSON: as_geojson,
		Writer:    wr,
	}

	uris := fs.Args()

	_, err := em.Emit(ctx, iterator_uri, uris...)

	if err != nil {
		return fmt.Errorf("Failed to emit records, %w", err)
	}

	return nil
}
