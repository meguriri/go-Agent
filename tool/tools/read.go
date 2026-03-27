package tools

import (
	"fmt"
	"os"
	"strings"

	"github.com/ollama/ollama/api"
)

type ReadTool struct{}

func (r ReadTool) GetTool() api.Tool {
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

func (r ReadTool) Run(args api.ToolCallFunctionArguments) string {
	pathRaw, _ := args.Get("path")
	limitRaw, ok := args.Get("limit")
	path := pathRaw.(string)
	var limit int = 0
	if ok {
		switch v := limitRaw.(type) {
		case int:
			limit = v
		case int8:
			limit = int(v)
		case int16:
			limit = int(v)
		case int32:
			limit = int(v)
		case int64:
			limit = int(v)
		case uint:
			limit = int(v)
		case uint8:
			limit = int(v)
		case uint16:
			limit = int(v)
		case uint32:
			limit = int(v)
		case uint64:
			limit = int(v)
		case float32:
			limit = int(v)
		case float64:
			limit = int(v)
		}
	}
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
