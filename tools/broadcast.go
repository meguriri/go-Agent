package tools

import (
	"github.com/ollama/ollama/api"
)

type BroadcastTool struct {
	TeamManager TeamManager
	Broadcaster Broadcaster
}

func (b BroadcastTool) GetTool() api.Tool {
	props := api.NewToolPropertiesMap()
	props.Set("content", api.ToolProperty{
		Type:        api.PropertyType{"string"},
		Description: "broadcasts message content",
	})

	BroadcastToolFunction := api.ToolFunction{
		Name:        "broadcast",
		Description: "Spawn a persistent teammate that runs in its own thread.",
		Parameters: api.ToolFunctionParameters{
			Type:       "object",
			Required:   []string{"content"},
			Properties: props,
		},
	}
	return api.Tool{
		Type:     "function",
		Function: BroadcastToolFunction,
	}
}

func (b BroadcastTool) Run(args api.ToolCallFunctionArguments) string {
	contentRaw, _ := args.Get("content")
	content := contentRaw.(string)
	members := b.TeamManager.MemberNames()
	return b.Broadcaster.Broadcast("lead", content, members)
}
