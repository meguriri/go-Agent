package tools

import (
	"fmt"
	"os"

	"github.com/ollama/ollama/api"
)

type WriteTool struct{}

func (w WriteTool) GetTool() api.Tool {
	props := api.NewToolPropertiesMap()
	props.Set("path", api.ToolProperty{
		Type:        api.PropertyType{"string"},
		Description: "file path",
	})
	props.Set("content", api.ToolProperty{
		Type:        api.PropertyType{"string"},
		Description: "writed content",
	})

	writeToolFunction := api.ToolFunction{
		Name:        "write_file",
		Description: "write content to file.",
		Parameters: api.ToolFunctionParameters{
			Type:       "object",
			Required:   []string{"path", "content"},
			Properties: props,
		},
	}
	return api.Tool{
		Type:     "function",
		Function: writeToolFunction,
	}
}

func (w WriteTool) Run(args api.ToolCallFunctionArguments) string {
	pathRaw, _ := args.Get("path")
	contentRaw, _ := args.Get("content")
	path := pathRaw.(string)
	content := contentRaw.(string)
	resolved, err := safePath(path)
	if err != nil {
		return fmt.Sprintf("Error: %v", err)
	}
	err = os.WriteFile(resolved, []byte(content), 0644)
	if err != nil {
		return fmt.Sprintf("Error: %v", err)
	}

	return fmt.Sprintf("Wrote %d bytes to %s", len(content), path)
}
