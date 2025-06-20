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
       "io"
       "log"

       "github.com/whosonfirst/go-whosonfirst-iterate/v2/emitter"       
       "github.com/whosonfirst/go-whosonfirst-iterate/v2/iterator"
)

func main() {

	emitter_uri := flag.String("emitter-uri", "repo://", "A valid whosonfirst/go-whosonfirst-iterate/emitter URI")
	
     	flag.Parse()

	ctx := context.Background()

	emitter_cb := func(ctx context.Context, path string, fh io.ReadSeeker, args ...interface{}) error {
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

The naming conventions (`iterator` and `emitter` and `publisher`) are not ideal. They may still be changed. Briefly:

* An "iterator" is a high-level construct that manages the dispatching and processing of multiple source URIs.

* An "emitter" is the code that walks (or "crawls") a given URI and emits documents to be proccesed by a user-defined callback function. Emitters are defined by the `emitter.Emitter` interface.

* A "publisher" is a higher-order construct that bundles an internal iterator with its own callback function to republish data derived from an iterator/emitter to an `io.Writer` target.

## URIs and Schemes (for emitters)

The following emitters are supported by default:

### cwd://

`CwdEmitter` implements the `Emitter` interface for crawling records in the current working directory.

### directory://

`DirectoryEmitter` implements the `Emitter` interface for crawling records in a directory.

### featurecollection://

`FeatureCollectionEmitter` implements the `Emitter` interface for crawling features in a GeoJSON FeatureCollection record.

### file://

`FileEmitter` implements the `Emitter` interface for crawling individual file records.

### filelist://

`FileListEmitter` implements the `Emitter` interface for crawling records listed in a "file list" (a plain text newline-delimted list of files).

### geojsonl://

`GeojsonLEmitter` implements the `Emitter` interface for crawling features in a line-separated GeoJSON record.

### null://

`NullEmitter` implements the `Emitter` interface for appearing to crawl records but not doing anything.

### repo://

`RepoEmitter` implements the `Emitter` interface for crawling records in a Who's On First style data directory.

## Query parameters

The following query parameters are honoured by all `emitter.Emitter` instances:

| Name | Value | Required | Notes
| --- | --- | --- | --- |
| include | String | No | One or more query filters (described below) to limit documents that will be processed. |
| exclude | String | No | One or more query filters (described below) for excluding documents from being processed. |

The following query paramters are honoured for `emitter.Emitter` URIs passed to the `iterator.NewIterator` method:

| Name | Value | Required | Notes
| --- | --- | --- | --- |
| _max_procs | Int | No | _To be written_ |
| _exclude | String (a valid regular expression) | No | _To be written_ |
| _exclude_alt | Bool | No | If true do not process "alternate geometry" files. |
| _retry | Bool | No | A boolean flag signaling that if a URI being walked fails it should be retried. Used in conjunction with the `_max_retries` and `_retry_after` parameters. |
| _max_retries | Int | No | The maximum number of attempts to walk any given URI. Defaults to "1" and the `_retry` parameter _must_ evaluate to a true value in order to change the default. |
| _retry_after | Int | The number of seconds to wait between attempts to walk any given URI. Defaults to "10" (seconds) and the `_retry` parameter _must_ evaluate to a true value in order to change the default. |
| _dedupe | Bool | No | A boolean value to track and skip records (specifically their relative URI) that have already been processed. |

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

## Notes about writing your own `iterate.Iterator` implementation.

Under the hood all `iterate.Iterate` instances are wrapped using the (private) `concurrentIterator` implementation. This is the code that implements throttling, file matching and other common tasks.

Importantly, it also takes care of automatically closing any `Record.Body` instances after that `Record` instance has been yielded (and the `yield` function completes).

As a consequence you should _not_ automatically close `Record.Body` instances in your own code using the common `defer rec.Body.Close()` idiom. This is unfortunate because it makes ensuring that those instances are closed after they are opened but, for whatever reasons, not scheduled to be yielded. This is a by-product of the way that Go `yield` functions work and the extra work to ensure that filehandles are closed (when not being yielded) is just the "cost of doing business" I guess.

_And yes, I did try using Go 1.24's `runtime.AddCleanup` but because it execute as part of the runtime.GC process it often gets triggered after the *Record instance has been purged without closing the underlying file handle. Basically what we need is a Python-style object level destructor but those don't exist yet so, again, here we are._

## Version

## v3

## v2

Version `2.x` of this package was released to address a problem with the way version 1.x was passing path names (or URIs) for files being processed: Namely [it wasn't thread-safe](https://github.com/whosonfirst/go-whosonfirst-iterate/issues/5) so it was possible to derive a path (from a context) that was associated with another file. Version 2.x changes the interface for local callback to include the string path (or URI) for the file being processed.

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