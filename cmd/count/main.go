package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"sync/atomic"
	"time"

	"github.com/whosonfirst/go-whosonfirst-iterate/v3/iterator"
)

func main() {

	var iterator_uri string

	flag.StringVar(&iterator_uri, "iterator-uri", "repo://", "...")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Count files in one or more whosonfirst/go-whosonfirst-iterate/v3/iterator sources.\n")
		fmt.Fprintf(os.Stderr, "Usage:\n\t %s [options] uri(N) uri(N)\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Valid options are:\n\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	ctx := context.Background()

	var count int64
	count = 0

	iter, err := iterator.NewIterator(ctx, iterator_uri)

	if err != nil {
		log.Fatal(err)
	}

	paths := flag.Args()

	t1 := time.Now()

	for _, err := range iter.Walk(ctx, paths...) {

		if err != nil {
			log.Fatal(err)
		}

		atomic.AddInt64(&count, 1)
	}

	log.Printf("Counted %d records (%d) in %v\n", count, count, time.Since(t1))
}
