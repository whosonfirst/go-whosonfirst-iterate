package main

import (
	"context"
	"log"

	"github.com/whosonfirst/go-whosonfirst-iterate/v3/app/emit"
)

func main() {

	ctx := context.Background()
	err := emit.Run(ctx)

	if err != nil {
		log.Fatalf("Failed to emit records, %v", err)
	}
}
