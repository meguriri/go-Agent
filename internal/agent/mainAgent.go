package agent

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"s01/internal/toolManager"
	"strings"

	"github.com/ollama/ollama/api"
)

// const (
// 	THRESHOLD      = 50000
// 	TRANSCRIPT_DIR = ".transcripts"
// 	TASKS_DIR      = ".tasks"
// 	KEEP_RECENT    = 3
// )

type MainAgent struct {
	client            *api.Client
	model             string
	ctx               context.Context
	threshold         int
	System            string
	tools             api.Tools
	ToolHandler       toolManager.ToolHandler
	rounds_since_todo int
	inbox             Inbox
	teamManager       TeamManager
}

func NewMainAgent(client *api.Client, model string, ctx context.Context, threshold int, system string,
	toolHandler toolManager.ToolHandler, inbox Inbox, teamManager TeamManager) *MainAgent {
	c := &MainAgent{
		client:            client,
		model:             model,
		ctx:               ctx,
		threshold:         threshold,
		System:            system,
		tools:             nil,
		ToolHandler:       toolHandler,
		rounds_since_todo: 0,
		inbox:             inbox,
		teamManager:       teamManager,
	}
	for _, v := range c.ToolHandler {
		c.tools = append(c.tools, v.GetTool())
	}
	return c
}

func (c MainAgent) Model() string {
	return c.model
}
func (c *MainAgent) Client() *api.Client {
	return c.client
}
func (c *MainAgent) Ctx() context.Context {
	return c.ctx
}

func (c *MainAgent) Chat() {
	history := make([]api.Message, 0)
	history = append(history, api.Message{
		Role:    "system",
		Content: c.System,
	})
	i := 1
	for true {
		fmt.Print("\033[36ms01 >> \033[0m")
		reader := bufio.NewReader(os.Stdin)
		query, err := reader.ReadString('\n')
		if err != nil {
			return
		}
		// var query string
		// if i == 1 {
		// 	query = "生成 Alice（程序员）和 Bob（测试员）。让 Alice 给 Bob 发送一条消息。"
		// } else if i == 2 {
		// 	query = "向所有队友广播：“状态更新：第一阶段已完成。”"
		// } else if i == 3 {
		// 	query = "请查看lead收件箱，看是否有任何消息。"
		// } else if i == 4 {
		// 	query = "/team"
		// } else if i == 5 {
		// 	query = "/inbox"
		// } else {
		// 	query = "exit"
		// }
		fmt.Printf("\033[36ms01 >>%s \033[0m\n", query)
		query = strings.ToLower(strings.Trim(query, " "))
		if query == "q" || query == "exit" {
			break
		}
		if query == "/team" {
			fmt.Println(c.teamManager.ListAll())
			i++
			continue
		}
		if query == "/inbox" {
			fmt.Println(c.inbox.ReadInboxText("lead"))
			i++
			continue
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

func (c *MainAgent) AgentLoop(messages []api.Message) []api.Message {
	for true {
		// messages = MicroCompact(messages)
		// if EstimateTokens(messages) > c.threshold {
		// 	fmt.Println("[auto_compact triggered]")
		// 	messages = AutoCompact(c, messages)
		// }

		// notifs := background.MyBackgroundManager.DrainNotifications()
		// if notifs != nil && len(messages) != 0 {
		// 	lines := make([]string, 0)
		// 	for _, n := range notifs {
		// 		lines = append(lines, fmt.Sprintf("[bg:%s] %s: %s", n.ID, n.Status, n.Result))
		// 	}
		// 	notifText := strings.Join(lines, "\n")
		// 	messages = append(messages, api.Message{
		// 		Role:    "user",
		// 		Content: fmt.Sprintf("<background-results>\n%s\n</background-results>", notifText),
		// 	})
		// 	messages = append(messages, api.Message{
		// 		Role:    "assistant",
		// 		Content: "Noted background results.",
		// 	})
		// }

		msg := c.inbox.ReadInboxText("lead")
		if msg != "" && len(msg) != 0 {
			messages = append(messages, api.Message{
				Role:    "user",
				Content: fmt.Sprintf("<inbox>%s</inbox>", msg),
			})
			messages = append(messages, api.Message{
				Role:    "assistant",
				Content: "Noted inbox messages.",
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
