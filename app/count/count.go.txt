package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"os"
	"strings"
	"sync/atomic"
	"time"

	"github.com/whosonfirst/go-whosonfirst-iterate/v3"
)

func main() {

	valid_schemes := strings.Join(iterate.IteratorSchemes(), ",")
	iterator_desc := fmt.Sprintf("A valid whosonfirst/go-whosonfirst-iterate/v3.Iterator URI. Supported iterator URI schemes are: %s", valid_schemes)

	var iterator_uri string

	flag.StringVar(&iterator_uri, "iterator-uri", "repo://", iterator_desc)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Count files in one or more whosonfirst/go-whosonfirst-iterate/emitter sources.\n")
		fmt.Fprintf(os.Stderr, "Usage:\n\t %s [options] uri(N) uri(N)\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Valid options are:\n\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	ctx := context.Background()

	paths := flag.Args()
	var count int64
	count = 0

	iter, err := iterate.NewIterator(ctx, iterator_uri)

	if err != nil {
		log.Fatal(err)
	}

	t1 := time.Now()

	for _, err := range iter.Iterate(ctx, paths...) {

		if err != nil {
			log.Fatal(err)
		}

		atomic.AddInt64(&count, 1)
	}

	slog.Info("Counted records", "count", count, "time", time.Since(t1))
}
