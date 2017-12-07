package interpolater

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"

	"github.com/Masterminds/sprig"
	yamlConverter "github.com/ghodss/yaml"
	yaml "gopkg.in/yaml.v2"
)

const (
	FormatPreserve = "preserve"
	FormatJSON     = "json"
	FormatYAML     = "yaml"
)

type Interpolater struct {
	Writer       io.Writer
	OutputFormat string
}

func (i Interpolater) Execute(basePath string, inputPaths []string) error {
	baseContents, err := ioutil.ReadFile(basePath)
	if err != nil {
		return fmt.Errorf("unable to read template file at '%s': %s", basePath, err)
	}

	inputVariables, err := i.readInputVars(inputPaths)
	if err != nil {
		return err
	}

	t, err := template.New("template").
		Funcs(sprig.FuncMap()).
		Option("missingkey=error").
		Parse(string(baseContents))
	if err != nil {
		return fmt.Errorf("template '%s' is not valid text/template format: %s", basePath, err)
	}

	var buffer bytes.Buffer
	err = t.Execute(&buffer, inputVariables)
	if err != nil {
		return fmt.Errorf("failed to render template '%s': %s", basePath, err)
	}

	return i.writeOutput(buffer.Bytes(), basePath)
}

func (i Interpolater) readInputVars(inputPaths []string) (map[string]interface{}, error) {
	inputVariables := map[string]interface{}{}
	for _, file := range inputPaths {
		fileContents, err := ioutil.ReadFile(file)
		if err != nil {
			return nil, fmt.Errorf("unable to input file at '%s': %s", file, err)
		}

		var fileVars map[string]interface{}
		err = yaml.Unmarshal(fileContents, &fileVars)
		if err != nil {
			return nil, fmt.Errorf("expected '%s' to be a valid YAML file: %s", file, err)
		}

		for k, v := range fileVars {
			inputVariables[k] = v
		}
	}
	return inputVariables, nil
}

func (i Interpolater) writeOutput(interpolatedContents []byte, basePath string) error {
	var output []byte
	var err error

	switch format := i.OutputFormat; format {
	case FormatPreserve:
		output = interpolatedContents
	case FormatJSON:
		output, err = yamlConverter.YAMLToJSON(interpolatedContents)
		if err != nil {
			return fmt.Errorf("template '%s' is not valid YAML/JSON: %s", basePath, err)
		}
	case FormatYAML:
		var v interface{}
		err = yaml.Unmarshal(interpolatedContents, &v)
		if err != nil {
			return fmt.Errorf("template '%s' is not valid YAML/JSON: %s", basePath, err)
		}
		output, err = yaml.Marshal(v)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported output type '%s'", format)
	}

	_, err = i.Writer.Write(output)
	if err != nil {
		return err
	}

	return nil
}
