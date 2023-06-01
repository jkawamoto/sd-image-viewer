// png.go
//
// Copyright (c) 2023 Junpei Kawamoto
//
// This software is released under the MIT License.
//
// http://opensource.org/licenses/mit-license.php

package image

import (
	"fmt"
	"image/png"
	"io"

	"github.com/jkawamoto/go-pngtext"
)

func ParsePNG(r io.ReadSeeker) (*Image, error) {
	cfg, err := png.DecodeConfig(r)
	if err != nil {
		return nil, err
	}

	_, err = r.Seek(0, io.SeekStart)
	if err != nil {
		return nil, err
	}

	list, err := pngtext.ParseTextualData(r)
	if err != nil {
		return nil, err
	}

	data := list.Find("parameters")
	if data == nil {
		return nil, errNoParameters
	}

	params, err := parseParameters(data.Text)
	if err != nil {
		return nil, err
	}

	res := &Image{
		Prompt:         params[promptKey],
		NegativePrompt: params[negativePromptKey],
		Checkpoint:     params[checkpointKey],
		Pixel:          cfg.Height * cfg.Width,
		Metadata:       params,
	}
	res.Metadata[sizeKey] = fmt.Sprintf("%vx%v", cfg.Width, cfg.Height)
	delete(res.Metadata, promptKey)
	delete(res.Metadata, negativePromptKey)
	delete(res.Metadata, checkpointKey)

	return res, nil
}
