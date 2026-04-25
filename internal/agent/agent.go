package agent

import (
	"context"

	"github.com/ollama/ollama/api"
)

type Agent interface {
	AgentLoop(messages []api.Message) []api.Message
	Model() string
	Client() *api.Client
	Ctx() context.Context
}

type TeamManager interface {
	ListAll() string
}

type Inbox interface {
	ReadInboxText(string) string
	ReadInboxMessages(string) []api.Message
}
