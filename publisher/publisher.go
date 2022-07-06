// Package publisher provides interfaces for (re)publishing documents that are emitted by an `Iterator` instance.
package publisher

import (
	"context"
)

// type Publisher provides an interface for (re)publishing documents that are emitted by an `Iterator` instance.
type Publisher interface {
	// Publish() publishes documents that are emitted by an `Iterator` instance. It takes as its arguments
	// a valid `emitter.Emitter` URI and a list of URIs to iterate through. It is assumed that implementations
	// of the Publisher interface will provide their own internal `emitter.EmitterCallbackFunc` callback
	// functions.
	Publish(context.Context, string, ...string) (int64, error)
}
