package iterator

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"iter"
	"log/slog"

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

	it := &GeojsonLIterator{
		filters: f,
	}

	return it, nil
}

func (it *GeojsonLIterator) Iterate(ctx context.Context, uris ...string) iter.Seq2[Record, error] {

	return func(yield func(Record, error) bool) {

		for _, uri := range uris {
			for r, err := range it.iterate(ctx, uri) {
				yield(r, err)
			}
		}
	}
}

func (it *GeojsonLIterator) iterate(ctx context.Context, uri string) iter.Seq2[Record, error] {

	logger := slog.Default()
	logger = logger.With("uri", uri)

	return func(yield func(Record, error) bool) {

		r, err := ReaderWithPath(ctx, uri)

		if err != nil {
			logger.Debug("Failed to create reader", "error", err)
			yield(nil, fmt.Errorf("Failed to create reader for '%s', %w", uri, err))
			return
		}

		defer r.Close()

		// see this - we're using ReadLine because it's entirely possible
		// that the raw GeoJSON (LS) will be too long for bufio.Scanner
		// see also - https://golang.org/pkg/bufio/#Reader.ReadLine
		// (20170822/thisisaaronland)

		reader := bufio.NewReader(r)
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

			logger := slog.Default()
			logger = logger.With("uri", uri)
			logger = logger.With("path", path)

			i += 1

			fragment, is_prefix, err := reader.ReadLine()

			if err == io.EOF {
				break
			}

			if err != nil {
				logger.Debug("Failed to readline", "error", err)
				yield(nil, fmt.Errorf("Failed to read line at '%s', %w", path, err))
				break
			}

			raw.Write(fragment)

			if is_prefix {
				logger.Debug("Line is prefix, skipping")
				continue
			}

			br := bytes.NewReader(raw.Bytes())
			ln_r, err := ioutil.NewReadSeekCloser(br)

			if err != nil {
				logger.Debug("Failed to create ReadSeekCloser from line", "error", err)
				yield(nil, fmt.Errorf("Failed to create new ReadSeekCloser for '%s', %w", path, err))
				break
			}

			defer ln_r.Close()

			if it.filters != nil {

				ok, err := it.filters.Apply(ctx, ln_r)

				if err != nil {
					logger.Debug("Failed to apply filters", "error", err)
					yield(nil, fmt.Errorf("Failed to apply filters for '%s', %w", path, err))
					break
				}

				if !ok {
					logger.Debug("No matches after applying filters, skipping")
					continue
				}

				_, err = ln_r.Seek(0, 0)

				if err != nil {
					logger.Debug("Failed to rewind reader", "error", err)
					yield(nil, fmt.Errorf("Failed to reset file handle for '%s', %w", path, err))
					break
				}
			}

			logger.Debug("Yield new record")
			iter_r := NewRecord(path, ln_r)
			yield(iter_r, nil)
		}

		raw.Reset()
	}
}
