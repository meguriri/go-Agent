package toolManager

import "s01/tools"

type Deps struct {
	Messenger      tools.Messenger
	InboxReader    tools.InboxReader
	SubagentRunner tools.SubagentRunner
	TeamManager    tools.TeamManager
	Broadcaster    tools.Broadcaster
}

type ToolHandler map[string]tools.Tool

func NewToolHandler(deps Deps) ToolHandler {
	t := make(ToolHandler)
	t["bash"] = tools.BashTool{}
	t["read_file"] = tools.ReadTool{}
	t["write_file"] = tools.WriteTool{}
	t["edit_file"] = tools.EditTool{}
	t["background_run"] = tools.BackgroundTool{}
	t["check_background"] = tools.CheckBackgroundTool{}
	t["spawn_teammate"] = tools.SpawnTeammateTool{TeamManager: deps.TeamManager}
	t["list_teammates"] = tools.ListTeammatesTool{TeamManager: deps.TeamManager}
	t["send_message"] = tools.SendMessageTool{Messenger: deps.Messenger}
	t["read_inbox"] = tools.ReadInboxTool{Reader: deps.InboxReader}
	t["broadcast"] = tools.BroadcastTool{TeamManager: deps.TeamManager, Broadcaster: deps.Broadcaster}
	// t["todo"] = TodoManager{}
	// t["task"] = TaskTool{Runner: deps.SubagentRunner}
	// t["load_skill"] = LoadSkillTool{}
	// t["compact"] = CompactTool{}
	// t["task_create"] = TasksCreateTool{}
	// t["task_update"] = TasksUpdateTool{}
	// t["task_list"] = TasksListTool{}
	// t["task_get"] = TasksGetTool{}
	return t
}

func NewSubToolHandler() ToolHandler {
	t := make(ToolHandler)
	t["bash"] = tools.BashTool{}
	t["read_file"] = tools.ReadTool{}
	t["write_file"] = tools.WriteTool{}
	t["edit_file"] = tools.EditTool{}
	return t
}

func NewTeamToolHandler(deps Deps) ToolHandler {
	t := make(ToolHandler)
	t["bash"] = tools.BashTool{}
	t["read_file"] = tools.ReadTool{}
	t["write_file"] = tools.WriteTool{}
	t["edit_file"] = tools.EditTool{}
	t["send_message"] = tools.SendMessageTool{Messenger: deps.Messenger}
	t["read_inbox"] = tools.ReadInboxTool{Reader: deps.InboxReader}
	return t
}

func (t ToolHandler) AddTools(name string, tool tools.Tool) {
	t[name] = tool
}
