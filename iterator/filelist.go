package iterator

import (
	"bufio"
	"context"
	"fmt"
	"iter"
	"path/filepath"

	"github.com/whosonfirst/go-whosonfirst-iterate/v3/filters"
)

func init() {
	ctx := context.Background()
	RegisterIterator(ctx, "filelist", NewFileListIterator)
}

// FileListIterator implements the `Iterator` interface for crawling records listed in a "file list" (a plain text newline-delimted list of files).
type FileListIterator struct {
	Iterator
	// filters is a `filters.Filters` instance used to include or exclude specific records from being crawled.
	filters filters.Filters
}

// NewFileListIterator() returns a new `FileListIterator` instance configured by 'uri' in the form of:
//
//	file://?{PARAMETERS}
//
// Where {PARAMETERS} may be:
// * `?include=` Zero or more `aaronland/go-json-query` query strings containing rules that must match for a document to be considered for further processing.
// * `?exclude=` Zero or more `aaronland/go-json-query`	query strings containing rules that if matched will prevent a document from being considered for further processing.
// * `?include_mode=` A valid `aaronland/go-json-query` query mode string for testing inclusion rules.
// * `?exclude_mode=` A valid `aaronland/go-json-query` query mode string for testing exclusion rules.
func NewFileListIterator(ctx context.Context, uri string) (Iterator, error) {

	f, err := filters.NewQueryFiltersFromURI(ctx, uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to create filters from query, %w", err)
	}

	idx := &FileListIterator{
		filters: f,
	}

	return idx, nil
}

func (idx *FileListIterator) Iterate(ctx context.Context, uris ...string) iter.Seq2[Record, error] {

	return func(yield func(Record, error) bool) {
		for _, uri := range uris {
			for r, err := range idx.iterate(ctx, uri) {
				yield(r, err)
			}
		}
	}
}

func (idx *FileListIterator) iterate(ctx context.Context, uri string) iter.Seq2[Record, error] {

	return func(yield func(Record, error) bool) {

		abs_path, err := filepath.Abs(uri)

		if err != nil {
			yield(nil, fmt.Errorf("Failed to derive absolute path for '%s', %w", uri, err))
			return
		}

		fh, err := ReaderWithPath(ctx, abs_path)

		if err != nil {
			yield(nil, fmt.Errorf("Failed to create reader for '%s', %w", abs_path, err))
			return
		}

		defer fh.Close()

		scanner := bufio.NewScanner(fh)

		for scanner.Scan() {

			select {
			case <-ctx.Done():
				break
			default:
				// pass
			}

			path := scanner.Text()

			fh, err := ReaderWithPath(ctx, path)

			if err != nil {
				yield(nil, fmt.Errorf("Failed to create reader for '%s', %w", path, err))
				break
			}

			if idx.filters != nil {

				ok, err := idx.filters.Apply(ctx, fh)

				if err != nil {
					yield(nil, fmt.Errorf("Failed to apply filters to '%s', %w", path, err))
					continue
				}

				if !ok {
					continue
				}

				_, err = fh.Seek(0, 0)

				if err != nil {
					yield(nil, fmt.Errorf("Failed to reset file handle for '%s', %w", path, err))
					continue
				}
			}

			iter_r := NewRecord(path, fh)
			yield(iter_r, nil)
		}

		err = scanner.Err()

		if err != nil {
			yield(nil, err)
		}
	}

}
