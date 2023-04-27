package image

import (
	"errors"
	"fmt"
	"image/png"
	"io"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/analysis/analyzer/simple"
	"github.com/blevesearch/bleve/v2/analysis/analyzer/standard"
	"github.com/blevesearch/bleve/v2/mapping"
	"github.com/jkawamoto/go-pngtext"
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

	img, err := Parse(f)
	if err != nil {
		return nil, err
	}

	img.CreationTime = info.ModTime()
	return img, nil
}

func Parse(r io.ReadSeeker) (*Image, error) {
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

var parametersRegexp = regexp.MustCompile("(.*?)(?:Negative prompt: (.+))?Steps: (\\d+), ")

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
		quoted bool
		key    string
	)
	for i := 0; i != len(additionalParameters); i++ {
		switch additionalParameters[i] {
		case '"':
			quoted = !quoted
		case ':':
			if !quoted {
				key = strings.Trim(additionalParameters[pos:i], " ")
				pos = i + 1
			}
		case ',':
			if !quoted {
				res[key] = strings.Trim(additionalParameters[pos:i], " ")
				pos = i + 1
			}
		}
	}
	res[key] = strings.Trim(additionalParameters[pos:], " ")

	return res, nil
}

func DocumentMapping() *mapping.DocumentMapping {
	textFieldMapping := bleve.NewTextFieldMapping()
	textFieldMapping.Analyzer = standard.Name

	keywordFieldMapping := bleve.NewTextFieldMapping()
	keywordFieldMapping.Analyzer = simple.Name

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
