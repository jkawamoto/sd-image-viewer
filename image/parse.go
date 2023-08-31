// params.go
//
// Copyright (c) 2023 Junpei Kawamoto
//
// This software is released under the MIT License.
//
// http://opensource.org/licenses/mit-license.php

package image

import (
	"fmt"
	"regexp"
	"strings"
)

var parametersRegexp = regexp.MustCompile(`(.*?)\s*(?:Negative prompt: (.+?)\s*)?Steps: (\d+), `)

func parseParameters(text string) (map[string]string, error) {
	text = strings.ReplaceAll(text, "\n", " ")
	m := parametersRegexp.FindAllStringSubmatch(text, -1)
	if len(m) != 1 {
		return nil, fmt.Errorf("%w: %v", errNotSupportedParameters, text)
	}

	res := make(map[string]string)

	sm := m[0]
	res[promptKey] = sm[1]
	res[negativePromptKey] = sm[2]
	res[stepsKey] = sm[3]

	additionalParameters := text[len(sm[0]):]
	var (
		pos    int
		json   bool
		quoted bool
		key    string
	)
	for i := 0; i != len(additionalParameters); i++ {
		switch additionalParameters[i] {
		case '"':
			if !json {
				quoted = !quoted
			}
		case ':':
			if !quoted && !json {
				key = strings.Trim(additionalParameters[pos:i], " ")
				pos = i + 1
			}
		case ',':
			if !quoted && !json {
				res[key] = strings.Trim(additionalParameters[pos:i], " ")
				pos = i + 1
				key = ""
			}
		case '{':
			json = true
		case '}':
			json = false
		}
	}
	if key != "" {
		res[key] = strings.Trim(additionalParameters[pos:], " ")
	}

	return res, nil
}
