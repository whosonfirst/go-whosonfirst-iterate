package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/whosonfirst/go-whosonfirst-iterate/v3"
)

func main() {

	valid_schemes := strings.Join(iterate.IteratorSchemes(), ",")
	iterator_desc := fmt.Sprintf("A valid whosonfirst/go-whosonfirst-iterate/v3.Iterator URI. Supported iterator URI schemes are: %s", valid_schemes)

	var iterator_uri string

	flag.StringVar(&iterator_uri, "iterator-uri", "repo://", iterator_desc)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Count files in one or more whosonfirst/go-whosonfirst-iterate.Iterator sources.\n")
		fmt.Fprintf(os.Stderr, "Usage:\n\t %s [options] uri(N) uri(N)\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Valid options are:\n\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	ctx := context.Background()

	paths := flag.Args()

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

}
