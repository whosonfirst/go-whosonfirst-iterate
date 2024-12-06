package iterator

import (
	"fmt"
	"io"
)

type RecordX interface {
	URI() string
	Body() io.ReadSeeker
}

type Record struct {
	URI  string
	Body io.ReadSeeker
}

func (r *Record) ReadAll() ([]byte, error) {

	_, err := r.Body.Seek(0, 0)

	if err != nil {
		return nil, fmt.Errorf("Failed to rewind body, %w", err)
	}

	return io.ReadAll(r.Body)
}
