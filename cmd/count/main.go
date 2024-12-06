package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"sync/atomic"
	"time"

	"github.com/whosonfirst/go-whosonfirst-iterate/v3"
)

func main() {

	var emitter_uri = flag.String("emitter-uri", "repo://", "")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Count files in one or more whosonfirst/go-whosonfirst-iterate/emitter sources.\n")
		fmt.Fprintf(os.Stderr, "Usage:\n\t %s [options] uri(N) uri(N)\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Valid options are:\n\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	ctx := context.Background()

	var count int64
	count = 0

	iter, err := iterate.NewIterator(ctx)

	if err != nil {
		log.Fatal(err)
	}

	paths := flag.Args()

	t1 := time.Now()

	for _, err := range iter.Iterate(ctx, *emitter_uri, paths...) {

		if err != nil {
			log.Fatal(err)
		}

		atomic.AddInt64(&count, 1)
	}

	log.Printf("Counted %d records (%d) in %v\n", count, iter.Seen, time.Since(t1))
}
