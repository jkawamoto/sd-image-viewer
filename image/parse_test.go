// image_test.go
//
// Copyright (c) 2023 Junpei Kawamoto
//
// This software is released under the MIT License.
//
// http://opensource.org/licenses/mit-license.php

package image

import (
	"errors"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
)

func Test_parseParameters(t *testing.T) {
	prompt := gofakeit.Paragraph(10, 1, 3, ", ")
	negativePrompt := gofakeit.Paragraph(10, 1, 3, ", ")
	steps := gofakeit.IntRange(1, 100)
	checkpoint := gofakeit.AppName()

	cases := []struct {
		text   string
		expect map[string]string
		err    error
	}{
		{
			text: fmt.Sprintf(
				"%v Negative prompt: %v Steps: %v, Model: %v", prompt, negativePrompt, steps, checkpoint),
			expect: map[string]string{
				promptKey:         prompt,
				negativePromptKey: negativePrompt,
				stepsKey:          fmt.Sprint(steps),
				checkpointKey:     checkpoint,
			},
		},
		{
			text: fmt.Sprintf(
				"%v,Negative prompt: %v Steps: %v, Model: %v", prompt, negativePrompt, steps, checkpoint),
			expect: map[string]string{
				promptKey:         prompt + ",",
				negativePromptKey: negativePrompt,
				stepsKey:          fmt.Sprint(steps),
				checkpointKey:     checkpoint,
			},
		},
		{
			text: fmt.Sprintf(
				"%v.Negative prompt: %v Steps: %v, Model: %v", prompt, negativePrompt, steps, checkpoint),
			expect: map[string]string{
				promptKey:         prompt + ".",
				negativePromptKey: negativePrompt,
				stepsKey:          fmt.Sprint(steps),
				checkpointKey:     checkpoint,
			},
		},
		{
			text: fmt.Sprintf(
				"%vNegative prompt: %v Steps: %v, Model: %v", prompt, negativePrompt, steps, checkpoint),
			expect: map[string]string{
				promptKey:         prompt,
				negativePromptKey: negativePrompt,
				stepsKey:          fmt.Sprint(steps),
				checkpointKey:     checkpoint,
			},
		},
		{
			text: fmt.Sprintf(
				"%v Negative prompt: %v,Steps: %v, Model: %v", prompt, negativePrompt, steps, checkpoint),
			expect: map[string]string{
				promptKey:         prompt,
				negativePromptKey: negativePrompt + ",",
				stepsKey:          fmt.Sprint(steps),
				checkpointKey:     checkpoint,
			},
		},
		{
			text: fmt.Sprintf(
				"%v Negative prompt: %v.Steps: %v, Model: %v", prompt, negativePrompt, steps, checkpoint),
			expect: map[string]string{
				promptKey:         prompt,
				negativePromptKey: negativePrompt + ".",
				stepsKey:          fmt.Sprint(steps),
				checkpointKey:     checkpoint,
			},
		},
		{
			text: fmt.Sprintf(
				"%v Negative prompt: %vSteps: %v, Model: %v", prompt, negativePrompt, steps, checkpoint),
			expect: map[string]string{
				promptKey:         prompt,
				negativePromptKey: negativePrompt,
				stepsKey:          fmt.Sprint(steps),
				checkpointKey:     checkpoint,
			},
		},
		{
			text: fmt.Sprintf(
				"%v Negative prompt: %v Steps: %v, Model: %v, ", prompt, negativePrompt, steps, checkpoint),
			expect: map[string]string{
				promptKey:         prompt,
				negativePromptKey: negativePrompt,
				stepsKey:          fmt.Sprint(steps),
				checkpointKey:     checkpoint,
			},
		},
		{
			text: fmt.Sprintf("%v Steps: %v, Model: %v", prompt, steps, checkpoint),
			expect: map[string]string{
				promptKey:         prompt,
				negativePromptKey: "",
				stepsKey:          fmt.Sprint(steps),
				checkpointKey:     checkpoint,
			},
		},
		{
			text: fmt.Sprintf(
				"%v Negative prompt: %v Steps: %v, Model: %v, Param1: abc, Param2: 123 456",
				prompt, negativePrompt, steps, checkpoint,
			),
			expect: map[string]string{
				promptKey:         prompt,
				negativePromptKey: negativePrompt,
				stepsKey:          fmt.Sprint(steps),
				checkpointKey:     checkpoint,
				"Param1":          "abc",
				"Param2":          "123 456",
			},
		},
		{
			text: fmt.Sprintf(
				`%v Negative prompt: %v Steps: %v, Model: %v, AdditionalParams: "p1: abc, p2: 456", Param1: abc`,
				prompt, negativePrompt, steps, checkpoint,
			),
			expect: map[string]string{
				promptKey:          prompt,
				negativePromptKey:  negativePrompt,
				stepsKey:           fmt.Sprint(steps),
				checkpointKey:      checkpoint,
				"AdditionalParams": `"p1: abc, p2: 456"`,
				"Param1":           "abc",
			},
		},
		{
			text: fmt.Sprintf(
				`%v Negative prompt: %v Steps: %v, Model: %v, Hashes: {"vae": "abc", "embed:123": "456"}, Param1: abc`,
				prompt, negativePrompt, steps, checkpoint,
			),
			expect: map[string]string{
				promptKey:         prompt,
				negativePromptKey: negativePrompt,
				stepsKey:          fmt.Sprint(steps),
				checkpointKey:     checkpoint,
				"Hashes":          `{"vae": "abc", "embed:123": "456"}`,
				"Param1":          "abc",
			},
		},
		{
			text: prompt,
			err:  errNotSupportedParameters,
		},
	}
	for _, c := range cases {
		t.Run(c.text, func(t *testing.T) {
			res, err := parseParameters(c.text)
			if !errors.Is(err, c.err) {
				t.Errorf("expect %v, got %v", c.expect, err)
			}
			if len(res) != len(c.expect) {
				t.Errorf("expect %v, got %v", c.expect, res)
			}
			for k, v := range c.expect {
				if res[k] != v {
					t.Errorf("expect %q, got %q", v, res[k])
				}
			}
		})
	}
}
