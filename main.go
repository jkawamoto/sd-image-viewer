package main

import (
	"context"
	"flag"
	"log"
	"sync"

	"github.com/blevesearch/bleve/v2"
	"github.com/jkawamoto/sd-image-viewer/server"
)

func main() {
	logger := log.Default()
	bleve.SetLog(logger)

	indexPath := flag.String("index", ".bleve", "path to the index")
	force := flag.Bool("force", false, "force reindexing all images")

	flag.Parse()
	if flag.NArg() == 0 {
		logger.Fatalln("one directory path is required")
	}
	dir := flag.Arg(0)

	index, err := newIndex(*indexPath)
	if err != nil {
		logger.Fatalf("failed to create an index: %v", err)
	}

	var wg sync.WaitGroup
	defer wg.Wait()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wg.Add(1)
	go func() {
		defer wg.Done()

		if err := indexDir(ctx, dir, index, *force, logger); err != nil {
			logger.Printf("failed to index files in %v: %v", dir, err)
		}
	}()

	s, err := server.NewServer(index, dir, logger)
	if err != nil {
		logger.Fatalf("failed to create a server: %v", err)
	}

	if err = s.Serve(); err != nil {
		logger.Fatalf("failed to serve: %v", err)
	}
}
