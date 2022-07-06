package filters

import (
	"bytes"
	"context"
	"io"
	"testing"
)

type TestFilters struct {
	Filters
}

func (tf *TestFilters) Apply(ctx context.Context, r io.ReadSeeker) (bool, error) {
	return true, nil
}

func TestFiltersInterface(t *testing.T) {

	ctx := context.Background()

	tf := &TestFilters{}

	var p interface{} = tf
	_, ok := p.(Filters)

	if !ok {
		t.Fatalf("Invalid interface")
	}

	r := bytes.NewReader([]byte("hello world"))

	ok, err := tf.Apply(ctx, r)

	if err != nil {
		t.Fatalf("Failed to apply filters, %v", err)
	}

	if !ok {
		t.Fatalf("Expected filters to pass")
	}
}
