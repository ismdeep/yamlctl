package yamlx

import (
	"fmt"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

func Get(root *yaml.Node, path string) (string, error) {
	segments, err := parsePath(path)
	if err != nil {
		return "", err
	}
	node, err := getNode(root, segments)
	if err != nil {
		return "", err
	}
	if node.Kind != yaml.ScalarNode {
		return "", fmt.Errorf("value at path is not a scalar")
	}
	return node.Value, nil
}

func Set(root *yaml.Node, path string, value string) (*yaml.Node, error) {
	segments, err := parsePath(path)
	if err != nil {
		return nil, err
	}
	return setNode(root, segments, value)
}

type pathSegment struct {
	key   *string
	index *int
}

func parsePath(path string) ([]pathSegment, error) {
	if path == "" {
		return nil, fmt.Errorf("invalid path")
	}

	parts := strings.Split(path, ".")
	segments := make([]pathSegment, 0, len(parts))

	for _, part := range parts {
		if part == "" {
			return nil, fmt.Errorf("invalid path")
		}

		for len(part) > 0 {
			bracket := strings.IndexByte(part, '[')
			if bracket == -1 {
				segments = append(segments, pathSegment{key: stringPtr(part)})
				break
			}

			if bracket > 0 {
				segments = append(segments, pathSegment{key: stringPtr(part[:bracket])})
			}

			end := strings.IndexByte(part[bracket:], ']')
			if end == -1 {
				return nil, fmt.Errorf("invalid path")
			}

			rawIndex := part[bracket+1 : bracket+end]
			if rawIndex == "" {
				return nil, fmt.Errorf("invalid path")
			}

			index, err := strconv.Atoi(rawIndex)
			if err != nil || index < 0 {
				return nil, fmt.Errorf("invalid array index: %s", rawIndex)
			}
			segments = append(segments, pathSegment{index: intPtr(index)})
			part = part[bracket+end+1:]
		}
	}

	return segments, nil
}

func getNode(node *yaml.Node, segments []pathSegment) (*yaml.Node, error) {
	if node.Kind == yaml.DocumentNode {
		if len(node.Content) == 0 {
			return nil, fmt.Errorf("empty document")
		}
		return getNode(node.Content[0], segments)
	}

	if len(segments) == 0 {
		return node, nil
	}

	segment := segments[0]

	if segment.key != nil {
		if node.Kind != yaml.MappingNode {
			return nil, fmt.Errorf("invalid path")
		}

		for i := 0; i < len(node.Content); i += 2 {
			k := node.Content[i]
			v := node.Content[i+1]

			if k.Value == *segment.key {
				return getNode(v, segments[1:])
			}
		}

		return nil, fmt.Errorf("key not found: %s", *segment.key)
	}

	if node.Kind != yaml.SequenceNode {
		return nil, fmt.Errorf("invalid path")
	}

	if *segment.index >= len(node.Content) {
		return nil, fmt.Errorf("index out of range: %d", *segment.index)
	}

	return getNode(node.Content[*segment.index], segments[1:])
}

func setNode(node *yaml.Node, segments []pathSegment, value string) (*yaml.Node, error) {
	if node.Kind == yaml.DocumentNode {
		return setNode(node.Content[0], segments, value)
	}

	if len(segments) == 0 {
		return nil, nil
	}

	segment := segments[0]

	if segment.key != nil {
		if node.Kind != yaml.MappingNode {
			return nil, fmt.Errorf("invalid path")
		}

		for i := 0; i < len(node.Content); i += 2 {
			k := node.Content[i]
			v := node.Content[i+1]

			if k.Value == *segment.key {
				if len(segments) == 1 {
					if v.Kind != yaml.ScalarNode {
						return nil, fmt.Errorf("only scalar values can be updated without rewriting formatting")
					}
					v.Value = value
					return v, nil
				}
				return setNode(v, segments[1:], value)
			}
		}

		return nil, fmt.Errorf("key not found: %s", *segment.key)
	}

	if node.Kind != yaml.SequenceNode {
		return nil, fmt.Errorf("invalid path")
	}

	if *segment.index >= len(node.Content) {
		return nil, fmt.Errorf("index out of range: %d", *segment.index)
	}

	if len(segments) == 1 {
		target := node.Content[*segment.index]
		if target.Kind != yaml.ScalarNode {
			return nil, fmt.Errorf("only scalar values can be updated without rewriting formatting")
		}
		target.Value = value
		return target, nil
	}

	return setNode(node.Content[*segment.index], segments[1:], value)
}

func stringPtr(value string) *string {
	return &value
}

func intPtr(value int) *int {
	return &value
}
