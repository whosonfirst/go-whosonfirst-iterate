package emitter

import (
	"context"
	"testing"
)

func TestRegisterEmitter(t *testing.T) {

	ctx := context.Background()

	err := RegisterEmitter(ctx, "null", NewNullEmitter)

	if err == nil {
		t.Fatalf("Expected NewNullEmitter to be registered already")
	}
}

func TestNewEmitter(t *testing.T) {

	ctx := context.Background()

	uri := "null://"

	_, err := NewEmitter(ctx, uri)

	if err != nil {
		t.Fatalf("Failed to create new emitter for '%s', %v", uri, err)
	}
}
