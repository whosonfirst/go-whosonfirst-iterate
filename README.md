# go-whosonfirst-iterate

Go package for iterating through a set of Who's On First documents

## Documentation

[![Go Reference](https://pkg.go.dev/badge/github.com/whosonfirst/go-whosonfirst-iterate.svg)](https://pkg.go.dev/github.com/whosonfirst/go-whosonfirst-iterate)

## Example

```
package main

import (
       "context"
       "flag"
       "log"

       "github.com/whosonfirst/go-whosonfirst-iterate/v3"
)

func main() {

	iter_uri := flag.String("iterator-uri", "repo://", "A valid whosonfirst/go-whosonfirst-iterate/v3/iterator URI.")
	
     	flag.Parse()

	iter_sources := flag.Args()

	ctx := context.Background()
	
	it, _ := iterate.NewIterator(ctx)

	for r, _ := range it.Iterate(ctx, iter_uri, iter_sources...){
		log.Printf("Iterating %s\n", r.URI())	    
	}
}
```

_Error handling removed for the sake of brevity._

## "v3"

The `/v3` release is a major, and backwards incompatible, refactoring of previous versions of this package.

## Related

* https://github.com/whosonfirst/go-whosonfirst-iterate-bucket
* https://github.com/whosonfirst/go-whosonfirst-iterate-git
* https://github.com/whosonfirst/go-whosonfirst-iterate-github
* https://github.com/whosonfirst/go-whosonfirst-iterate-organization
* https://github.com/whosonfirst/go-whosonfirst-iterate-reader
* https://github.com/whosonfirst/go-whosonfirst-iterate-sqlite
* https://github.com/whosonfirst/go-whosonfirst-iterate-fs

## See also

* https://github.com/aaronland/go-json-query
* https://github.com/aaronland/go-roster