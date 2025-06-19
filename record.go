package iterate

import (
	"io"
)

type Record struct {
	Path string
	Body io.ReadSeekCloser
}

func NewRecord(path string, r io.ReadSeekCloser) *Record {

	rec := &Record{
		Path: path,
		Body: r,
	}

	return rec
}
