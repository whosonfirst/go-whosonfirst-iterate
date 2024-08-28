package iterator

import (
	"context"
	"fmt"
	"io/fs"
	"iter"

	"github.com/whosonfirst/go-ioutil"
	"github.com/whosonfirst/go-whosonfirst-iterate/v3/filters"
)

// DirectoryIterator implements the `Iterator` interface for crawling records in a directory.
type FSIterator struct {
	Iterator
	options *FSIteratorOptions
}

type FSIteratorOptions struct {
	Filters filters.Filters
	FS      fs.FS
}

func NewFSIteratorWithOptions(ctx context.Context, opts *FSIteratorOptions) (Iterator, error) {

	idx := &FSIterator{
		options: opts,
	}

	return idx, nil
}

func (idx *FSIterator) Walk(ctx context.Context, uris ...string) iter.Seq2[*Candidate, error] {

	return func(yield func(*Candidate, error) bool) {

		walk_func := func(path string, d fs.DirEntry, err error) error {

			select {
			case <-ctx.Done():
				return nil
			default:
				// pass
			}

			if err != nil {
				return err
			}

			if d.IsDir() {
				return nil
			}

			// logger := slog.Default()
			// logger = logger.With("path", path")

			r, err := idx.options.FS.Open(path)

			if err != nil {
				return fmt.Errorf("Failed to create reader for '%s', %w", path, err)
			}

			rsc, err := ioutil.NewReadSeekCloser(r)

			if err != nil {
				return fmt.Errorf("Failed to create new ReadSeekCloser for %s, %w", path, err)
			}

			defer rsc.Close()

			if idx.options.Filters != nil {

				ok, err := idx.options.Filters.Apply(ctx, rsc)

				if err != nil {
					return fmt.Errorf("Failed to apply filters for '%s', %w", path, err)
				}

				if !ok {
					return nil
				}

				_, err = rsc.Seek(0, 0)

				if err != nil {
					return fmt.Errorf("Failed to seek(0, 0) on reader for '%s', %w", path, err)
				}
			}

			c := &Candidate{
				Path:   path,
				Reader: rsc,
			}

			yield(c, nil)
			return nil
		}

		err := fs.WalkDir(idx.options.FS, ".", walk_func)

		if err != nil {
			yield(nil, err)
		}
	}
}
