package tools

import (
	"github.com/ollama/ollama/api"
)

type TasksCreateTool struct{}

func (t TasksCreateTool) GetTool() api.Tool {
	props := api.NewToolPropertiesMap()
	props.Set("subject", api.ToolProperty{
		Type:        api.PropertyType{"string"},
		Description: "task's subject",
	})
	props.Set("description", api.ToolProperty{
		Type:        api.PropertyType{"string"},
		Description: "task's  description",
	})

	tasksCreateToolFunction := api.ToolFunction{
		Name:        "task_create",
		Description: "Create a new task.",
		Parameters: api.ToolFunctionParameters{
			Type:       "object",
			Required:   []string{"subject"},
			Properties: props,
		},
	}
	return api.Tool{
		Type:     "function",
		Function: tasksCreateToolFunction,
	}
}

func (t TasksCreateTool) Run(args api.ToolCallFunctionArguments) string {
	return ""
}

type TasksUpdateTool struct{}

func (t TasksUpdateTool) GetTool() api.Tool {
	props := api.NewToolPropertiesMap()
	props.Set("task_id", api.ToolProperty{
		Type:        api.PropertyType{"integer"},
		Description: "ID of the task to update",
	})
	props.Set("status", api.ToolProperty{
		Type:        api.PropertyType{"string"},
		Description: "task's status one of 'pending', 'in_progress', 'completed'",
		Enum:        []interface{}{"pending", "in_progress", "completed"},
	})
	props.Set("addBlockedBy", api.ToolProperty{
		Type:        api.PropertyType{"array"},
		Description: "",
		Items: &api.ToolProperty{
			Type: api.PropertyType{"integer"},
		},
	})
	props.Set("addBlocks", api.ToolProperty{
		Type:        api.PropertyType{"array"},
		Description: "",
		Items: &api.ToolProperty{
			Type: api.PropertyType{"integer"},
		},
	})

	tasksUpdateToolFunction := api.ToolFunction{
		Name:        "task_update",
		Description: "Update a task's status or dependencies.",
		Parameters: api.ToolFunctionParameters{
			Type:       "object",
			Required:   []string{"task_id"},
			Properties: props,
		},
	}
	return api.Tool{
		Type:     "function",
		Function: tasksUpdateToolFunction,
	}
}

func (t TasksUpdateTool) Run(args api.ToolCallFunctionArguments) string {
	return ""
}

type TasksListTool struct{}

func (t TasksListTool) GetTool() api.Tool {
	props := api.NewToolPropertiesMap()

	tasksListToolFunction := api.ToolFunction{
		Name:        "task_list",
		Description: "List all tasks with status summary.",
		Parameters: api.ToolFunctionParameters{
			Type:       "object",
			Properties: props,
		},
	}
	return api.Tool{
		Type:     "function",
		Function: tasksListToolFunction,
	}
}

func (t TasksListTool) Run(args api.ToolCallFunctionArguments) string {
	return ""
}

type TasksGetTool struct{}

func (t TasksGetTool) GetTool() api.Tool {
	props := api.NewToolPropertiesMap()
	props.Set("task_id", api.ToolProperty{
		Type:        api.PropertyType{"integer"},
		Description: "ID of the task to retrieve",
	})

	tasksGetToolFunction := api.ToolFunction{
		Name:        "task_get",
		Description: "Get full details of a task by ID.",
		Parameters: api.ToolFunctionParameters{
			Type:       "object",
			Required:   []string{"task_id"},
			Properties: props,
		},
	}
	return api.Tool{
		Type:     "function",
		Function: tasksGetToolFunction,
	}
}

func (t TasksGetTool) Run(args api.ToolCallFunctionArguments) string {
	return ""
}
