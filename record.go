package iterate

import (
	"io"
)

type Record struct {
	Path string
	Body io.ReadSeeker
}
