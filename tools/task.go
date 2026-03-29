package tools

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/ollama/ollama/api"
)

type TaskTool struct{}

type SubAgent struct {
	client            *api.Client
	model             string
	ctx               context.Context
	System            string
	tools             api.Tools
	ToolHandler       ToolHandler
	rounds_since_todo int
}

func (t TaskTool) GetTool() api.Tool {
	props := api.NewToolPropertiesMap()
	props.Set("prompt", api.ToolProperty{
		Type:        api.PropertyType{"string"},
		Description: "prompt for the subagent",
	})
	props.Set("description", api.ToolProperty{
		Type:        api.PropertyType{"string"},
		Description: "Short description of the task",
	})

	taskToolFunction := api.ToolFunction{
		Name:        "task",
		Description: "生成一个具有全新上下文的subagent。它共享文件系统，但不共享对话历史。",
		Parameters: api.ToolFunctionParameters{
			Type:       "object",
			Required:   []string{"prompt"},
			Properties: props,
		},
	}
	return api.Tool{
		Type:     "function",
		Function: taskToolFunction,
	}
}

func (t TaskTool) Run(args api.ToolCallFunctionArguments) string {
	promptRaw, _ := args.Get("prompt")
	descriptionRaw, ok := args.Get("description")
	var description string
	if ok {
		description = descriptionRaw.(string)
	}
	prompt := promptRaw.(string)
	ctx := context.Background()
	result := Run_subagent(prompt, description, ctx)
	return result
}

func NewSubAgent(client *api.Client, model string, ctx context.Context, system string,
	toolHandler ToolHandler) *SubAgent {
	c := &SubAgent{
		client:            client,
		model:             model,
		ctx:               ctx,
		System:            system,
		tools:             nil,
		ToolHandler:       toolHandler,
		rounds_since_todo: 0,
	}
	for _, v := range c.ToolHandler {
		c.tools = append(c.tools, v.GetTool())
	}
	return c
}

func Run_subagent(prompt string, description string, ctx context.Context) string {
	modelID := os.Getenv("OLLAMA_MODELS")
	c, err := api.ClientFromEnvironment()
	if err != nil {
		log.Fatalf("Create sub client error: %v\n", err)
	}
	dir, _ := os.Getwd()
	var systemPrompt string
	if description == "" {
		systemPrompt = fmt.Sprintf("You are a coding subagent at %s. Complete the given task, then summarize your findings.", dir)
	} else {
		systemPrompt = fmt.Sprintf("You are a coding subagent at %s. Complete the given task. The task's description is %s, then summarize your findings.", dir, description)
	}
	subAgent := NewSubAgent(c, modelID, ctx,
		systemPrompt,
		NewSubToolHandler(),
	)
	sub_messages := []api.Message{{
		Role:    "user",
		Content: prompt,
	}}
	for i := 1; i <= 30; i++ {
		req := &api.ChatRequest{
			Model:    subAgent.model,
			Messages: sub_messages,
			Tools:    subAgent.tools,
		}
		var fullContent strings.Builder
		var thinkingContent strings.Builder
		var assistantMsg api.Message
		use_todo := false
		err := subAgent.client.Chat(subAgent.ctx, req, func(resp api.ChatResponse) error {
			fullContent.WriteString(resp.Message.Content)
			thinkingContent.WriteString(resp.Message.Thinking)
			if len(resp.Message.ToolCalls) > 0 {
				assistantMsg.ToolCalls = append(assistantMsg.ToolCalls, resp.Message.ToolCalls...)
			}
			return nil
		})
		if err != nil {
			log.Fatalf("get llm response error: %v\n", err)
			return "get llm response error: " + err.Error()
		}
		assistantMsg.Role = "assistant"
		assistantMsg.Thinking = thinkingContent.String()
		assistantMsg.Content = fullContent.String()
		sub_messages = append(sub_messages, assistantMsg)
		if len(assistantMsg.ToolCalls) == 0 {
			break
		}
		if assistantMsg.Thinking != "" {
			fmt.Printf("\033[90m正在思考：%s\033[0m\n", assistantMsg.Thinking)
		}
		for _, tc := range assistantMsg.ToolCalls {
			fmt.Printf("\033[33m$ 正在执行工具: %s\033[0m\n", tc.Function.Name)
			var output string
			handler, ok := subAgent.ToolHandler[tc.Function.Name]
			if !ok {
				output = "Unknown tool: " + tc.Function.Name
			} else {
				if tc.Function.Name == "todo" {
					use_todo = true
				}
				output = handler.Run(tc.Function.Arguments)
			}
			if tc.Function.Name != "todo" {
				fmt.Printf("执行结果摘要: %s\n", strings.Split(output, "\n")[0])
			} else {
				fmt.Printf("\033[32m 更新后的待办事项:\n%s \033[0m\n", output)
			}
			toolResultMsg := api.Message{
				Role:    "tool",
				Content: output,
			}
			sub_messages = append(sub_messages, toolResultMsg)
		}
		if use_todo {
			subAgent.rounds_since_todo = 0
		} else {
			subAgent.rounds_since_todo++
		}
		if subAgent.rounds_since_todo >= 3 {
			sub_messages = append(sub_messages, api.Message{
				Role:    "user",
				Content: "<reminder>Update your todos.</reminder>",
			})
		}
	}
	result := sub_messages[len(sub_messages)-1].Content
	if result == "" {
		return "(no summary)"
	}
	return result
}
