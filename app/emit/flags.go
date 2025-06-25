package emit

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

var as_json bool
var as_geojson bool

var to_stdout bool
var to_devnull bool

// DefaultFlagSet returns a default `flag.FlagSet` for executing a command line application
// to emitting records with a `go-whosonfirst-iterate/v3.Iterator` instance.
func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("emit")

	valid_schemes := strings.Join(iterate.IteratorSchemes(), ",")
	iterator_desc := fmt.Sprintf("A valid whosonfirst/go-whosonfirst-iterate/v3.Iterator URI. Supported iterator URI schemes are: %s", valid_schemes)

	fs.StringVar(&iterator_uri, "iterator-uri", "repo://", iterator_desc)
	fs.BoolVar(&verbose, "verbose", false, "Enable verbose (debug) logging.")

	fs.BoolVar(&as_json, "json", false, "Emit features as a well-formed JSON array.")
	fs.BoolVar(&as_geojson, "geojson", false, "Emit features as a well-formed GeoJSON FeatureCollection record.")

	fs.BoolVar(&to_stdout, "stdout", true, "Publish features to STDOUT.")
	fs.BoolVar(&to_devnull, "null", false, "Publish features to /dev/null")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Emit records in one or more whosonfirst/go-whosonfirst-iterate/v3.Iterator sources as structured data.\n")
		fmt.Fprintf(os.Stderr, "Usage:\n\t %s [options] uri(N) uri(N)\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Valid options are:\n\n")
		fs.PrintDefaults()
	}

	return fs
}
