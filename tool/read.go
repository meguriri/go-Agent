package tool

import (
	"fmt"
	"os"
	"strings"

	"github.com/ollama/ollama/api"
)

func GetReadTool() api.Tool {
	props := api.NewToolPropertiesMap()
	props.Set("path", api.ToolProperty{
		Type:        api.PropertyType{"string"},
		Description: "file path",
	})
	props.Set("limit", api.ToolProperty{
		Type:        api.PropertyType{"integer"},
		Description: "content limit",
	})

	readToolFunction := api.ToolFunction{
		Name:        "read_file",
		Description: "Read file contents.",
		Parameters: api.ToolFunctionParameters{
			Type:       "object",
			Required:   []string{"path"},
			Properties: props,
		},
	}
	return api.Tool{
		Type:     "function",
		Function: readToolFunction,
	}
}

func RunRead(path string, limit int) string {
	resolved, err := safePath(path)
	if err != nil {
		return fmt.Sprintf("Error: %v", err)
	}
	data, err := os.ReadFile(resolved)
	if err != nil {
		return fmt.Sprintf("Error: %v", err)
	}
	text := strings.ReplaceAll(string(data), "\r\n", "\n")
	lines := strings.Split(text, "\n")
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	if limit > 0 && limit < len(lines) {
		more := len(lines) - limit
		lines = append(lines[:limit], fmt.Sprintf("... (%d more lines)", more))
	}
	out := strings.Join(lines, "\n")
	runes := []rune(out)
	if len(runes) > 50000 {
		return string(runes[:50000])
	}
	return out
}
