package tools

import (
	"github.com/ollama/ollama/api"
)

type ReadInboxTool struct {
	Reader InboxReader
}

func (r ReadInboxTool) GetTool() api.Tool {
	props := api.NewToolPropertiesMap()
	props.Set("sender", api.ToolProperty{
		Type:        api.PropertyType{"string"},
		Description: "which agent send this message",
	})
	readToolFunction := api.ToolFunction{
		Name:        "read_inbox",
		Description: "Read and drain your inbox.",
		Parameters: api.ToolFunctionParameters{
			Type:       "object",
			Required:   []string{"sender"},
			Properties: props,
		},
	}
	return api.Tool{
		Type:     "function",
		Function: readToolFunction,
	}
}

func (r ReadInboxTool) Run(args api.ToolCallFunctionArguments) string {
	senderRaw, _ := args.Get("sender")
	sender := senderRaw.(string)
	message := r.Reader.ReadInboxText(sender)
	return message

}
