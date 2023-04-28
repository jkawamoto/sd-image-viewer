// frontend.go
//
// Copyright (c) 2023 Junpei Kawamoto
//
// This software is released under the MIT License.
//
// http://opensource.org/licenses/mit-license.php

package frontend

import "embed"

//go:embed dist/*
var Contents embed.FS
