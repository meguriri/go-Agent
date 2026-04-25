package tools

import (
	"s01/internal/task"

	"github.com/ollama/ollama/api"
)

// var taskManager = task.NewTaskManager(".tasks")

type TasksCreateTool struct {
}

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
	subjectRaw, _ := args.Get("subject")
	subject := subjectRaw.(string)
	descriptionRaw, ok := args.Get("description")
	description := ""
	if ok {
		description = descriptionRaw.(string)
	}
	return task.MyTaskManager.Create(subject, description)
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
		Description: "前置依赖任务id的列表",
		Items: &api.ToolProperty{
			Type: api.PropertyType{"integer"},
		},
	})
	props.Set("addBlocks", api.ToolProperty{
		Type:        api.PropertyType{"array"},
		Description: "后置任务id的列表",
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
	taskIdRaw, ok := args.Get("task_id")
	var taskId int
	if ok {
		switch v := taskIdRaw.(type) {
		case int:
			taskId = v
		case int8:
			taskId = int(v)
		case int16:
			taskId = int(v)
		case int32:
			taskId = int(v)
		case int64:
			taskId = int(v)
		case uint:
			taskId = int(v)
		case uint8:
			taskId = int(v)
		case uint16:
			taskId = int(v)
		case uint32:
			taskId = int(v)
		case uint64:
			taskId = int(v)
		case float32:
			taskId = int(v)
		case float64:
			taskId = int(v)
		}
	}
	statusRaw, ok := args.Get("status")
	var status string
	if ok {
		status = statusRaw.(string)
	}
	addBlockByRaw, ok := args.Get("addBlockedBy")
	var addBlockBy []int = nil
	if ok {
		addBlockByInterface := addBlockByRaw.([]interface{})
		addBlockBy = make([]int, len(addBlockByInterface))
		for i, v := range addBlockByInterface {
			switch val := v.(type) {
			case int:
				addBlockBy[i] = val
			case int8:
				addBlockBy[i] = int(val)
			case int16:
				addBlockBy[i] = int(val)
			case int32:
				addBlockBy[i] = int(val)
			case int64:
				addBlockBy[i] = int(val)
			case uint:
				addBlockBy[i] = int(val)
			case uint8:
				addBlockBy[i] = int(val)
			case uint16:
				addBlockBy[i] = int(val)
			case uint32:
				addBlockBy[i] = int(val)
			case uint64:
				addBlockBy[i] = int(val)
			case float32:
				addBlockBy[i] = int(val)
			case float64:
				addBlockBy[i] = int(val)
			}
		}
	}
	addBlocksRaw, ok := args.Get("addBlocks")
	var addBlocks []int = nil
	if ok {
		addBlocksInterface := addBlocksRaw.([]interface{})
		addBlocks = make([]int, len(addBlocksInterface))
		for i, v := range addBlocksInterface {
			switch val := v.(type) {
			case int:
				addBlocks[i] = val
			case int8:
				addBlocks[i] = int(val)
			case int16:
				addBlocks[i] = int(val)
			case int32:
				addBlocks[i] = int(val)
			case int64:
				addBlocks[i] = int(val)
			case uint:
				addBlocks[i] = int(val)
			case uint8:
				addBlocks[i] = int(val)
			case uint16:
				addBlocks[i] = int(val)
			case uint32:
				addBlocks[i] = int(val)
			case uint64:
				addBlocks[i] = int(val)
			case float32:
				addBlocks[i] = int(val)
			case float64:
				addBlocks[i] = int(val)
			}
		}
	}
	return task.MyTaskManager.Update(taskId, status, addBlockBy, addBlocks)
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
	return task.MyTaskManager.ListAll()
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
	taskIdRaw, ok := args.Get("task_id")
	var taskId int
	if ok {
		switch v := taskIdRaw.(type) {
		case int:
			taskId = v
		case int8:
			taskId = int(v)
		case int16:
			taskId = int(v)
		case int32:
			taskId = int(v)
		case int64:
			taskId = int(v)
		case uint:
			taskId = int(v)
		case uint8:
			taskId = int(v)
		case uint16:
			taskId = int(v)
		case uint32:
			taskId = int(v)
		case uint64:
			taskId = int(v)
		case float32:
			taskId = int(v)
		case float64:
			taskId = int(v)
		}
	}
	return task.MyTaskManager.Get(taskId)
}
