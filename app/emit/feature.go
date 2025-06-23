package emit

import (
	"context"
	"fmt"
	"io"
	// "sync"
	"sync/atomic"

	"github.com/whosonfirst/go-whosonfirst-iterate/v3"
)

// FeatureEmitter implements the Emitter interface for (re)publishing GeoJSON Feature documents that are emitted by an `Iterator` instance.
type FeatureEmitter struct {
	Emitter
	// AsJSON is a boolean flag signaling that the final output should be published as a JSON array.
	AsJSON bool
	// AsGeoJSON is a boolean flag signaling that the final output should be published as a GeoJSON FeatureCollection.
	AsGeoJSON bool
	// Writer is the underlying `io.Writer` instance where published data will be written to.
	Writer io.Writer
}

// Emit() will (re)publish all the documents emitted from an `Iterator` instance derived from 'iterator_uri' and 'uris'.
func (pub *FeatureEmitter) Emit(ctx context.Context, iterator_uri string, uris ...string) (int64, error) {

	var count int64
	var count_bytes int64

	count = 0
	count_bytes = 0

	it, err := iterate.NewIterator(ctx, iterator_uri)

	if err != nil {
		return atomic.LoadInt64(&count_bytes), fmt.Errorf("Failed to create new iterator, %w", err)
	}

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

	for rec, err := range it.Iterate(ctx, uris...) {

		select {
		case <-ctx.Done():
			return atomic.LoadInt64(&count_bytes), nil
		default:
			// pass
		}

		if err != nil {
			return atomic.LoadInt64(&count_bytes), err
		}

		atomic.AddInt64(&count, 1)

		if pub.AsGeoJSON || pub.AsJSON {
			if atomic.LoadInt64(&count) > 1 {

				b, err := pub.Writer.Write([]byte(`,`))

				if err != nil {
					return atomic.LoadInt64(&count_bytes), fmt.Errorf("Failed to write JSON array separator, %w", err)
				}

				atomic.AddInt64(&count_bytes, int64(b))
			}
		}

		b, err := io.Copy(pub.Writer, rec.Body)

		if err != nil {
			return atomic.LoadInt64(&count_bytes), fmt.Errorf("Failed to copy data from %s, %w", rec.Path, err)
		}

		atomic.AddInt64(&count_bytes, int64(b))
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
