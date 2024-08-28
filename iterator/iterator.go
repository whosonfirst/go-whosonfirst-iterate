package iterator

import (
	"context"
	"fmt"
	"io"
	"iter"
	"net/url"
	"sort"
	"strings"

	"github.com/aaronland/go-roster"
)

type Candidate struct {
	Path   string
	Reader io.ReadSeeker
}

type Iterator interface {
	Walk(context.Context, ...string) iter.Seq2[*Candidate, error]
}

// type IteratorInitializeFunc is a function used to initialize an implementation of the `Iterator` interface.
type IteratorInitializeFunc func(context.Context, string) (Iterator, error)

// iterators is a `aaronland/go-roster.Roster` instance used to maintain a list of registered `IteratorInitializeFunc` initialization functions.
var iterators roster.Roster

// RegisterIterator() associates 'scheme' with 'init_func' in an internal list of avilable `Iterator` implementations.
func RegisterIterator(ctx context.Context, scheme string, f IteratorInitializeFunc) error {

	err := ensureSpatialRoster()

	if err != nil {
		return fmt.Errorf("Failed to register %s scheme, %w", scheme, err)
	}

	return iterators.Register(ctx, scheme, f)
}

// NewIterator() returns a new `Iterator` instance derived from 'uri'. The semantics of and requirements for
// 'uri' as specific to the package implementing the interface.
func NewIterator(ctx context.Context, uri string) (Iterator, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse URI, %w", err)
	}

	scheme := u.Scheme

	if scheme == "" {
		return nil, fmt.Errorf("Emittter URI is missing scheme '%s'", uri)
	}

	i, err := iterators.Driver(ctx, scheme)

	if err != nil {
		return nil, fmt.Errorf("Failed to retrieve driver for '%s' scheme, %w", scheme, err)
	}

	fn := i.(IteratorInitializeFunc)

	if fn == nil {
		return nil, fmt.Errorf("Unregistered initialization function for '%s' scheme", scheme)
	}

	return fn(ctx, uri)
}

// Schemes() returns the list of schemes that have been "registered".
func Schemes() []string {

	ctx := context.Background()
	schemes := []string{}

	err := ensureSpatialRoster()

	if err != nil {
		return schemes
	}

	for _, dr := range iterators.Drivers(ctx) {
		scheme := fmt.Sprintf("%s://", strings.ToLower(dr))
		schemes = append(schemes, scheme)
	}

	sort.Strings(schemes)
	return schemes
}

// ensureDispatcherRoster() ensures that a `aaronland/go-roster.Roster` instance used to maintain a list of registered `IteratorInitializeFunc`
// initialization functions is present
func ensureSpatialRoster() error {

	if iterators == nil {

		r, err := roster.NewDefaultRoster()

		if err != nil {
			return fmt.Errorf("Failed to create new roster, %w", err)
		}

		iterators = r
	}

	return nil
}
