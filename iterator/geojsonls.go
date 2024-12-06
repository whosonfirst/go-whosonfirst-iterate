package iterator

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"iter"

	"github.com/whosonfirst/go-ioutil"
	"github.com/whosonfirst/go-whosonfirst-iterate/v3/filters"
)

func init() {
	ctx := context.Background()
	RegisterIterator(ctx, "geojsonl", NewGeoJSONLIterator)
}

// GeojsonLIterator implements the `Iterator` interface for crawling features in a line-separated GeoJSON record.
type GeojsonLIterator struct {
	Iterator
	// filters is a `filters.Filters` instance used to include or exclude specific records from being crawled.
	filters filters.Filters
}

// NewGeojsonLIterator() returns a new `GeojsonLIterator` instance configured by 'uri' in the form of:
//
//	geojsonl://?{PARAMETERS}
//
// Where {PARAMETERS} may be:
// * `?include=` Zero or more `aaronland/go-json-query` query strings containing rules that must match for a document to be considered for further processing.
// * `?exclude=` Zero or more `aaronland/go-json-query`	query strings containing rules that if matched will prevent a document from being considered for further processing.
// * `?include_mode=` A valid `aaronland/go-json-query` query mode string for testing inclusion rules.
// * `?exclude_mode=` A valid `aaronland/go-json-query` query mode string for testing exclusion rules.
func NewGeoJSONLIterator(ctx context.Context, uri string) (Iterator, error) {

	f, err := filters.NewQueryFiltersFromURI(ctx, uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to create filters from query, %w", err)
	}

	idx := &GeojsonLIterator{
		filters: f,
	}

	return idx, nil
}

func (idx *GeojsonLIterator) Iterate(ctx context.Context, uris ...string) iter.Seq2[Record, error] {

	return func(yield func(Record, error) bool) {

		for _, uri := range uris {
			for r, err := range idx.iterate(ctx, uri) {
				yield(r, err)
			}
		}
	}
}

func (idx *GeojsonLIterator) iterate(ctx context.Context, uri string) iter.Seq2[Record, error] {

	return func(yield func(Record, error) bool) {

		fh, err := ReaderWithPath(ctx, uri)

		if err != nil {
			yield(nil, fmt.Errorf("Failed to create reader for '%s', %w", uri, err))
			return
		}

		defer fh.Close()

		// see this - we're using ReadLine because it's entirely possible
		// that the raw GeoJSON (LS) will be too long for bufio.Scanner
		// see also - https://golang.org/pkg/bufio/#Reader.ReadLine
		// (20170822/thisisaaronland)

		reader := bufio.NewReader(fh)
		raw := bytes.NewBuffer([]byte(""))

		i := 0

		for {

			select {
			case <-ctx.Done():
				break
			default:
				// pass
			}

			path := fmt.Sprintf("%s#%d", uri, i)
			i += 1

			fragment, is_prefix, err := reader.ReadLine()

			if err == io.EOF {
				break
			}

			if err != nil {
				yield(nil, fmt.Errorf("Failed to read line at '%s', %w", path, err))
				break
			}

			raw.Write(fragment)

			if is_prefix {
				continue
			}

			br := bytes.NewReader(raw.Bytes())
			fh, err := ioutil.NewReadSeekCloser(br)

			if err != nil {
				yield(nil, fmt.Errorf("Failed to create new ReadSeekCloser for '%s', %w", path, err))
				break
			}

			defer fh.Close()

			if idx.filters != nil {

				ok, err := idx.filters.Apply(ctx, fh)

				if err != nil {
					yield(nil, fmt.Errorf("Failed to apply filters for '%s', %w", path, err))
					break
				}

				if !ok {
					continue
				}

				_, err = fh.Seek(0, 0)

				if err != nil {
					yield(nil, fmt.Errorf("Failed to reset file handle for '%s', %w", path, err))
					break
				}
			}

			iter_r := NewRecord(path, fh)
			yield(iter_r, nil)
		}

		raw.Reset()
	}
}
