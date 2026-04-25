package tools

import (
	"github.com/ollama/ollama/api"
)

type SpawnTeammateTool struct {
	TeamManager TeamManager
}

func (s SpawnTeammateTool) GetTool() api.Tool {
	props := api.NewToolPropertiesMap()
	props.Set("name", api.ToolProperty{
		Type:        api.PropertyType{"string"},
		Description: "this teamAgent's name",
	})
	props.Set("role", api.ToolProperty{
		Type:        api.PropertyType{"string"},
		Description: "this teamAgent's role",
	})
	props.Set("prompt", api.ToolProperty{
		Type:        api.PropertyType{"string"},
		Description: "this teamAgent's prompt",
	})

	SpawnTeammateToolFunction := api.ToolFunction{
		Name:        "spawn_teammate",
		Description: "Spawn a persistent teammate that runs in its own thread.",
		Parameters: api.ToolFunctionParameters{
			Type:       "object",
			Required:   []string{"name", "role", "prompt"},
			Properties: props,
		},
	}
	return api.Tool{
		Type:     "function",
		Function: SpawnTeammateToolFunction,
	}
}

func (s SpawnTeammateTool) Run(args api.ToolCallFunctionArguments) string {
	nameRaw, _ := args.Get("name")
	roleRaw, _ := args.Get("role")
	promptRaw, _ := args.Get("prompt")
	name := nameRaw.(string)
	role := roleRaw.(string)
	prompt := promptRaw.(string)
	return s.TeamManager.Spawn(name, role, prompt)
}
