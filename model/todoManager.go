package model

// import (
// 	"fmt"
// 	"strings"
// )

// type item struct {
// 	id     string
// 	status string
// 	text   string
// }

// type TodoManager struct {
// 	items []item
// }

// // 把items渲染成多行字符串返回
// func (t *TodoManager) Render() string {
// 	if len(t.items) == 0 {
// 		return "No todos."
// 	}
// 	lines := make([]string, 0)
// 	done := 0
// 	for _, item := range t.items {
// 		var marker string
// 		switch item.status {
// 		case "pending":
// 			marker = "[ ]"
// 		case "in_progress":
// 			marker = "[>]"
// 		case "completed":
// 			marker = "[x]"
// 		default:
// 			marker = "[?]"
// 		}
// 		lines = append(lines, fmt.Sprintf("%s #%s: %s", marker, item.id, item.text))
// 		if item.status == "completed" {
// 			done++
// 		}
// 	}
// 	lines = append(lines, fmt.Sprintf("\n(%d/%d completed)", done, len(t.items)))
// 	return strings.Join(lines, "\n")
// }

// // 更新item的状态，并返回更新后的渲染结果
// func (t *TodoManager) Update(items []item) string {
// 	if len(t.items) > 20 {
// 		return "Max 20 todos allowed"
// 	}
// 	validated := make([]item, 0)
// 	inProgressCount := 0
// 	for i, it := range items {
// 		text := strings.TrimSpace(it.text)
// 		status := strings.ToLower(it.status)
// 		var item_id string
// 		if it.id == "" {
// 			it.id = fmt.Sprintf("%d", i+1)
// 		} else {
// 			item_id = it.id
// 		}
// 		if text == "" {
// 			return fmt.Sprintf("Item %d: text required", item_id)
// 		}
// 		if status != "pending" && status != "in_progress" && status != "completed" {
// 			return fmt.Sprintf("Item %d: invalid status '%s'", item_id, it.status)
// 		}
// 		if status == "in_progress" {
// 			inProgressCount++
// 		}
// 		validated = append(validated, item{
// 			id:     item_id,
// 			text:   text,
// 			status: status,
// 		})
// 	}
// 	if inProgressCount > 1 {
// 		return "Only one task can be in_progress at a time"
// 	}
// 	t.items = validated
// 	return t.Render()
// }
