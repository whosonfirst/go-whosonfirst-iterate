package publisher

import (
	"bufio"
	"bytes"
	"context"
	"crypto/sha256"
	"fmt"
	"path/filepath"
	"testing"
)

func TestFeaturePublisher(t *testing.T) {

	ctx := context.Background()

	emitter_uri := "repo://?include=properties.sfomuseum:uri=2019"

	rel_path := "../fixtures"
	abs_path, err := filepath.Abs(rel_path)

	if err != nil {
		t.Fatalf("Failed to derive absolute path for %s, %v", rel_path, err)
	}

	var buf bytes.Buffer
	wr := bufio.NewWriter(&buf)

	fp := &FeaturePublisher{
		AsGeoJSON: true,
		Writer:    wr,
	}

	_, err = fp.Publish(ctx, emitter_uri, abs_path)

	if err != nil {
		t.Fatalf("Failed to publish %s, %v", abs_path, err)
	}

	wr.Flush()

	expected_hash := "30b5663f9d196a3b2164bd429f58023ff31813800da9804f64004fe16a5436e9"

	hash := sha256.Sum256(buf.Bytes())
	str_hash := fmt.Sprintf("%x", hash)

	if str_hash != expected_hash {
		t.Fatalf("Unexpected hash: %s", str_hash)
	}

}
