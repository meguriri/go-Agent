package tools

import (
	"s01/internal/skill"

	"github.com/ollama/ollama/api"
)

type LoadSkillTool struct{}

func (r LoadSkillTool) GetTool() api.Tool {
	props := api.NewToolPropertiesMap()
	props.Set("name", api.ToolProperty{
		Type:        api.PropertyType{"string"},
		Description: "Skill name to load",
	})

	loadSkillToolFunction := api.ToolFunction{
		Name:        "load_skill",
		Description: "Load specialized knowledge by name.",
		Parameters: api.ToolFunctionParameters{
			Type:       "object",
			Required:   []string{"name"},
			Properties: props,
		},
	}
	return api.Tool{
		Type:     "function",
		Function: loadSkillToolFunction,
	}
}

func (r LoadSkillTool) Run(args api.ToolCallFunctionArguments) string {
	nameRaw, _ := args.Get("name")
	name := nameRaw.(string)
	skillsLoader := skill.NewSkillLoader("./skills")
	return skillsLoader.GetContent(name)
}
