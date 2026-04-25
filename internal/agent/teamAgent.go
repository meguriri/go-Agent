package agent

import (
	"context"
	"fmt"
	"log"
	"s01/internal/toolManager"
	"strings"

	"github.com/ollama/ollama/api"
)

type TeamAgent struct {
	name              string
	role              string
	status            string
	client            *api.Client
	model             string
	ctx               context.Context
	threshold         int
	System            string
	tools             api.Tools
	ToolHandler       toolManager.ToolHandler
	rounds_since_todo int
	inbox             Inbox
}

func NewTeamAgent(name string, role string, status string, client *api.Client, model string, ctx context.Context, threshold int, system string,
	toolHandler toolManager.ToolHandler, inbox Inbox) *TeamAgent {
	c := &TeamAgent{
		name:              name,
		role:              role,
		status:            status,
		client:            client,
		model:             model,
		ctx:               ctx,
		threshold:         threshold,
		System:            system,
		tools:             nil,
		ToolHandler:       toolHandler,
		rounds_since_todo: 0,
		inbox:             inbox,
	}
	for _, v := range c.ToolHandler {
		c.tools = append(c.tools, v.GetTool())
	}
	return c
}

func (c TeamAgent) Model() string {
	return c.model
}
func (c *TeamAgent) Client() *api.Client {
	return c.client
}
func (c *TeamAgent) Ctx() context.Context {
	return c.ctx
}

func (c *TeamAgent) AgentLoop(messages []api.Message) []api.Message {
	for i := 0; i < 50; i++ {
		inbox := c.inbox.ReadInboxMessages(c.name)
		// for _, msg := range inbox {
		// 	data, err := json.Marshal(msg)
		// 	if err != nil {
		// 		log.Fatalln("json marshal error: " + err.Error())
		// 	}
		// 	messages = append(messages, api.Message{
		// 		Role:    "user",
		// 		Content: string(data),
		// 	})
		// }
		messages = append(messages, inbox...)
		req := &api.ChatRequest{
			Model:    c.model,
			Messages: messages,
			Tools:    c.tools,
		}

		var assistantMsg api.Message
		first_thinking := 0
		err := c.client.Chat(c.ctx, req, func(resp api.ChatResponse) error {
			if resp.Message.Thinking != "" {

				assistantMsg.Thinking += resp.Message.Thinking
				if first_thinking == 0 {
					fmt.Printf("\033[35m%s:\033[0m\033[35m正在思考：\033[0m", c.name)
				}
				fmt.Printf("\033[35m%s\033[0m", resp.Message.Thinking)
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
			fmt.Printf("\033[35m%s $ 正在执行工具: %s\033[0m\n", c.name, tc.Function.Name)
			var output string
			handler, ok := c.ToolHandler[tc.Function.Name]
			if !ok {
				output = "Unknown tool: " + tc.Function.Name
			} else {
				output = handler.Run(tc.Function.Arguments)
			}
			if tc.Function.Name != "todo" {
				fmt.Printf("\033[35m%s 执行结果：%s\033[0m\n", c.name, strings.Split(output, "\n")[0])
			}
			toolResultMsg := api.Message{
				Role:      "tool",
				Content:   output,
				ToolCalls: []api.ToolCall{tc},
			}
			messages = append(messages, toolResultMsg)
		}
	}
	return messages
}
