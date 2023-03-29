package emitter

import (
	"context"
	"fmt"
	"io/fs"

	"github.com/whosonfirst/go-whosonfirst-crawl"
	"github.com/whosonfirst/go-whosonfirst-iterate/v2/filters"
)

// FsEmitter implements the `Emitter` interface for crawling records in a fs.
type FsEmitter struct {
	Emitter
	// filters is a `filters.Filters` instance used to include or exclude specific records from being crawled.
	filters filters.Filters
	fs      fs.FS
}

func NewFsEmitter(ctx context.Context, uri string, iterator_fs fs.FS) (Emitter, error) {

	f, err := filters.NewQueryFiltersFromURI(ctx, uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to create filters from query, %w", err)
	}

	idx := &FsEmitter{
		fs:      iterator_fs,
		filters: f,
	}

	return idx, nil
}

// WalkURI() walks (crawls) the fs named 'uri' and for each file (not excluded by any filters specified
// when `idx` was created) invokes 'index_cb'.
func (idx *FsEmitter) WalkURI(ctx context.Context, index_cb EmitterCallbackFunc, uri string) error {

	var walk_func func(path string, d fs.DirEntry, err error) error

	walk_func = func(path string, d fs.DirEntry, err error) error {

		if err != nil {
			return fmt.Errorf("Failed to walk %s, %w", path, err)
		}

		if d.IsDir() {

			if path == "." {
				return nil
			}

			return fs.WalkDir(idx.fs, path, walk_func)
		}

		r, err := idx.fs.Open(path)

		if err != nil {
			return fmt.Errorf("Failed to open %s for reading, %w", path, err)
		}

		defer r.Close()

		if idx.filters != nil {

			ok, err := idx.filters.Apply(ctx, r)

			if err != nil {
				return fmt.Errorf("Failed to apply filters for '%s', %w", path, err)
			}

			if !ok {
				return nil
			}

			_, err = r.Seek(0, 0)

			if err != nil {
				return fmt.Errorf("Failed to seek(0, 0) on reader for '%s', %w", path, err)
			}
		}

		err = index_cb(ctx, path, r)

		if err != nil {
			return fmt.Errorf("Failed to invoke callback for '%s', %w", path, err)
		}

		return nil
	}

	err := fs.WalkDir(idx.fs, uri, walk_func)

	if err != nil {
		return fmt.Errorf("Failed to walk filesystem, %w", err)
	}

	return nil
}
