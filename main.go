package main

import (
	"context"
	"log"
	"os"

	"github.com/blevesearch/bleve/v2"
	"github.com/jkawamoto/sd-image-viewer/server"
)

func main() {
	ctx := context.Background()
	logger := log.Default()

	if len(os.Args) == 1 {
		logger.Fatalln("one directory path is required")
	}

	bleve.SetLog(logger)

	index, err := newIndex(".bleve")
	if err != nil {
		logger.Fatalf("failed to create an index: %v", err)
	}

	dir := os.Args[1]
	err = indexDir(ctx, dir, index, false, logger)
	if err != nil {
		logger.Fatalf("failed to read files in %v: %v", os.Args[1], err)
	}

	s, err := server.NewServer(index, dir, logger)
	if err != nil {
		logger.Fatalf("failed to create a server: %v", err)
	}

	if err = s.Serve(); err != nil {
		logger.Fatalf("failed to serve: %v", err)
	}
}
