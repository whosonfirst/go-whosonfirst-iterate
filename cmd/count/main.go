package main

import (
	"context"
	"log"

	"github.com/whosonfirst/go-whosonfirst-iterate/v3/app/count"
)

func main() {

	ctx := context.Background()
	err := count.Run(ctx)

	if err != nil {
		log.Fatalf("Failed to count records, %v", err)
	}
}
