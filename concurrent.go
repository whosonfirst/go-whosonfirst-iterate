// Package iterator provides methods and utilities for iterating over a collection of records
// (presumed but not required to be Who's On First records) from a variety of sources and dispatching
// processing to user-defined callback functions.
package iterate

import (
	"context"
	"fmt"
	_ "io"
	"iter"
	"log/slog"
	"net/url"
	"regexp"
	"runtime"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/whosonfirst/go-whosonfirst-uri"
)

type ConcurrentIterator struct {
	Iterator
	iterator Iterator
	// seen is the count of documents that have been seen (or emitted).
	seen int64
	// ...
	iterating *atomic.Bool
	// max_procs is the number maximum (CPU) processes to used to process documents simultaneously.
	max_procs int
	// exclude_paths is a `regexp.Regexp` instance used to test and exclude (if matching) the paths of documents as they are iterated through.
	exclude_paths     *regexp.Regexp
	exclude_alt_files bool
	// ...
	include_paths *regexp.Regexp
	max_attempts  int
	retry_after   int
	// skip records (specifically their relative URI) that have already been processed
	dedupe bool
	// lookup table to track records (specifically their relative URI) that have been processed
	dedupe_map *sync.Map
}

// NewIterator() returns a new `Iterator` instance derived from 'emitter_uri' and 'emitter_cb'. The former is expected
// to be a valid `whosonfirst/go-whosonfirst-iterate/v2/emitter.Emitter` URI whose semantics are defined by the underlying
// implementation of the `emitter.Emitter` interface. The following iterator-specific query parameters are also accepted:
// * `?_max_procs=` Explicitly set the number maximum processes to use for iterating documents simultaneously. (Default is the value of `runtime.NumCPU()`.)
// * `?_exclude=` A valid regular expresion used to test and exclude (if matching) the paths of documents as they are iterated through.
// * `?_dedupe=` A boolean value to track and skip records (specifically their relative URI) that have already been processed.
func NewConcurrentIterator(ctx context.Context, iterator_uri string, it Iterator) (Iterator, error) {

	u, err := url.Parse(iterator_uri)

	if err != nil {
		return nil, err
	}

	q := u.Query()

	max_procs := runtime.NumCPU()

	retry := false
	max_attempts := 1
	retry_after := 10 // seconds

	if q.Has("_max_procs") {

		max, err := strconv.ParseInt(q.Get("_max_procs"), 10, 64)

		if err != nil {
			return nil, fmt.Errorf("Failed to parse '_max_procs' parameter, %w", err)
		}

		max_procs = int(max)
	}

	if q.Has("_retry") {

		v, err := strconv.ParseBool(q.Get("_retry"))

		if err != nil {
			return nil, fmt.Errorf("Failed to parse '_retry' parameter, %w", err)
		}

		retry = v
	}

	if retry {

		if q.Has("_max_retries") {

			v, err := strconv.Atoi(q.Get("_max_retries"))

			if err != nil {
				return nil, fmt.Errorf("Failed to parse '_max_retries' parameter, %w", err)
			}

			max_attempts = v
		}

		if q.Has("_retry_after") {

			v, err := strconv.Atoi(q.Get("_retry_after"))

			if err != nil {
				return nil, fmt.Errorf("Failed to parse '_retry_after' parameter, %w", err)
			}

			retry_after = v
		}
	}

	i := &ConcurrentIterator{
		iterator:     it,
		seen:         int64(0),
		iterating:    new(atomic.Bool),
		max_procs:    max_procs,
		max_attempts: max_attempts,
		retry_after:  retry_after,
	}

	if q.Has("_include") {

		re_include, err := regexp.Compile(q.Get("_include"))

		if err != nil {
			return nil, fmt.Errorf("Failed to parse '_include' parameter, %w", err)
		}

		i.include_paths = re_include
	}

	if q.Has("_exclude") {

		re_exclude, err := regexp.Compile(q.Get("_exclude"))

		if err != nil {
			return nil, fmt.Errorf("Failed to parse '_exclude' parameter, %w", err)
		}

		i.exclude_paths = re_exclude
	}

	if q.Has("_exclude_alt") {

		v, err := strconv.ParseBool(q.Get("_exclude_alt"))

		if err != nil {
			return nil, fmt.Errorf("Failed to parse '_exclude_alt' parameter, %w", err)
		}

		i.exclude_alt_files = v
	}

	if q.Has("_dedupe") {

		v, err := strconv.ParseBool(q.Get("_dedupe"))

		if err != nil {
			return nil, fmt.Errorf("Failed to parse '_dedupe' parameter, %w", err)
		}

		if v {
			i.dedupe = true
			i.dedupe_map = new(sync.Map)
		}

	}

	return i, nil
}

