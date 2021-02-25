# go-whosonfirst-iterate

Go package for iterating through a set of Who's On First documents

## Important

Documentation for this package is incomplete and will be updated shortly.

## Example

```
package main

import (
       "context"
       "flag"
       "github.com/whosonfirst/go-whosonfirst-iterate/emitter"       
       "github.com/whosonfirst/go-whosonfirst-iterate/indexer"
       "io"
       "log"
)

func main() {

	emitter_uri := flag.String("emitter-uri", "repo://", "A valid whosonfirst/go-whosonfirst-iterate/emitter URI")
	
     	flag.Parse()

	ctx := context.Background()

	emitter_cb := func(ctx context.Context, fh io.ReadSeeker, args ...interface{}) error {
		path, _ := index.PathForContext(ctx)
		log.Printf("Indexing %s\n", path)
		return nil
	}

	iter, _ := iterator.NewIterator(ctx, *emitter_uri, cb)

	uris := flag.Args()
	iter.IterateURIs(ctx, uris...)
}
```

_Error handling removed for the sake of brevity._

## Concepts

### Iterators

### Emitters

_To be written_

## Interfaces

```
type EmitterInitializeFunc func(context.Context, string) (Emitter, error)

type EmitterCallbackFunc func(context.Context, io.ReadSeekCloser, ...interface{}) error

type Emitter interface {
	WalkURI(context.Context, EmitterCallbackFunc, string) error
}
```

_To be written_

## URIs and Schemes 

_To be written_

### directory://

### featurecollection://

### file://

### filelist://

### geojsonls://

### repo://

## Query parameters

| Name | Value | Required | Notes
| --- | --- | --- | --- |
| _max_procs | Int | No | _To be written_ |
| _exclude | String (a valid regular expression) | No | _To be written_ |
| include | String | No | One or more query filters (described below) to limit documents that will be processed. |
| exclude | String | No | One or more query filters (described below) for excluding documents from being processed. |

## Filters

### QueryFilters

You can also specify inline queries by appending one or more `include` or `exclude` parameters to a `emitter.Emitter` URI, where the value is a string in the format of:

```
{PATH}={REGULAR EXPRESSION}
```

Paths follow the dot notation syntax used by the [tidwall/gjson](https://github.com/tidwall/gjson) package and regular expressions are any valid [Go language regular expression](https://golang.org/pkg/regexp/). Successful path lookups will be treated as a list of candidates and each candidate's string value will be tested against the regular expression's [MatchString](https://golang.org/pkg/regexp/#Regexp.MatchString) method.

For example:

```
repo://?include=properties.wof:placetype=region
```

You can pass multiple query parameters. For example:

```
repo://?include=properties.wof:placetype=region&include=properties.wof:name=(?i)new.*
```

The default query mode is to ensure that all queries match but you can also specify that only one or more queries need to match by appending a `include_mode` or `exclude_mode` parameter where the value is either "ANY" or "ALL".

## Tools

```
$> make cli
go build -mod vendor -o bin/count cmd/count/main.go
go build -mod vendor -o bin/emit cmd/emit/main.go
```

### count

Count files in one or more whosonfirst/go-whosonfirst-index/v2/emitter sources.

```
> ./bin/count -h
Count files in one or more whosonfirst/go-whosonfirst-iterate/emitter sources.
Usage:
	 ./bin/count [options] uri(N) uri(N)
Valid options are:

  -emitter-uri string
    	A valid whosonfirst/go-whosonfirst-iterate/emitter URI. Supported emitter URI schemes are: directory://,featurecollection://,file://,filelist://,geojsonl://,repo:// (default "repo://")
```

For example:

```
$> ./bin/count \
	/usr/local/data/sfomuseum-data-architecture/

2021/02/17 14:07:01 time to index paths (1) 87.908997ms
2021/02/17 14:07:01 Counted 1072 records (1072) in 88.045771ms
```

Or:

```
$> ./bin/count \
	-emitter-uri 'repo://?include=properties.sfomuseum:placetype=terminal&include=properties.mz:is_current=1' \
	/usr/local/data/sfomuseum-data-architecture/
	
2021/02/17 14:09:18 time to index paths (1) 71.06355ms
2021/02/17 14:09:18 Counted 4 records (4) in 71.184227ms
```

### emit

Publish features from one or more whosonfirst/go-whosonfirst-index/v2/emitter sources.

```
> ./bin/emit -h
Publish features from one or more whosonfirst/go-whosonfirst-iterate/emitter sources.
Usage:
	 ./bin/emit [options] uri(N) uri(N)
Valid options are:

  -emitter-uri string
    	A valid whosonfirst/go-whosonfirst-iterator/emitter URI. Supported emitter URI schemes are: directory://,featurecollection://,file://,filelist://,geojsonl://,repo:// (default "repo://")
  -geojson
    	Emit features as a well-formed GeoJSON FeatureCollection record.
  -json
    	Emit features as a well-formed JSON array.
  -null
    	Publish features to /dev/null
  -stdout
    	Publish features to STDOUT. (default true)
```

For example:

```
$> ./bin/emit \
	-emitter-uri 'repo://?include=properties.sfomuseum:placetype=museum' \
	-geojson \	
	/usr/local/data/sfomuseum-data-architecture/ \

| jq '.features[]["properties"]["wof:id"]'

1729813675
1477855937
1360521563
1360521569
1360521565
1360521571
1159157863
```

## See also

* https://github.com/aaronland/go-json-query
* https://github.com/aaronland/go-roster