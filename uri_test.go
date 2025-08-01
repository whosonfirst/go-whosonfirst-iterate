package iterate

import (
	"net/url"
	"testing"
)

func TestScrubURI(t *testing.T) {

	has_token := []string{
		"example://?access_token=1234",
	}

	for _, uri := range has_token {

		new_uri, err := ScrubURI(uri)

		if err != nil {
			t.Fatalf("Failed to scrub '%s', %v", uri, err)
		}

		u, err := url.Parse(new_uri)

		if err != nil {
			t.Fatalf("Failed to parse new URI '%s' (derived from '%s'), %v", new_uri, uri, err)
		}

		q := u.Query()

		if q.Get("access_token") != "..." {
			t.Fatalf("Invalid access_token query parameter: %s", q.Get("access_token"))
		}
	}
}
