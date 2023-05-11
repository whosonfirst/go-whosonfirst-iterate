package publisher

import (
	"context"
	"fmt"
	"io"
	"sync"
	"sync/atomic"

	"github.com/whosonfirst/go-whosonfirst-iterate/v2/iterator"
)

// FeaturePublisher implements the Publisher interface for (re)publishing GeoJSON Feature documents that are emitted by an `Iterator` instance.
type FeaturePublisher struct {
	Publisher
	// AsJSON is a boolean flag signaling that the final output should be published as a JSON array.
	AsJSON bool
	// AsGeoJSON is a boolean flag signaling that the final output should be published as a GeoJSON FeatureCollection.
	AsGeoJSON bool
	// Writer is the underlying `io.Writer` instance where published data will be written to.
	Writer io.Writer
}

// Publish() will (re)publish all the documents emitted from an `Iterator` instance derived from 'emitter_uri' and 'uris'.
func (pub *FeaturePublisher) Publish(ctx context.Context, emitter_uri string, uris ...string) (int64, error) {

	mu := new(sync.RWMutex)

	var count int64
	var count_bytes int64

	count = 0
	count_bytes = 0

	if pub.AsGeoJSON {

		b, err := pub.Writer.Write([]byte(`{"type":"FeatureCollection", "features":`))

		if err != nil {
			return atomic.LoadInt64(&count_bytes), fmt.Errorf("Failed to write GeoJSON header, %w", err)
		}

		atomic.AddInt64(&count_bytes, int64(b))
	}

	if pub.AsGeoJSON || pub.AsJSON {

		b, err := pub.Writer.Write([]byte(`[`))

		if err != nil {
			return atomic.LoadInt64(&count_bytes), fmt.Errorf("Failed to write JSON array header, %w", err)
		}

		atomic.AddInt64(&count_bytes, int64(b))
	}

	emitter_cb := func(ctx context.Context, path string, fh io.ReadSeeker, args ...interface{}) error {

		select {
		case <-ctx.Done():
			return nil
		default:
			// pass
		}

		mu.Lock()
		defer mu.Unlock()

		atomic.AddInt64(&count, 1)

		if pub.AsGeoJSON || pub.AsJSON {
			if atomic.LoadInt64(&count) > 1 {

				b, err := pub.Writer.Write([]byte(`,`))

				if err != nil {
					return fmt.Errorf("Failed to write JSON array separator, %w", err)
				}

				atomic.AddInt64(&count_bytes, int64(b))
			}
		}

		b, err := io.Copy(pub.Writer, fh)

		if err != nil {
			return fmt.Errorf("Failed to copy data from %s, %w", path, err)
		}

		atomic.AddInt64(&count_bytes, int64(b))
		return nil
	}

	iter, err := iterator.NewIterator(ctx, emitter_uri, emitter_cb)

	if err != nil {
		return atomic.LoadInt64(&count_bytes), fmt.Errorf("Failed to create new iterator, %w", err)
	}

	err = iter.IterateURIs(ctx, uris...)

	if err != nil {
		return atomic.LoadInt64(&count_bytes), fmt.Errorf("Failed to iterate URIs, %w", err)
	}

	if pub.AsGeoJSON || pub.AsJSON {

		b, err := pub.Writer.Write([]byte(`]`))

		if err != nil {
			return atomic.LoadInt64(&count_bytes), fmt.Errorf("Failed to close JSON array, %w", err)
		}

		atomic.AddInt64(&count_bytes, int64(b))
	}

	if pub.AsGeoJSON {

		b, err := pub.Writer.Write([]byte(`}`))

		if err != nil {
			return atomic.LoadInt64(&count_bytes), fmt.Errorf("Failed to close GeoJSON FeatureCollection, %w", err)
		}

		atomic.AddInt64(&count_bytes, int64(b))
	}

	return atomic.LoadInt64(&count_bytes), nil
}
