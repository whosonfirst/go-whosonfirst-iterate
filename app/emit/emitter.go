package emit

import (
	"context"
)

// type Emitter provides an interface for (re)publishing documents that are emitted by an `Iterator` instance.
type Emitter interface {
	// Emits() writes documents that are emitted by an `Iterator` instance. It takes as its arguments
	// a valid `iterate.Iterator` URI and a list of URIs to iterate through.
	Emit(context.Context, string, ...string) (int64, error)
}
