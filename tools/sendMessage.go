package tools

import (
	"github.com/ollama/ollama/api"
)

type SendMessageTool struct {
	Messenger Messenger
}

func (s SendMessageTool) GetTool() api.Tool {
	props := api.NewToolPropertiesMap()
	props.Set("sender", api.ToolProperty{
		Type:        api.PropertyType{"string"},
		Description: "which agent send this message",
	})
	props.Set("to", api.ToolProperty{
		Type:        api.PropertyType{"string"},
		Description: "send to which agent",
	})
	props.Set("content", api.ToolProperty{
		Type:        api.PropertyType{"string"},
		Description: "message content",
	})
	props.Set("msg_type", api.ToolProperty{
		Type:        api.PropertyType{"string"},
		Description: "message type",
		Enum: []any{
			"message", "broadcast", "shutdown_request", "shutdown_response", "plan_approval_response",
		},
	})

	sendMessageToolFunction := api.ToolFunction{
		Name:        "send_message",
		Description: "Send message to a teammate.",
		Parameters: api.ToolFunctionParameters{
			Type:       "object",
			Required:   []string{"sender", "to", "content"},
			Properties: props,
		},
	}
	return api.Tool{
		Type:     "function",
		Function: sendMessageToolFunction,
	}
}

func (s SendMessageTool) Run(args api.ToolCallFunctionArguments) string {
	senderRaw, _ := args.Get("sender")
	toRaw, _ := args.Get("to")
	contentRaw, _ := args.Get("content")
	msgTypeRaw, ok := args.Get("msg_type")
	sender := senderRaw.(string)
	to := toRaw.(string)
	content := contentRaw.(string)
	var msgType string = "message"
	if ok {
		msgType = msgTypeRaw.(string)
	}
	return s.Messenger.Send(sender, to, content, msgType, nil)

	// return team.BUS.Send(sender, to, content, msgType, nil)
}
