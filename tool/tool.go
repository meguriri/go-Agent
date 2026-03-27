package tool

import "github.com/ollama/ollama/api"

type tool interface {
	GetTool() api.Tool
	Run(api.ToolCallFunctionArguments) string
}

type ToolHandler map[string]tool

func NewToolHandler() ToolHandler {
	t := make(ToolHandler)
	t["bash"] = bashTool
	t["read"] = RunRead
	return t
}

// func (t *toolHandler) addTools(tool api.ToolFunction) {
// 	t[tool.Name] = tool.Func
// }
