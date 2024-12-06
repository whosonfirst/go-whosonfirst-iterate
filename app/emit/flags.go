package emit

import (
	"flag"
	"fmt"
	"os"

	"github.com/sfomuseum/go-flags/flagset"
)

var iterator_uri string

func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("count")
	fs.StringVar(&iterator_uri, "iterator-uri", "", "")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Count files in one or more whosonfirst/go-whosonfirst-iterate/v3 sources.\n")
		fmt.Fprintf(os.Stderr, "Usage:\n\t %s [options] uri(N) uri(N)\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Valid options are:\n\n")
		fs.PrintDefaults()
	}

	return fs
}
