package filters

import (
	"context"
	"os"
	"testing"
)

func TestNewQueryFiltersFromURI(t *testing.T) {

	ctx := context.Background()

	_, err := NewQueryFiltersFromURI(ctx, "example://?include=properities.mz:is_current=1&exclude=properties.edtf:deprecated=.*")

	if err != nil {
		t.Fatalf("Failed to create new query from URI, %v", err)
	}
}

func TestIncludeQueryFilters(t *testing.T) {

	ctx := context.Background()

	qf_uri := "example://?include=properties.sfomuseum:placetype=map&include=properties.sfomuseum:uri=2019"

	qf, err := NewQueryFiltersFromURI(ctx, qf_uri)

	if err != nil {
		t.Fatalf("Failed to create new query from URI, %v", err)
	}

	path := "../fixtures/data/151/183/838/5/1511838385.geojson"

	r, err := os.Open(path)

	if err != nil {
		t.Fatalf("Failed to open %s, %v", path, err)
	}

	defer r.Close()

	ok, err := qf.Apply(ctx, r)

	if err != nil {
		t.Fatalf("Failed to apply query filters to %s, %v", path, err)
	}

	if !ok {
		t.Fatalf("Expected %s to pass query filters (%s) but did not.", path, qf_uri)
	}
}

func TestExcludeQueryFilters(t *testing.T) {

	ctx := context.Background()

	qf_uri := "example://?include=properties.sfomuseum:placetype=map&exclude=properties.sfomuseum:uri=2019"

	qf, err := NewQueryFiltersFromURI(ctx, qf_uri)

	if err != nil {
		t.Fatalf("Failed to create new query from URI, %v", err)
	}

	path := "../fixtures/data/151/183/838/5/1511838385.geojson"

	r, err := os.Open(path)

	if err != nil {
		t.Fatalf("Failed to open %s, %v", path, err)
	}

	defer r.Close()

	ok, err := qf.Apply(ctx, r)

	if err != nil {
		t.Fatalf("Failed to apply query filters to %s, %v", path, err)
	}

	if ok {
		t.Fatalf("Expected %s to fail query filters (%s) but did not.", path, qf_uri)
	}
}