func (it *ConcurrentIterator) Iterate(ctx context.Context, uris ...string) iter.Seq2[*Record, error] {

	return func(yield func(rec *Record, err error) bool) {

		t1 := time.Now()

		defer func() {
			slog.Debug("Time to process paths", "count", len(uris), "time", time.Since(t1))
		}()

		it.iterating.Swap(true)
		defer it.iterating.Swap(false)

		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		procs := it.max_procs
		throttle := make(chan bool, procs)

		for i := 0; i < procs; i++ {
			throttle <- true
		}

		done_ch := make(chan bool)
		err_ch := make(chan error)
		rec_ch := make(chan *Record)

		remaining := len(uris)

		for _, uri := range uris {

			go func(uri string) {

				t2 := time.Now()
				defer func() {
					slog.Debug("Time to iterate uri", "uri", uri, "time", time.Since(t2))
				}()

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

				for rec, err := range it.iterator.Iterate(ctx, uri) {

					if err != nil {
						err_ch <- err
						continue
					}

					atomic.AddInt64(&it.seen, 1)

					// Notes about automatically closing rec.Body
					// It would be nice to be able to just use the common
					// defer rec.Body.Close() here but the mechanics of
					// the way yield functions works means we need to do
					// after the yield function. This happens below in the
					// unsurprisingly named `do_yield` function. What all
					// of this means is that we need to be more attentive
					// than usual about closing filehandles, when necessary,
					// before the *Record instance is yield-ed. I guess this
					// is just the price of using iterators for the time
					// being. And yes, I did try using runtime.AddCleanup
					// but because it execute as part of the runtime.GC process
					// it often gets triggered after the *Record instance
					// has been purged without closing the underlying file
					// handle. Basically what we need is a Python-style object
					// level destructor but those don't exist yet so, again,
					// here we are.

					ok, err := it.shouldYieldRecord(ctx, rec)

					if err != nil {
						err_ch <- err
						continue
					}

					if !ok {
						rec.Body.Close()
						continue
					}

					rec_ch <- rec
				}

			}(uri)
		}

		do_yield := func(rec *Record, err error) bool {

			// This bit is important. See notes above.

			if rec != nil {
				defer func() {
					slog.Debug("Close record", "path", rec.Path)
					rec.Body.Close()
				}()
			}

			return yield(rec, err)
		}

		for remaining > 0 {
			select {
			case <-done_ch:
				remaining -= 1
			case err := <-err_ch:
				do_yield(nil, err)
				return
			case rec := <-rec_ch:
				if !do_yield(rec, nil) {
					return
				}
			default:
				// pass
			}
		}

	}
}

// Seen() returns the total number of records processed so far.
func (it ConcurrentIterator) Seen() int64 {
	return atomic.LoadInt64(&it.seen)
}

// IsIterating() returns a boolean value indicating whether 'it' is still processing documents.
func (it ConcurrentIterator) IsIterating() bool {
	return it.iterating.Load()
}

func (it ConcurrentIterator) shouldYieldRecord(ctx context.Context, rec *Record) (bool, error) {

	if it.include_paths != nil {

		if !it.include_paths.MatchString(rec.Path) {
			return false, nil
		}
	}

	if it.exclude_paths != nil {

		if it.exclude_paths.MatchString(rec.Path) {
			return false, nil
		}
	}

	if it.exclude_alt_files {

		is_alt, err := uri.IsAltFile(rec.Path)

		if err != nil {
			return false, err
		}

		if is_alt {
			return false, nil
		}
	}

	if it.dedupe {

		id, uri_args, err := uri.ParseURI(rec.Path)

		if err != nil {
			return false, fmt.Errorf("Failed to parse %s, %w", rec.Path, err)
		}

		rel_path, err := uri.Id2RelPath(id, uri_args)

		if err != nil {
			return false, fmt.Errorf("Failed to derive relative path for %s, %w", rec.Path, err)
		}

		_, seen := it.dedupe_map.LoadOrStore(rel_path, true)

		if seen {
			slog.Debug("Skip record", "path", rel_path)
			return false, nil
		}
	}

	return true, nil
}
