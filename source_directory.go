package iterate

import (
	"context"
	"fmt"
	"io"
	"iter"
	"os"
	"path/filepath"
	_ "log/slog"
	"sync"
	
	"github.com/whosonfirst/go-whosonfirst-crawl"
	"github.com/whosonfirst/go-whosonfirst-iterate/v2/filters"
)

func init() {
	// ctx := context.Background()
	// RegisterSource(ctx, "directory", NewDirectorySource)
}

// DirectorySource implements the `Source` interface for crawling records in a directory.
type DirectorySource struct {
	Source
	// filters is a `filters.Filters` instance used to include or exclude specific records from being crawled.
	filters filters.Filters
}

// NewDirectorySource() returns a new `DirectorySource` instance configured by 'uri' in the form of:
//
//	directory://?{PARAMETERS}
//
// Where {PARAMETERS} may be:
// * `?include=` Zero or more `aaronland/go-json-query` query strings containing rules that must match for a document to be considered for further processing.
// * `?exclude=` Zero or more `aaronland/go-json-query`	query strings containing rules that if matched will prevent a document from being considered for further processing.
// * `?include_mode=` A valid `aaronland/go-json-query` query mode string for testing inclusion rules.
// * `?exclude_mode=` A valid `aaronland/go-json-query` query mode string for testing exclusion rules.
func NewDirectorySource(ctx context.Context, uri string) (Source, error) {

	f, err := filters.NewQueryFiltersFromURI(ctx, uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to create filters from query, %w", err)
	}

	idx := &DirectorySource{
		filters: f,
	}

	return idx, nil
}

func (idx *DirectorySource) Walk(ctx context.Context, uri string) iter.Seq2[*Record, error] {

	return func(yield func(rec *Record, err error) bool) {

		abs_path, err := filepath.Abs(uri)

		if err != nil {
			yield(nil, fmt.Errorf("Failed to derive absolute path for '%s', %w", uri, err))
		}

		mu := new(sync.RWMutex)
		
		crawl_cb := func(path string, info os.FileInfo) error {

			select {
			case <-ctx.Done():
				return nil
			default:
				// pass
			}

			if info.IsDir() {
				return nil
			}

			r, err := ReaderWithPath(ctx, path)

			if err != nil {
				return fmt.Errorf("Failed to create reader for '%s', %w", abs_path, err)
			}

			defer r.Close()

			if idx.filters != nil {

				ok, err := idx.filters.Apply(ctx, r)

				if err != nil {
					return fmt.Errorf("Failed to apply filters for '%s', %w", abs_path, err)
				}

				if !ok {
					return nil
				}

				_, err = r.Seek(0, 0)

				if err != nil {
					return fmt.Errorf("Failed to seek(0, 0) on reader for '%s', %w", abs_path, err)
				}
			}

			rec := &Record{
				Path: path,
				Body: r,
			}

			mu.Lock()
			defer mu.Unlock()
			
			if !yield(rec, nil){
			   return io.EOF
			}
			
			return nil
		}

		c := crawl.NewCrawler(abs_path)
		err = c.Crawl(crawl_cb)

		if err != nil && err != io.EOF {
			yield(nil, err)
		}
	}
}
