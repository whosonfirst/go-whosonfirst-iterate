// Package iterator provides methods and utilities for iterating over a collection of records
// (presumed but not required to be Who's On First records) from a variety of sources and dispatching
// processing to user-defined callback functions.
package iterate

import (
	"context"
	"fmt"
	"iter"
	"log/slog"
	"net/url"
	"regexp"
	"runtime"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/whosonfirst/go-whosonfirst-iterate/v3/iterator"
)

// type Iterator provides a struct that can be used for iterating over a collection of records
// (presumed but not required to be Who's On First records) from a variety of sources and dispatching
// processing to user-defined callback functions.
type Iterator struct {
	// Seen is the count of documents that have been seen (or emitted).
	Seen int64
	// count is the current number of documents being processed used to signal where an `Iterator` instance is still indexing (processing) documents.
	count int64
	// max_procs is the number maximum (CPU) processes to used to process documents simultaneously.
	max_procs int
	// exclude_paths is a `regexp.Regexp` instance used to test and exclude (if matching) the paths of documents as they are iterated through.
	exclude_paths *regexp.Regexp
}

func NewIterator(ctx context.Context) (*Iterator, error) {
	return NewIteratorFromURI(ctx, "iter://")
}

func NewIteratorFromURI(ctx context.Context, uri string) (*Iterator, error) {

	max_procs := runtime.NumCPU()

	i := Iterator{
		Seen:      0,
		count:     0,
		max_procs: max_procs,
	}

	u, err := url.Parse(uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse URI, %w", err)
	}

	q := u.Query()

	if q.Has("max_procs") {

		max, err := strconv.ParseInt(q.Get("max_procs"), 10, 64)

		if err != nil {
			return nil, fmt.Errorf("Failed to parse 'max_procs' parameter, %w", err)
		}

		max_procs = int(max)
	}

	if q.Has("exclude") {

		re_exclude, err := regexp.Compile(q.Get("exclude"))

		if err != nil {
			return nil, fmt.Errorf("Failed to parse 'exclude' parameter, %w", err)
		}

		i.exclude_paths = re_exclude
	}

	return &i, nil
}

// IterateURIs processes 'uris' concurrent dispatching each URI to the iterator's underlying `Emitter.WalkURI`
// method and `EmitterCallbackFunc` callback function.
func (idx *Iterator) Iterate(ctx context.Context, provider_uri string, provider_sources ...string) iter.Seq2[*iterator.Record, error] {

	return func(yield func(*iterator.Record, error) bool) {

		it, err := iterator.NewIterator(ctx, provider_uri)

		if err != nil {
			yield(nil, err)
			return
		}

		logger := slog.Default()

		t1 := time.Now()

		defer func() {
			logger.Debug("time to index paths", "count", len(provider_sources), "time", time.Since(t1))
		}()

		idx.increment()
		defer idx.decrement()

		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		procs := idx.max_procs
		throttle := make(chan bool, procs)

		for i := 0; i < procs; i++ {
			throttle <- true
		}

		done_ch := make(chan bool)
		err_ch := make(chan error)

		remaining := len(provider_sources)

		for _, uri := range provider_sources {

			go func(uri string) {

				logger := slog.Default()
				logger = logger.With("uri", uri)

				<-throttle

				defer func() {
					throttle <- true
					done_ch <- true
				}()

				select {
				case <-ctx.Done():
					return
				default:
					// pass
				}

				for r, err := range it.Iterate(ctx, uri) {

					go atomic.AddInt64(&idx.Seen, 1)

					if idx.exclude_paths != nil {

						if idx.exclude_paths.MatchString(r.URI) {
							continue
						}
					}

					yield(r, err)
				}

			}(uri)
		}

		for remaining > 0 {
			select {
			case <-done_ch:
				remaining -= 1
			case err := <-err_ch:
				logger.Error(err.Error())
			default:
				// pass
			}
		}
	}

}

// IsIndexing() returns a boolean value indicating whether 'idx' is still processing documents.
func (idx *Iterator) IsIndexing() bool {

	if atomic.LoadInt64(&idx.count) > 0 {
		return true
	}

	return false
}

// increment() increments the count of documents being processed.
func (idx *Iterator) increment() {
	atomic.AddInt64(&idx.count, 1)
}

// decrement() decrements the count of documents being processed.
func (idx *Iterator) decrement() {
	atomic.AddInt64(&idx.count, -1)
}
