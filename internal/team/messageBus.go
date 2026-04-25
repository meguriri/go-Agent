package team

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ollama/ollama/api"
)

var VALID_MSG_TYPES = map[string]struct{}{
	"message":                struct{}{},
	"broadcast":              struct{}{},
	"shutdown_request":       struct{}{},
	"shutdown_response":      struct{}{},
	"plan_approval_response": struct{}{},
}

const (
	TEAM_DIR = "./.team"
)

type Message struct {
	Type      string
	From      string
	Content   string
	Timestamp time.Time
}

type MessageBus struct {
	Dir string
}

func NewMessageBus(dir string) *MessageBus {
	os.MkdirAll(dir, 0755)
	b := &MessageBus{
		Dir: dir,
	}
	return b
}

func (b *MessageBus) Send(sender string, to string,
	content string, msgtype string, extra map[string]any) string {
	if _, ok := VALID_MSG_TYPES[msgtype]; !ok {
		return fmt.Sprintf("Error: Invalid type %s. Valid: %v", msgtype, VALID_MSG_TYPES)
	}
	msg := Message{
		Type:      msgtype,
		From:      sender,
		Content:   content,
		Timestamp: time.Now(),
	}
	if extra != nil && len(extra) != 0 {
		for k, v := range extra {
			if k == "Type" {
				msg.Type = v.(string)
			} else if k == "From" {
				msg.From = v.(string)
			} else if k == "Content" {
				msg.Content = v.(string)
			} else if k == "Timestamp" {
				msg.Timestamp = v.(time.Time)
			}
		}
	}
	inboxPath := filepath.Join(b.Dir, fmt.Sprintf("%s.jsonl", to))
	f, err := os.OpenFile(inboxPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalln("read file error: " + err.Error())
		return ""
	}
	defer f.Close()
	data, err := json.Marshal(msg)
	_, err = f.WriteString(string(data) + "\n")

	return fmt.Sprintf("Sent %s to %s", msgtype, to)
}

func (b *MessageBus) ReadInboxText(name string) string {
	messages := b.ReadInbox(name)
	lines := make([]string, 0)
	for _, msg := range messages {
		data, err := json.Marshal(msg)
		if err != nil {
			continue
		}
		lines = append(lines, string(data))
	}
	return strings.Join(lines, "\n")
}

func (b *MessageBus) ReadInboxMessages(name string) []api.Message {
	msgs := b.ReadInbox(name)
	messages := make([]api.Message, len(msgs))
	for i, msg := range msgs {
		data, err := json.Marshal(msg)
		if err != nil {
			log.Fatalln("json marshal error: " + err.Error())
			continue
		}
		messages[i] = api.Message{
			Role:    "user",
			Content: string(data),
		}
	}
	return messages
}

func (b *MessageBus) ReadInbox(name string) []Message {
	inboxPath := filepath.Join(b.Dir, fmt.Sprintf("%s.jsonl", name))
	_, err := os.Stat(inboxPath)
	if os.IsNotExist(err) {
		return nil
	}
	messages := make([]Message, 0)
	file, err := os.Open(inboxPath)
	if err != nil {
		log.Fatalln("read file error: " + err.Error())
		return nil
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}
		var msg Message
		if err := json.Unmarshal(line, &msg); err != nil {
			fmt.Printf("JSON 解析失败: %v | 内容: %s\n", err, string(line))
			continue
		} else {
			messages = append(messages, msg)
		}
	}
	err = os.Truncate(inboxPath, 0)
	if err != nil {
		log.Fatalln("清空inbox失败: " + err.Error())
		return nil
	}
	return messages
}

func (b *MessageBus) Broadcast(sender string, content string, teamMates []string) string {
	count := 0
	extra := make(map[string]any)
	for _, name := range teamMates {
		if name != sender {
			b.Send(sender, name, content, "broadcast", extra)
			count++
		}
	}
	return fmt.Sprintf("Broadcast to %d teammates", count)
}
