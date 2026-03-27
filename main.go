package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"s01/tool"
	"strings"

	"github.com/joho/godotenv"

	"github.com/ollama/ollama/api"
)

type Client struct {
	client      *api.Client
	model       string
	ctx         context.Context
	System      string
	tools       api.Tools
	ToolHandler tool.ToolHandler
}

func (c *Client) agentLoop(messages []api.Message) []api.Message {
	for true {
		req := &api.ChatRequest{
			Model:    c.model,
			Messages: messages,
			Tools:    c.tools,
		}
		var fullContent strings.Builder
		var thinkingContent strings.Builder
		var assistantMsg api.Message
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
		if len(assistantMsg.ToolCalls) == 0 {
			return messages
		}
		if assistantMsg.Thinking != "" {
			fmt.Printf("正在思考：%s\n", assistantMsg.Thinking)
		}
		for _, tc := range assistantMsg.ToolCalls {
			fmt.Printf("\033[33m$ 正在执行工具: %s\033[0m\n", tc.Function.Name)
			var output string
			handler, ok := c.ToolHandler[tc.Function.Name]
			if !ok {
				output = "Unknown tool: " + tc.Function.Name
			} else {
				output = handler.Run(tc.Function.Arguments)
			}

			fmt.Printf("执行结果摘要: %s\n", strings.Split(output, "\n")[0])
			toolResultMsg := api.Message{
				Role:    "tool",
				Content: output,
			}
			messages = append(messages, toolResultMsg)
		}
	}
	return messages
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("读取 .env 文件失败，请检查文件是否存在")
	}
	ollamaHost := os.Getenv("OLLAMA_HOST")
	modelID := os.Getenv("OLLAMA_MODELS")
	fmt.Printf("正在连接 Ollama: %s，使用模型: %s\n", ollamaHost, modelID)

	c, err := api.ClientFromEnvironment()
	if err != nil {
		log.Fatalf("Create client error: %v\n", err)
	}
	ctx := context.Background()
	dir, _ := os.Getwd()
	client := Client{
		client:      c,
		model:       "qwen3:8b",
		ctx:         ctx,
		System:      fmt.Sprintf("You are a coding agent at %s. Use bash to solve tasks. Act, don't explain.", dir),
		tools:       nil,
		ToolHandler: tool.NewToolHandler(),
	}
	for _, v := range client.ToolHandler {
		client.tools = append(client.tools, v.GetTool())
	}

	history := make([]api.Message, 0)
	history = append(history, api.Message{
		Role:    "system",
		Content: client.System,
	})
	var query string
	for true {
		fmt.Print("\033[36ms01 >> \033[0m")
		// fmt.Scan(&query)
		query = "我的test目录下有一个hello.py，然后它是从小到大的排序，我想给他改成从大到小的排序"
		query = strings.ToLower(strings.Trim(query, " "))
		if query == "q" || query == "exit" {
			break
		}
		history = append(history, api.Message{
			Role:    "user",
			Content: query,
		})
		history = client.agentLoop(history)
		responses_content := history[len(history)-1].Content
		fmt.Println(responses_content)
	}

}
