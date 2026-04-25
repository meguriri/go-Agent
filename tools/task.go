package tools

import (
	"context"

	"github.com/ollama/ollama/api"
)

type TaskTool struct {
	Runner SubagentRunner
}

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
		Description: "生成一个具有全新上下文的subagent。它共享文件系统，但不共享对话历史。",
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
	promptRaw, _ := args.Get("prompt")
	descriptionRaw, ok := args.Get("description")
	var description string
	if ok {
		description = descriptionRaw.(string)
	}
	prompt := promptRaw.(string)
	ctx := context.Background()
	message := t.Runner.RunSubagent(ctx, prompt, description)
	return message
}
