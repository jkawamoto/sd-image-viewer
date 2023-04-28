package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/blevesearch/bleve/v2"

	"github.com/jkawamoto/sd-image-viewer/server"
)

const AppName = "sd-image-viewer"
const Version = "v0.0.0"

func main() {
	logger := log.Default()
	bleve.SetLog(logger)

	cacheDir, err := os.UserCacheDir()
	if err != nil {
		logger.Printf("Failed to get the cache directory: %v", err)
	}

	host := flag.String("host", "localhost", "the IP to listen on")
	port := flag.Int("port", 0, "the port to listen on for insecure connections, defaults to a random value")
	indexPath := flag.String("index", filepath.Join(cacheDir, AppName), "path to the index")
	force := flag.Bool("force", false, "force reindexing all images")
	prune := flag.Bool("prune", false, "remove non exiting images from the index")
	version := flag.Bool("v", false, "prints current version")

	flag.Parse()
	if *version {
		fmt.Printf("%v %v\n", AppName, Version)
		os.Exit(0)
	}
	if flag.NArg() == 0 {
		logger.Fatalln("One directory path is required")
	}
	dir := flag.Arg(0)

	index, created, err := newIndex(*indexPath)
	if err != nil {
		logger.Fatalf("Failed to create an index: %v", err)
	}
	if created {
		// if a new index is created, force reindexing all images.
		*force = true
	}
	defer func() {
		logger.Printf("Closing the index")
		if err = index.Close(); err != nil {
			logger.Printf("Failed to close the index: %v", err)
		}
	}()

	var wg sync.WaitGroup
	defer wg.Wait()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if *prune {
			err := pruneIndex(ctx, index, logger)
			if errors.Is(err, context.Canceled) {
				return
			} else if err != nil {
				logger.Printf("Failed to prune index: %v", err)
			}
		}
		for {
			err := indexDir(ctx, dir, index, *force, logger)
			if errors.Is(err, context.Canceled) {
				break
			} else if err != nil {
				logger.Printf("Failed to index files in %v: %v", dir, err)
			}
			*force = false
			select {
			case <-ctx.Done():
				break
			case <-time.After(1 * time.Hour):
			}
		}
	}()

	s, err := server.NewServer(*host, *port, index, dir, logger)
	if err != nil {
		logger.Fatalf("Failed to create a server: %v", err)
	}

	if err = s.Serve(); err != nil {
		logger.Fatalf("Failed to serve: %v", err)
	}
}
