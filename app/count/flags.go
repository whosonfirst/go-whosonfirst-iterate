package count

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/sfomuseum/go-flags/flagset"
	"github.com/whosonfirst/go-whosonfirst-iterate/v3"
)

var iterator_uri string
var verbose bool

// DefaultFlagSet returns a default `flag.FlagSet` for executing a command line application
// to count records with a `go-whosonfirst-iterate/v3.Iterator` instance.
func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("emit")

	valid_schemes := strings.Join(iterate.IteratorSchemes(), ",")
	iterator_desc := fmt.Sprintf("A valid whosonfirst/go-whosonfirst-iterate/v3.Iterator URI. Supported iterator URI schemes are: %s", valid_schemes)

	fs.StringVar(&iterator_uri, "iterator-uri", "repo://", iterator_desc)
	fs.BoolVar(&verbose, "verbose", false, "Enable verbose (debug) logging.")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Count files in one or more whosonfirst/go-whosonfirst-iterate/v3.Iterator sources.\n")
		fmt.Fprintf(os.Stderr, "Usage:\n\t %s [options] uri(N) uri(N)\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Valid options are:\n\n")
		fs.PrintDefaults()
	}

	return fs
}
