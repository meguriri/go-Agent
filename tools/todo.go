package tools

import (
	"fmt"
	"strings"

	"github.com/ollama/ollama/api"
)

type item struct {
	id     string
	status string
	text   string
}

type TodoManager struct {
	items []item
}

// 把items渲染成多行字符串返回
func (t *TodoManager) Render() string {
	if len(t.items) == 0 {
		return "No todos."
	}
	lines := make([]string, 0)
	done := 0
	for _, item := range t.items {
		var marker string
		switch item.status {
		case "pending":
			marker = "[ ]"
		case "in_progress":
			marker = "[>]"
		case "completed":
			marker = "[x]"
		default:
			marker = "[?]"
		}
		lines = append(lines, fmt.Sprintf("%s #%s: %s", marker, item.id, item.text))
		if item.status == "completed" {
			done++
		}
	}
	lines = append(lines, fmt.Sprintf("\n(%d/%d completed)", done, len(t.items)))
	return strings.Join(lines, "\n")
}

// 更新item的状态，并返回更新后的渲染结果
func (t *TodoManager) Update(items []item) string {
	if len(t.items) > 20 {
		return "Max 20 todos allowed"
	}
	validated := make([]item, 0)
	inProgressCount := 0
	for i, it := range items {
		text := strings.TrimSpace(it.text)
		status := strings.ToLower(it.status)
		var item_id string
		if it.id == "" {
			it.id = fmt.Sprintf("%d", i+1)
		} else {
			item_id = it.id
		}
		if text == "" {
			return fmt.Sprintf("Item %d: text required", item_id)
		}
		if status != "pending" && status != "in_progress" && status != "completed" {
			return fmt.Sprintf("Item %d: invalid status '%s'", item_id, it.status)
		}
		if status == "in_progress" {
			inProgressCount++
		}
		validated = append(validated, item{
			id:     item_id,
			text:   text,
			status: status,
		})
	}
	if inProgressCount > 1 {
		return "Only one task can be in_progress at a time"
	}
	t.items = validated
	return t.Render()
}

func (t TodoManager) GetTool() api.Tool {
	itemProps := api.NewToolPropertiesMap()
	itemProps.Set("id", api.ToolProperty{
		Type:        api.PropertyType{"string"},
		Description: "Task unique identifier",
	})
	itemProps.Set("text", api.ToolProperty{
		Type:        api.PropertyType{"string"},
		Description: "Task description",
	})
	itemProps.Set("status", api.ToolProperty{
		Type:        api.PropertyType{"string"},
		Description: "Task status, one of 'pending', 'in_progress', 'completed'",
		Enum:        []interface{}{"pending", "in_progress", "completed"},
	})
	todoItemSchema := api.ToolFunctionParameters{
		Type:       "object",
		Properties: itemProps,
		Required:   []string{"id", "text", "status"},
	}
	props := api.NewToolPropertiesMap()
	props.Set("items", api.ToolProperty{
		Type:        api.PropertyType{"array"},
		Description: "List of todo items to update",
		Items:       &todoItemSchema,
	})

	todoToolFunction := api.ToolFunction{
		Name:        "todo",
		Description: "Update task list. Track progress on multi-step tasks.",
		Parameters: api.ToolFunctionParameters{
			Type:       "object",
			Required:   []string{"items"},
			Properties: props,
		},
	}
	return api.Tool{
		Type:     "function",
		Function: todoToolFunction,
	}
}

func (t TodoManager) Run(args api.ToolCallFunctionArguments) string {
	itemsRaw, _ := args.Get("items")
	itemsI := itemsRaw.([]interface{})
	items := make([]item, 0)
	for _, it := range itemsI {
		itMap := it.(map[string]interface{})
		item := item{
			id:     itMap["id"].(string),
			text:   itMap["text"].(string),
			status: itMap["status"].(string),
		}
		items = append(items, item)
	}
	return t.Update(items)
}
