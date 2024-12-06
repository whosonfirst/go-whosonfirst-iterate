package count

import (
	"flag"

	"github.com/sfomuseum/go-flags/flagset"
)

type RunOptions struct {
	IteratorURI     string
	IteratorSources []string
}

func RunOptionsFromFlagSet(fs *flag.FlagSet) *RunOptions {

	flagset.Parse(fs)

	opts := &RunOptions{
		IteratorURI:     iterator_uri,
		IteratorSources: fs.Args(),
	}

	return opts
}
