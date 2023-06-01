// image.go
//
// Copyright (c) 2023 Junpei Kawamoto
//
// This software is released under the MIT License.
//
// http://opensource.org/licenses/mit-license.php

package image

import (
	"errors"
	"os"
	"path/filepath"
	"time"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/analysis/analyzer/simple"
	"github.com/blevesearch/bleve/v2/analysis/analyzer/standard"
	"github.com/blevesearch/bleve/v2/mapping"
)

const (
	DocType = "Image"

	promptKey         = "Prompt"
	negativePromptKey = "Negative Prompt"
	checkpointKey     = "Model"
	stepsKey          = "Steps"
	sizeKey           = "Size"
)

var (
	errNoParameters           = errors.New("no parameters found")
	errNotSupportedParameters = errors.New("parameter format is not supported")
)

type Image struct {
	Prompt         string            `json:"prompt"`
	NegativePrompt string            `json:"negative-prompt"`
	Checkpoint     string            `json:"checkpoint"`
	Pixel          int               `json:"pixel"`
	CreationTime   time.Time         `json:"creation-time"`
	Metadata       map[string]string `json:"metadata"`
}

func (*Image) Type() string {
	return DocType
}

func Analyzer(field string) string {
	switch field {
	case "prompt", "negative-prompt":
		return standard.Name
	case "checkpoint":
		return simple.Name
	default:
		return ""
	}
}

func ParseImageFile(name string) (_ *Image, err error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = errors.Join(err, f.Close())
	}()

	info, err := f.Stat()
	if err != nil {
		return nil, err
	}

	var img *Image
	switch filepath.Ext(name) {
	case ".png":
		img, err = ParsePNG(f)
		if err != nil {
			return nil, err
		}
	case ".webp":
		img, err = ParseWebP(f)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("filetype not supported")
	}

	img.CreationTime = info.ModTime()
	return img, nil
}

func DocumentMapping() *mapping.DocumentMapping {
	textFieldMapping := bleve.NewTextFieldMapping()
	textFieldMapping.Analyzer = standard.Name

	keywordFieldMapping := bleve.NewKeywordFieldMapping()

	intFieldMapping := bleve.NewNumericFieldMapping()

	dateTimeFieldMapping := bleve.NewDateTimeFieldMapping()

	docMapping := bleve.NewDocumentMapping()
	docMapping.AddFieldMappingsAt("prompt", textFieldMapping)
	docMapping.AddFieldMappingsAt("negative-prompt", textFieldMapping)
	docMapping.AddFieldMappingsAt("checkpoint", keywordFieldMapping)
	docMapping.AddFieldMappingsAt("pixel", intFieldMapping)
	docMapping.AddFieldMappingsAt("creation-time", dateTimeFieldMapping)
	docMapping.AddSubDocumentMapping("metadata", bleve.NewDocumentMapping())

	return docMapping
}
