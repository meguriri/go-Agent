package tools

import (
	"fmt"
	"os"
	"strings"

	"github.com/ollama/ollama/api"
)

type EditTool struct{}

func (e EditTool) GetTool() api.Tool {
	props := api.NewToolPropertiesMap()
	props.Set("path", api.ToolProperty{
		Type:        api.PropertyType{"string"},
		Description: "file path",
	})
	props.Set("old_text", api.ToolProperty{
		Type:        api.PropertyType{"string"},
		Description: "need to replace old_text",
	})
	props.Set("new_text", api.ToolProperty{
		Type:        api.PropertyType{"string"},
		Description: "edit to new_text",
	})

	editToolFunction := api.ToolFunction{
		Name:        "edit_file",
		Description: "Replace exact text in file.",
		Parameters: api.ToolFunctionParameters{
			Type:       "object",
			Required:   []string{"path", "old_text", "new_text"},
			Properties: props,
		},
	}
	return api.Tool{
		Type:     "function",
		Function: editToolFunction,
	}
}

func (e EditTool) Run(args api.ToolCallFunctionArguments) string {
	pathRaw, _ := args.Get("path")
	oldTextRaw, _ := args.Get("old_text")
	newTextRaw, _ := args.Get("new_text")
	path := pathRaw.(string)
	oldText := oldTextRaw.(string)
	newText := newTextRaw.(string)
	resolved, err := safePath(path)
	if err != nil {
		return fmt.Sprintf("Error: %v", err)
	}
	data, err := os.ReadFile(resolved)
	if err != nil {
		return fmt.Sprintf("Error: %v", err)
	}
	if !strings.Contains(string(data), oldText) {
		return fmt.Sprintf("Error: Text not found in %s", path)
	}
	newContent := strings.Replace(string(data), oldText, newText, 1)
	err = os.WriteFile(path, []byte(newContent), 0644)
	if err != nil {
		return fmt.Sprintf("Error: %v", err)
	}
	return fmt.Sprintf("Edited %s", path)
}
