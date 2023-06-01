// webp.go
//
// Copyright (c) 2023 Junpei Kawamoto
//
// This software is released under the MIT License.
//
// http://opensource.org/licenses/mit-license.php

package image

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/gohugoio/hugo/resources/images/exif"
	"golang.org/x/image/riff"
)

func ParseWebP(r io.Reader) (*Image, error) {
	formType, rr, err := riff.NewReader(r)
	if err != nil {
		return nil, err
	}
	if string(formType[:]) != "WEBP" {
		return nil, errors.New("not webp file")
	}

	var (
		width, height uint32
		text          string
	)
	for {
		id, chunkLen, chunkReader, err := rr.Next()
		if errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			return nil, err
		}

		switch string(id[:]) {
		case "VP8X":
			data, err := io.ReadAll(chunkReader)
			if err != nil {
				return nil, err
			}

			width = binary.LittleEndian.Uint32(append(data[chunkLen-6:chunkLen-3], 0)) + 1
			height = binary.LittleEndian.Uint32(append(data[chunkLen-3:], 0)) + 1

		case "EXIF":
			decoder, err := exif.NewDecoder()
			if err != nil {
				return nil, err
			}

			ex, err := decoder.Decode(chunkReader)
			if err != nil {
				return nil, err
			}
			if tag, ok := ex.Tags["UserComment"]; ok {
				text = strings.TrimPrefix(fmt.Sprint(tag), "UNICODE")
			}
		}
	}

	params, err := parseParameters(text)
	if err != nil {
		return nil, err
	}

	res := &Image{
		Prompt:         params[promptKey],
		NegativePrompt: params[negativePromptKey],
		Checkpoint:     params[checkpointKey],
		Pixel:          int(width * height),
		Metadata:       params,
	}
	res.Metadata[sizeKey] = fmt.Sprintf("%vx%v", width, height)
	delete(res.Metadata, promptKey)
	delete(res.Metadata, negativePromptKey)
	delete(res.Metadata, checkpointKey)

	return res, nil
}
