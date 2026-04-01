package model

import (
	"context"
	"fmt"
	"log"
	"s01/tools"
	"strings"

	"github.com/ollama/ollama/api"
)

const (
	THRESHOLD      = 50000
	TRANSCRIPT_DIR = ".transcripts"
	TASKS_DIR      = ".tasks"
	KEEP_RECENT    = 3
)

type Agent struct {
	client            *api.Client
	model             string
	ctx               context.Context
	System            string
	tools             api.Tools
	ToolHandler       tools.ToolHandler
	rounds_since_todo int
}

func NewAgent(client *api.Client, model string, ctx context.Context, system string,
	toolHandler tools.ToolHandler) *Agent {
	c := &Agent{
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

func (c *Agent) Chat() {
	history := make([]api.Message, 0)
	history = append(history, api.Message{
		Role:    "system",
		Content: c.System,
	})
	var query string
	i := 1
	for true {
		fmt.Print("\033[36ms01 >> \033[0m")
		// fmt.Scan(&query)
		if i == 1 {
			query = "逐个读取 sandbox 目录下的所有 cpp或 h 文件，分析它们的功能（不要修改文件，只分析功能）"
		} else if i == 2 {
			query = "持续读取sandbox的文件，直到压缩自动触发。"
		} else if i == 3 {
			query = "使用压缩工具（compact tool ）手动压缩对话"
		} else {
			break
		}
		query = strings.ToLower(strings.Trim(query, " "))
		if query == "q" || query == "exit" {
			break
		}
		history = append(history, api.Message{
			Role:    "user",
			Content: query,
		})
		history = c.AgentLoop(history)
		responses_content := history[len(history)-1].Content
		fmt.Println(responses_content)
		i++
	}
}

func (c *Agent) AgentLoop(messages []api.Message) []api.Message {
	for true {
		messages = MicroCompact(messages)
		if EstimateTokens(messages) > THRESHOLD {
			fmt.Println("[auto_compact triggered]")
			messages = AutoCompact(c, messages)
		}

		req := &api.ChatRequest{
			Model:    c.model,
			Messages: messages,
			Tools:    c.tools,
		}
		var fullContent strings.Builder
		var thinkingContent strings.Builder
		var assistantMsg api.Message
		use_todo := false
		manual_compact := false

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
			return messages
		}

		assistantMsg.Role = "assistant"
		assistantMsg.Thinking = thinkingContent.String()
		assistantMsg.Content = fullContent.String()
		messages = append(messages, assistantMsg)
		if assistantMsg.Thinking != "" {
			fmt.Printf("\033[90m正在思考：%s\033[0m\n", assistantMsg.Thinking)
		}
		if len(assistantMsg.ToolCalls) == 0 {
			return messages
		}
		for _, tc := range assistantMsg.ToolCalls {
			fmt.Printf("\033[33m$ 正在执行工具: %s\033[0m\n", tc.Function.Name)
			var output string
			handler, ok := c.ToolHandler[tc.Function.Name]
			if !ok {
				output = "Unknown tool: " + tc.Function.Name
			} else {
				if tc.Function.Name == "todo" {
					use_todo = true
					output = handler.Run(tc.Function.Arguments)
				} else if tc.Function.Name == "compact" {
					manual_compact = true
					output = "Compressing..."
					continue
				} else {
					output = handler.Run(tc.Function.Arguments)
				}
			}
			if tc.Function.Name != "todo" {
				fmt.Printf("执行结果摘要: %s\n", strings.Split(output, "\n")[0])
			} else {
				fmt.Printf("\033[32m 更新后的待办事项:\n%s \033[0m\n", output)
			}
			toolResultMsg := api.Message{
				Role:      "tool",
				Content:   output,
				ToolCalls: []api.ToolCall{tc},
			}
			messages = append(messages, toolResultMsg)
		}
		if use_todo {
			c.rounds_since_todo = 0
		} else {
			c.rounds_since_todo++
		}
		if manual_compact {
			fmt.Println("[manual compact]")
			messages = AutoCompact(c, messages)
		}
		// if c.rounds_since_todo >= 3 {
		// 	messages = append(messages, api.Message{
		// 		Role:    "user",
		// 		Content: "<reminder>Update your todos.</reminder>",
		// 	})
		// }
	}
	return messages
}
