package tool

import (
	"s01/tool/tools"

	"github.com/ollama/ollama/api"
)

type tool interface {
	GetTool() api.Tool
	Run(api.ToolCallFunctionArguments) string
}

type ToolHandler map[string]tool

func NewToolHandler() ToolHandler {
	t := make(ToolHandler)
	t["bash"] = tools.BashTool{}
	t["read_file"] = tools.ReadTool{}
	t["write_file"] = tools.WriteTool{}
	t["edit_file"] = tools.EditTool{}
	return t
}

func (t ToolHandler) AddTools(name string, tool tool) {
	t[name] = tool
}
