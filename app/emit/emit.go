package emit

import (
	"context"
	"flag"
	"io"
	"log"
	"os"

	"github.com/sfomuseum/go-flags/flagset"
	"github.com/whosonfirst/go-whosonfirst-iterate/v3"
)

func Run(ctx context.Context) error {
	fs := DefaultFlagSet()
	return RunWithFlagSet(ctx, fs)
}

func RunWithFlagSet(ctx context.Context, fs *flag.FlagSet) error {

	flagset.Parse(fs)

	paths := fs.Args()

	iter, err := iterate.NewIterator(ctx, iterator_uri)

	if err != nil {
		log.Fatal(err)
	}

	for rec, err := range iter.Iterate(ctx, paths...) {

		if err != nil {
			log.Fatal(err)
		}

		_, err = io.Copy(os.Stdout, rec.Body)

		if err != nil {
			log.Fatal(err)
		}
	}

	return nil
}
