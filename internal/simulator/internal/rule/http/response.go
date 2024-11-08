package http

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"gopkg.in/yaml.v3"
)

type Response struct {
	StatusCode int
	Body       []byte
}

func (r *Response) MarshalYAML() (any, error) {
	return &yaml.Node{
		Kind: yaml.MappingNode,
		Content: []*yaml.Node{
			{
				Kind:  yaml.ScalarNode,
				Tag:   "!!str",
				Value: "status",
			},
			{
				Kind:  yaml.ScalarNode,
				Tag:   "!!int",
				Value: strconv.Itoa(r.StatusCode),
			},
			{
				Kind:  yaml.ScalarNode,
				Tag:   "!!str",
				Value: "body",
			},
			{
				Kind:  yaml.ScalarNode,
				Style: yaml.LiteralStyle,
				Tag:   "!!str",
				Value: string(r.Body),
			},
		},
	}, nil
}

func (r *Response) UnmarshalYAML(value *yaml.Node) error {
	if value.Kind != yaml.MappingNode {
		return errors.New("expected a mapping node")
	}

	var codeStr, bodyStr string
	for i := 0; i < len(value.Content); i += 2 {
		key := value.Content[i].Value
		val := value.Content[i+1].Value

		switch key {
		case "status":
			codeStr = val
		case "body":
			bodyStr = val
		default:
			return fmt.Errorf("unexpected key: %s", key)
		}
	}

	code, err := strconv.Atoi(codeStr)
	if err != nil {
		return fmt.Errorf("invalid status code: %d", code)
	}

	r.StatusCode = code
	r.Body = []byte(bodyStr)
	return nil
}

func WriteToFile(path string, response Response) error {
	data, err := yaml.Marshal(&response)
	if err != nil {
		return err
	}

	err = os.WriteFile(path, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func ReadFromFile(path string) (Response, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Response{}, err
	}

	var r Response
	err = yaml.Unmarshal(data, &r)
	if err != nil {
		return Response{}, err
	}

	return r, nil
}
