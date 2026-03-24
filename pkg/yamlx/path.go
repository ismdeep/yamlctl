package yamlx

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

func Get(root *yaml.Node, path string) (string, error) {
	node, err := getNode(root, strings.Split(path, "."))
	if err != nil {
		return "", err
	}
	if node.Kind != yaml.ScalarNode {
		return "", fmt.Errorf("value at path is not a scalar")
	}
	return node.Value, nil
}

func Set(root *yaml.Node, path string, value string) (*yaml.Node, error) {
	keys := strings.Split(path, ".")
	return setNode(root, keys, value)
}

func getNode(node *yaml.Node, keys []string) (*yaml.Node, error) {
	if node.Kind == yaml.DocumentNode {
		if len(node.Content) == 0 {
			return nil, fmt.Errorf("empty document")
		}
		return getNode(node.Content[0], keys)
	}

	if len(keys) == 0 {
		return node, nil
	}

	if node.Kind != yaml.MappingNode {
		return nil, fmt.Errorf("invalid path")
	}

	for i := 0; i < len(node.Content); i += 2 {
		k := node.Content[i]
		v := node.Content[i+1]

		if k.Value == keys[0] {
			return getNode(v, keys[1:])
		}
	}

	return nil, fmt.Errorf("key not found: %s", keys[0])
}

func setNode(node *yaml.Node, keys []string, value string) (*yaml.Node, error) {
	if node.Kind == yaml.DocumentNode {
		return setNode(node.Content[0], keys, value)
	}

	if len(keys) == 0 {
		return nil, nil
	}

	if node.Kind != yaml.MappingNode {
		return nil, fmt.Errorf("invalid path")
	}

	for i := 0; i < len(node.Content); i += 2 {
		k := node.Content[i]
		v := node.Content[i+1]

		if k.Value == keys[0] {
			if len(keys) == 1 {
				if v.Kind != yaml.ScalarNode {
					return nil, fmt.Errorf("only scalar values can be updated without rewriting formatting")
				}
				v.Value = value
				return v, nil
			}
			return setNode(v, keys[1:], value)
		}
	}

	return nil, fmt.Errorf("key not found: %s", keys[0])
}
