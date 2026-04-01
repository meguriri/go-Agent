package tools

import (
	"github.com/ollama/ollama/api"
)

type CompactTool struct{}

func (c CompactTool) GetTool() api.Tool {
	props := api.NewToolPropertiesMap()
	props.Set("focus", api.ToolProperty{
		Type:        api.PropertyType{"string"},
		Description: "What to preserve in the summary",
	})

	compactToolFunction := api.ToolFunction{
		Name:        "compact",
		Description: "Trigger manual conversation compression.",
		Parameters: api.ToolFunctionParameters{
			Type:       "object",
			Required:   []string{"focus"},
			Properties: props,
		},
	}
	return api.Tool{
		Type:     "function",
		Function: compactToolFunction,
	}
}

func (c CompactTool) Run(args api.ToolCallFunctionArguments) string {
	return "Manual compression requested."
}
