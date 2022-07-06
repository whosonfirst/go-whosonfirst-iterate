package publisher

import (
	"context"
	"testing"
)

type TestPublisher struct {
	Publisher
}

func (tp *TestPublisher) Publish(ctx context.Context, emitter_uri string, uris ...string) (int64, error) {
	return 0, nil
}

func TestPublisherInterface(t *testing.T) {

	ctx := context.Background()

	tp := &TestPublisher{}

	var p interface{} = tp
	_, ok := p.(Publisher)

	if !ok {
		t.Fatalf("Invalid interface")
	}

	_, err := tp.Publish(ctx, "repo://", "example")

	if err != nil {
		t.Fatalf("Failed to publish results, %v", err)
	}
}
