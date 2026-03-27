package tools

import (
	"github.com/ollama/ollama/api"
)

type TaskTool struct{}

func (t TaskTool) GetTool() api.Tool {
	props := api.NewToolPropertiesMap()
	props.Set("prompt", api.ToolProperty{
		Type:        api.PropertyType{"string"},
		Description: "prompt for the subagent",
	})
	props.Set("description", api.ToolProperty{
		Type:        api.PropertyType{"string"},
		Description: "Short description of the task",
	})

	taskToolFunction := api.ToolFunction{
		Name:        "task",
		Description: "Spawn a subagent with fresh context. It shares the filesystem but not conversation history.",
		Parameters: api.ToolFunctionParameters{
			Type:       "object",
			Required:   []string{"prompt"},
			Properties: props,
		},
	}
	return api.Tool{
		Type:     "function",
		Function: taskToolFunction,
	}
}

func (t TaskTool) Run(args api.ToolCallFunctionArguments) string {
	//run_subagent
	return ""
}
