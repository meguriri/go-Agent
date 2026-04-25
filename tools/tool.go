package tools

import (
	"context"

	"github.com/ollama/ollama/api"
)

type Tool interface {
	GetTool() api.Tool
	Run(api.ToolCallFunctionArguments) string
}
type Messenger interface {
	Send(string, string, string, string, map[string]any) string
}
type InboxReader interface {
	ReadInboxText(string) string
}

type SubagentRunner interface {
	RunSubagent(context.Context, string, string) string
}
type TeamManager interface {
	Spawn(name string, role string, prompt string) string
	ListAll() string
	MemberNames() []string
}
type Broadcaster interface {
	Broadcast(sender string, content string, teamMates []string) string
}
