package main

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/search/query"
	"github.com/jkawamoto/sd-image-viewer/image"
)

const (
	targetExt = ".png"
	posFile   = ".pos"
)

func newIndex(name string) (_ bleve.Index, created bool, err error) {
	index, err := bleve.Open(name)
	if errors.Is(err, bleve.ErrorIndexPathDoesNotExist) {
		indexMapping := bleve.NewIndexMapping()
		indexMapping.AddDocumentMapping(image.DocType, image.DocumentMapping())

		index, err = bleve.New(name, indexMapping)
		created = true
	}
	if err != nil {
		return nil, false, err
	}

	return index, created, nil
}

func indexDir(ctx context.Context, dir string, index bleve.Index, force bool, logger *log.Logger) (err error) {
	var lastIndexed time.Time

	posFileName := filepath.Join(dir, posFile)
	info, err := os.Stat(posFileName)
	if os.IsNotExist(err) {
		// not index yet
	} else if err != nil {
		return err
	} else if !force {
		lastIndexed = info.ModTime()
	}
	defer func() {
		if err == nil {
			f, e := os.Create(posFileName)
			err = errors.Join(err, e, f.Close())
		}
	}()

	return filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if filepath.Ext(path) != targetExt {
			return nil
		}
		info, err := d.Info()
		if err != nil {
			logger.Printf("failed to stat an image file: %v", err)
			return nil
		}

		if info.ModTime().Before(lastIndexed) {
			return nil
		}

		img, err := image.ParseImageFile(path)
		if err != nil {
			logger.Printf("failed to parse an image file: %v", err)
			return nil
		}

		logger.Printf("indexing %v", path)
		err = index.Index(path, img)
		if err != nil {
			logger.Printf("failed to index an image: %v", err)
		}

		return nil
	})
}

func pruneIndex(ctx context.Context, index bleve.Index, logger *log.Logger) error {
	size := 100
	from := 0

	logger.Println("pruning index")
	for {
		if ctx.Err() != nil {
			return ctx.Err()
		}

		res, err := index.Search(bleve.NewSearchRequestOptions(query.NewMatchAllQuery(), size, from, false))
		if err != nil {
			return err
		}
		if len(res.Hits) == 0 {
			logger.Println("finished pruning index")
			return nil
		}
		from += size

		b := index.NewBatch()
		for _, v := range res.Hits {
			_, err = os.Stat(v.ID)
			if os.IsNotExist(err) {
				logger.Printf("removing %v from index", v.ID)
				b.Delete(v.ID)
				from--
			} else if err != nil {
				logger.Printf("failed to stat a file: %v", err)
			}
		}
		if err = index.Batch(b); err != nil {
			return fmt.Errorf("failed to remove items: %w", err)
		}
	}
}
