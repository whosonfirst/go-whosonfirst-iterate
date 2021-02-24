package ioutil

import (
	"io"
	"log"
	wof_ioutil "github.com/whosonfirst/go-ioutil"
)


func NewReadSeekCloser(fh interface{}) (io.ReadSeekCloser, error) {
	log.Println("The 'whosonfirst/go-reader/ioutil' package is deprecated and will be removed in v2.x. Please use the 'whosonfirst/go-ioutil' package instead.")
	return wof_ioutil.NewReadSeekCloser(fh)
}
