package iterwriter

import (
	"context"
	"flag"
	"fmt"
	"github.com/sfomuseum/go-flags/flagset"
	iter_writer "github.com/whosonfirst/go-whosonfirst-iterate/v2/writer"
	wof_writer "github.com/whosonfirst/go-writer/v2"
	"log"
)

var writer_uri string
var iterator_uri string

func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("es")

	fs.StringVar(&writer_uri, "writer-uri", "", "")
	fs.StringVar(&iterator_uri, "iterator-uri", "repo://", "")

	return fs
}

func Run(ctx context.Context, logger *log.Logger) error {
	fs := DefaultFlagSet()
	return RunWithFlagSet(ctx, fs, logger)
}

func RunWithFlagSet(ctx context.Context, fs *flag.FlagSet, logger *log.Logger) error {

	flagset.Parse(fs)

	iterator_paths := fs.Args()

	wr, err := wof_writer.NewWriter(ctx, writer_uri)

	if err != nil {
		return fmt.Errorf("Failed to create new writer, %w", err)
	}

	err = iter_writer.IterateWithWriter(ctx, wr, iterator_uri, iterator_paths...)

	if err != nil {
		return fmt.Errorf("Failed to iterate with writer, %w", err)
	}

	return nil
}
