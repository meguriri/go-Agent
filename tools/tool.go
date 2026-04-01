package tools

import (
	"github.com/ollama/ollama/api"
)

type tool interface {
	GetTool() api.Tool
	Run(api.ToolCallFunctionArguments) string
}

type ToolHandler map[string]tool

func NewToolHandler() ToolHandler {
	t := make(ToolHandler)
	t["bash"] = BashTool{}
	// t["read_file"] = ReadTool{}
	// t["write_file"] = WriteTool{}
	// t["edit_file"] = EditTool{}
	// t["todo"] = TodoManager{}
	// t["task"] = TaskTool{}
	t["load_skill"] = LoadSkillTool{}
	return t
}

func NewSubToolHandler() ToolHandler {
	t := make(ToolHandler)
	t["bash"] = BashTool{}
	t["read_file"] = ReadTool{}
	t["write_file"] = WriteTool{}
	t["edit_file"] = EditTool{}
	t["todo"] = TodoManager{}
	return t
}

func (t ToolHandler) AddTools(name string, tool tool) {
	t[name] = tool
}
