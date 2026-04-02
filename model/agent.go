package model

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"s01/background"
	"s01/tools"
	"strings"

	"github.com/ollama/ollama/api"
)

// const (
// 	THRESHOLD      = 50000
// 	TRANSCRIPT_DIR = ".transcripts"
// 	TASKS_DIR      = ".tasks"
// 	KEEP_RECENT    = 3
// )

type Agent struct {
	client            *api.Client
	model             string
	ctx               context.Context
	threshold         int
	System            string
	tools             api.Tools
	ToolHandler       tools.ToolHandler
	rounds_since_todo int
}

func NewAgent(client *api.Client, model string, ctx context.Context, threshold int, system string,
	toolHandler tools.ToolHandler) *Agent {
	c := &Agent{
		client:            client,
		model:             model,
		ctx:               ctx,
		threshold:         threshold,
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
	i := 1
	for true {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("\033[36ms01 >> \033[0m")
		query, err := reader.ReadString('\n')
		if err != nil {
			return
		}
		// if i == 1 {
		// 	query = "在后台运行 `sleep 5 && echo done`，然后在其运行期间创建一个文件。"
		// } else if i == 2 {
		// 	query = "启动 3 个后台任务：“sleep 2”、“sleep 4”、“sleep 6”。检查它们的状态。"
		// } else if i == 3 {
		// 	query = "在后台运行一个在sandbox目录下运行test.py,同时继续接受我的其他任务"
		// } else {
		// 	query = "exit"
		// }
		// fmt.Printf("\033[36ms01 >>%s \033[0m\n", query)
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
		fmt.Println("\n" + responses_content)
		i++
	}
}

func (c *Agent) AgentLoop(messages []api.Message) []api.Message {
	for true {
		messages = MicroCompact(messages)
		if EstimateTokens(messages) > c.threshold {
			fmt.Println("[auto_compact triggered]")
			messages = AutoCompact(c, messages)
		}

		notifs := background.MyBackgroundManager.DrainNotifications()
		if notifs != nil && len(messages) != 0 {
			lines := make([]string, 0)
			for _, n := range notifs {
				lines = append(lines, fmt.Sprintf("[bg:%s] %s: %s", n.ID, n.Status, n.Result))
			}
			notifText := strings.Join(lines, "\n")
			messages = append(messages, api.Message{
				Role:    "user",
				Content: fmt.Sprintf("<background-results>\n%s\n</background-results>", notifText),
			})
			messages = append(messages, api.Message{
				Role:    "assistant",
				Content: "Noted background results.",
			})
		}

		req := &api.ChatRequest{
			Model:    c.model,
			Messages: messages,
			Tools:    c.tools,
		}

		var assistantMsg api.Message
		use_todo := false
		manual_compact := false
		first_thinking := 0
		err := c.client.Chat(c.ctx, req, func(resp api.ChatResponse) error {
			if resp.Message.Thinking != "" {

				assistantMsg.Thinking += resp.Message.Thinking
				if first_thinking == 0 {
					fmt.Printf("\033[90m正在思考：\033[0m")
				}
				fmt.Printf("\033[90m%s\033[0m", resp.Message.Thinking)
				first_thinking++
			}
			if resp.Message.Content != "" {
				assistantMsg.Content += resp.Message.Content
			}
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
		messages = append(messages, assistantMsg)

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
