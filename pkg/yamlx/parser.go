package yamlx

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"

	"gopkg.in/yaml.v3"
)

func Load(path string) (*yaml.Node, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var root yaml.Node
	if err := yaml.Unmarshal(data, &root); err != nil {
		return nil, err
	}

	return &root, nil
}

func Save(path string, root *yaml.Node) error {
	out, err := yaml.Marshal(root)
	if err != nil {
		return err
	}

	return os.WriteFile(path, out, 0644)
}

func SaveScalar(path string, node *yaml.Node) error {
	if node == nil {
		return fmt.Errorf("nil node")
	}
	if node.Kind != yaml.ScalarNode {
		return fmt.Errorf("only scalar nodes are supported")
	}
	if node.Style&(yaml.LiteralStyle|yaml.FoldedStyle) != 0 {
		return fmt.Errorf("block scalar values are not supported for in-place updates")
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	start, err := offsetForLineColumn(data, node.Line, node.Column)
	if err != nil {
		return err
	}

	end, err := scalarEndOffset(data, start, node.Style)
	if err != nil {
		return err
	}

	replacement := renderScalar(node.Value, node.Style)
	out := append([]byte{}, data[:start]...)
	out = append(out, replacement...)
	out = append(out, data[end:]...)

	return os.WriteFile(path, out, 0644)
}

func offsetForLineColumn(data []byte, line, column int) (int, error) {
	if line < 1 || column < 1 {
		return 0, fmt.Errorf("invalid node position")
	}

	currentLine := 1
	currentColumn := 1

	for i := 0; i < len(data); {
		if currentLine == line && currentColumn == column {
			return i, nil
		}

		r, size := utf8.DecodeRune(data[i:])
		if r == utf8.RuneError && size == 1 {
			return 0, fmt.Errorf("invalid utf-8 in YAML file")
		}

		i += size
		if r == '\n' {
			currentLine++
			currentColumn = 1
			continue
		}
		currentColumn++
	}

	if currentLine == line && currentColumn == column {
		return len(data), nil
	}

	return 0, fmt.Errorf("node position %d:%d not found in source", line, column)
}

func scalarEndOffset(data []byte, start int, style yaml.Style) (int, error) {
	if start < 0 || start > len(data) {
		return 0, fmt.Errorf("invalid start offset")
	}

	switch {
	case style&yaml.SingleQuotedStyle != 0:
		return singleQuotedEndOffset(data, start)
	case style&yaml.DoubleQuotedStyle != 0:
		return doubleQuotedEndOffset(data, start)
	default:
		return plainScalarEndOffset(data, start), nil
	}
}

func singleQuotedEndOffset(data []byte, start int) (int, error) {
	if start >= len(data) || data[start] != '\'' {
		return 0, fmt.Errorf("single-quoted scalar not found at expected position")
	}
	for i := start + 1; i < len(data); i++ {
		if data[i] != '\'' {
			continue
		}
		if i+1 < len(data) && data[i+1] == '\'' {
			i++
			continue
		}
		return i + 1, nil
	}
	return 0, fmt.Errorf("unterminated single-quoted scalar")
}

func doubleQuotedEndOffset(data []byte, start int) (int, error) {
	if start >= len(data) || data[start] != '"' {
		return 0, fmt.Errorf("double-quoted scalar not found at expected position")
	}
	escaped := false
	for i := start + 1; i < len(data); i++ {
		if escaped {
			escaped = false
			continue
		}
		switch data[i] {
		case '\\':
			escaped = true
		case '"':
			return i + 1, nil
		}
	}
	return 0, fmt.Errorf("unterminated double-quoted scalar")
}

func plainScalarEndOffset(data []byte, start int) int {
	lineEnd := start
	for lineEnd < len(data) && data[lineEnd] != '\n' {
		lineEnd++
	}

	commentStart := -1
	for i := start; i < lineEnd; i++ {
		if data[i] != '#' {
			continue
		}
		if i == start || data[i-1] == ' ' || data[i-1] == '\t' {
			commentStart = i
			break
		}
	}
	if commentStart == -1 {
		commentStart = lineEnd
	}

	end := commentStart
	for end > start && (data[end-1] == ' ' || data[end-1] == '\t') {
		end--
	}
	return end
}

func renderScalar(value string, style yaml.Style) []byte {
	switch {
	case style&yaml.SingleQuotedStyle != 0:
		return []byte("'" + strings.ReplaceAll(value, "'", "''") + "'")
	case style&yaml.DoubleQuotedStyle != 0:
		return []byte(strconv.Quote(value))
	default:
		return []byte(value)
	}
}
