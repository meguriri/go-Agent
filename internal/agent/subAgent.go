package agent

import (
	"context"
	"fmt"
	"log"
	"os"
	internal "s01/internal/toolManager"
	"strings"

	"github.com/ollama/ollama/api"
)

type SubagentService struct {
	Client *api.Client
	Model  string
}

func (s *SubagentService) RunSubagent(ctx context.Context, prompt string, description string) string {
	var systemPrompt string
	dir, _ := os.Getwd()
	if description == "" {
		systemPrompt = fmt.Sprintf("You are a coding subagent at %s. Complete the given task, then summarize your findings.", dir)
	} else {
		systemPrompt = fmt.Sprintf("You are a coding subagent at %s. Complete the given task. The task's description is %s, then summarize your findings.", dir, description)
	}
	sub := NewSubAgent(s.Client, s.Model, ctx, systemPrompt, internal.NewSubToolHandler(), prompt, description)
	out := sub.AgentLoop([]api.Message{{Role: "user", Content: prompt}})
	if len(out) == 0 {
		return "(no summary)"
	}
	return out[len(out)-1].Content
}

type SubAgent struct {
	client            *api.Client
	model             string
	ctx               context.Context
	System            string
	tools             api.Tools
	ToolHandler       internal.ToolHandler
	prompt            string
	description       string
	rounds_since_todo int
}

func NewSubAgent(client *api.Client, model string, ctx context.Context, system string,
	toolHandler internal.ToolHandler, prompt string, description string) *SubAgent {
	c := &SubAgent{
		client:            client,
		model:             model,
		ctx:               ctx,
		System:            system,
		tools:             nil,
		ToolHandler:       toolHandler,
		prompt:            prompt,
		description:       description,
		rounds_since_todo: 0,
	}
	for _, v := range c.ToolHandler {
		c.tools = append(c.tools, v.GetTool())
	}
	return c
}

func (c SubAgent) Model() string {
	return c.model
}
func (c *SubAgent) Client() *api.Client {
	return c.client
}
func (c *SubAgent) Ctx() context.Context {
	return c.ctx
}

func (c *SubAgent) AgentLoop(messages []api.Message) []api.Message {
	for i := 1; i <= 30; i++ {
		req := &api.ChatRequest{
			Model:    c.model,
			Messages: messages,
			Tools:    c.tools,
		}
		var fullContent strings.Builder
		var thinkingContent strings.Builder
		var assistantMsg api.Message
		use_todo := false
		err := c.client.Chat(c.ctx, req, func(resp api.ChatResponse) error {
			fullContent.WriteString(resp.Message.Content)
			thinkingContent.WriteString(resp.Message.Thinking)
			if len(resp.Message.ToolCalls) > 0 {
				assistantMsg.ToolCalls = append(assistantMsg.ToolCalls, resp.Message.ToolCalls...)
			}
			return nil
		})
		if err != nil {
			log.Fatalf("get llm response error: %v\n", err)
			return []api.Message{
				api.Message{
					Role:    "user",
					Content: "get llm response error: " + err.Error(),
				},
			}
		}
		assistantMsg.Role = "assistant"
		assistantMsg.Thinking = thinkingContent.String()
		assistantMsg.Content = fullContent.String()
		messages = append(messages, assistantMsg)
		if len(assistantMsg.ToolCalls) == 0 {
			break
		}
		for _, tc := range assistantMsg.ToolCalls {
			fmt.Printf("	\033[33m$ 正在执行工具: %s\033[0m\n", tc.Function.Name)
			var output string
			handler, ok := c.ToolHandler[tc.Function.Name]
			if !ok {
				output = "Unknown tool: " + tc.Function.Name
			} else {
				if tc.Function.Name == "todo" {
					use_todo = true
				}
				output = handler.Run(tc.Function.Arguments)
			}
			if tc.Function.Name != "todo" {
				fmt.Printf("	执行结果摘要: %s\n", strings.Split(output, "\n")[0])
			} else {
				fmt.Printf("	\033[32m 更新后的待办事项:\n%s \033[0m\n", output)
			}
			toolResultMsg := api.Message{
				Role:    "tool",
				Content: output,
			}
			messages = append(messages, toolResultMsg)
		}
		if use_todo {
			c.rounds_since_todo = 0
		} else {
			c.rounds_since_todo++
		}
		if c.rounds_since_todo >= 3 {
			messages = append(messages, api.Message{
				Role:    "user",
				Content: "<reminder>Update your todos.</reminder>",
			})
		}
	}
	result := messages[len(messages)-1].Content
	if result == "" {
		return []api.Message{
			api.Message{
				Role:    "user",
				Content: "(no summary)",
			},
		}
	}
	return messages[len(messages)-1:]
}
