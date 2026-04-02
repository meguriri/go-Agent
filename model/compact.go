package model

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ollama/ollama/api"
)

const (
	TRANSCRIPT_DIR = ".transcripts"
	KEEP_RECENT    = 3
)

func EstimateTokens(messages []api.Message) int {
	//Rough token count: ~4 chars per token.
	tokens := 0
	for _, message := range messages {
		tokens += len(message.Content) + len(message.Thinking)
	}
	return tokens / 4
}

func MicroCompact(messages []api.Message) []api.Message {
	toolResults := make([]int, 0)
	for i, message := range messages {
		if message.Role == "tool" {
			toolResults = append(toolResults, i)
		}
	}
	if len(toolResults) <= KEEP_RECENT {
		return messages
	}
	toolResults = toolResults[:len(toolResults)-KEEP_RECENT]
	for _, idx := range toolResults {
		messages[idx].Content = fmt.Sprintf("[Previous: used %s]\n", messages[idx].ToolCalls[0].Function.Name)
	}
	return messages
}

func AutoCompact(c *Agent, messages []api.Message) []api.Message {
	os.MkdirAll(TRANSCRIPT_DIR, 0755)
	fileName := fmt.Sprintf("transcript_%d.jsonl", time.Now().Unix())
	transcriptPath := filepath.Join(TRANSCRIPT_DIR, fileName)
	f, err := os.Create(transcriptPath)
	if err != nil {
		return nil
	}
	defer f.Close()
	encoder := json.NewEncoder(f)
	for _, msg := range messages {
		encoder.Encode(msg) // Encode 会自动添加换行符
	}
	fmt.Printf("[transcript saved: %s]\n", transcriptPath)
	fullData, _ := json.Marshal(messages)
	conversationText := string(fullData)
	if len(conversationText) > 80000 {
		conversationText = conversationText[:80000]
	}

	req := &api.ChatRequest{
		Model: c.model,
		Messages: []api.Message{{
			Role:    "user",
			Content: "Summarize this conversation for continuity. Include: 1) What was accomplished, 2) Current state, 3) Key decisions made. Be concise but preserve critical details.\n\n" + conversationText},
		},
	}
	var fullContent strings.Builder
	var assistantMsg api.Message

	err = c.client.Chat(c.ctx, req, func(resp api.ChatResponse) error {
		fullContent.WriteString(resp.Message.Content)
		return nil
	})
	if err != nil {
		log.Fatalf("get llm response error: %v\n", err)
		return nil
	}
	assistantMsg.Role = "assistant"
	assistantMsg.Content = fullContent.String()
	summary := assistantMsg.Content
	return []api.Message{
		{Role: "user", Content: fmt.Sprintf("[Conversation compressed. Transcript: %s]\n\n%s", transcriptPath, summary)},
		{Role: "assistant", Content: "Understood. I have the context from the summary. Continuing."},
	}
}
