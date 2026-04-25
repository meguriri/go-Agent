package tools

import (
	"github.com/ollama/ollama/api"
)

type ReadInboxTool struct {
	Reader InboxReader
}

func (r ReadInboxTool) GetTool() api.Tool {
	props := api.NewToolPropertiesMap()
	readToolFunction := api.ToolFunction{
		Name:        "read_inbox",
		Description: "Read and drain the lead's inbox.",
		Parameters: api.ToolFunctionParameters{
			Type:       "object",
			Properties: props,
		},
	}
	return api.Tool{
		Type:     "function",
		Function: readToolFunction,
	}
}

func (r ReadInboxTool) Run(args api.ToolCallFunctionArguments) string {
	message := r.Reader.ReadInboxText("lead")
	return message

}
