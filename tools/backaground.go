package tools

import (
	"s01/background"

	"github.com/ollama/ollama/api"
)

type BackgroundTool struct{}

func (b BackgroundTool) GetTool() api.Tool {
	props := api.NewToolPropertiesMap()
	props.Set("command", api.ToolProperty{
		Type:        api.PropertyType{"string"},
		Description: "shell command",
	})

	backgroundToolFunction := api.ToolFunction{
		Name:        "background_run",
		Description: "Run command in background coroutine. Returns task_id immediately.",
		Parameters: api.ToolFunctionParameters{
			Type:       "object",
			Required:   []string{"command"},
			Properties: props,
		},
	}
	return api.Tool{
		Type:     "function",
		Function: backgroundToolFunction,
	}
}

func (b BackgroundTool) Run(args api.ToolCallFunctionArguments) string {
	commandRaw, _ := args.Get("command")
	commad := commandRaw.(string)
	return background.MyBackgroundManager.Run(commad)
}

type CheckBackgroundTool struct{}

func (c CheckBackgroundTool) GetTool() api.Tool {
	props := api.NewToolPropertiesMap()
	props.Set("task_id", api.ToolProperty{
		Type:        api.PropertyType{"string"},
		Description: "task's id",
	})

	checkBackgroundToolFunction := api.ToolFunction{
		Name:        "check_background",
		Description: "Check background task status. Omit task_id to list all.",
		Parameters: api.ToolFunctionParameters{
			Type:       "object",
			Properties: props,
		},
	}
	return api.Tool{
		Type:     "function",
		Function: checkBackgroundToolFunction,
	}
}

func (c CheckBackgroundTool) Run(args api.ToolCallFunctionArguments) string {
	taskIDRaw, ok := args.Get("task_id")
	var taskID string = ""
	if ok {
		taskID = taskIDRaw.(string)
	}
	return background.MyBackgroundManager.Check(taskID)
}
