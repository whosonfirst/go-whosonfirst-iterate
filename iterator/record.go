package iterator

import (
	"io"
)

type Record interface {
	URI() string
	Body() io.ReadSeeker
}

type IteratorRecord struct {
	Record
	uri  string
	body io.ReadSeeker
}

func (r *IteratorRecord) URI() string {
	return r.uri
}

func (r *IteratorRecord) Body() io.ReadSeeker {
	return r.body
}

func NewRecord(uri string, body io.ReadSeeker) Record {

	r := &IteratorRecord{
		uri:  uri,
		body: body,
	}

	return r
}
