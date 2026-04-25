package tools

import (
	"github.com/ollama/ollama/api"
)

type ListTeammatesTool struct {
	TeamManager TeamManager
}

func (l ListTeammatesTool) GetTool() api.Tool {
	props := api.NewToolPropertiesMap()

	ListTeammatesToolFunction := api.ToolFunction{
		Name:        "list_teammates",
		Description: "List all teammates with name, role, status.",
		Parameters: api.ToolFunctionParameters{
			Type:       "object",
			Properties: props,
		},
	}
	return api.Tool{
		Type:     "function",
		Function: ListTeammatesToolFunction,
	}
}

func (l ListTeammatesTool) Run(args api.ToolCallFunctionArguments) string {
	return l.TeamManager.ListAll()
}
