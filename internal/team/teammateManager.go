package team

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"s01/internal/agent"
	"s01/internal/toolManager"
	"strings"

	"github.com/ollama/ollama/api"
)

const (
	WORKDIR = "./"
)

type Member struct {
	Name   string `json:"name"`
	Role   string `json:"role"`
	Status string `json:"status"`
}

type Config struct {
	TeamName string    `json:"teamName"`
	Members  []*Member `json:"members"`
}

type TeammateManager struct {
	dir        string
	configPath string
	config     Config
	client     *api.Client
	bus        *MessageBus
}

func NewTeammateManager(team_dir string, bus *MessageBus) *TeammateManager {
	t := &TeammateManager{
		dir:        team_dir,
		configPath: team_dir + "/config.json",
		bus:        bus,
	}
	c, err := api.ClientFromEnvironment()
	if err != nil {
		log.Fatalf("Create client error: %v\n", err)
	}
	t.client = c
	os.MkdirAll(team_dir, 0755)
	t.config = *t.loadConfig()
	return t
}

func (t *TeammateManager) loadConfig() *Config {
	if _, err := os.Stat(t.configPath); err != nil && errors.Is(err, os.ErrNotExist) {
		log.Println("file not exist")
		return &Config{
			TeamName: "default",
			Members:  make([]*Member, 0),
		}
	}
	data, err := os.ReadFile(t.configPath)
	if err != nil {
		log.Fatalln("read file error: " + err.Error())
		return nil
	}
	config := &Config{}
	err = json.Unmarshal(data, config)
	if err != nil {
		log.Fatalln("json unmarshal error: " + err.Error())
		return nil
	}
	return config
}

func (t *TeammateManager) saveConfig() {
	data, err := json.Marshal(t.config)
	if err != nil {
		log.Fatalln("json marshal error: " + err.Error())
	}
	err = os.WriteFile(t.configPath, data, 0644)
	if err != nil {
		log.Fatalln("write file error: " + err.Error())
	}
}

func (t *TeammateManager) findMember(name string) *Member {
	for _, m := range t.config.Members {
		if m.Name == name {
			return m
		}
	}
	return nil
}

func (t *TeammateManager) Spawn(name string, role string, prompt string) string {
	member := t.findMember(name)
	if member != nil {
		if member.Status != "idle" && member.Status != "shutdown" {
			return fmt.Sprintf("Error: '%s' is currently %s", name, member.Status)
		}
		member.Status = "working"
		member.Role = role
	} else {
		member = &Member{
			Name:   name,
			Role:   role,
			Status: "working",
		}
		t.config.Members = append(t.config.Members, member)
	}

	t.saveConfig()
	go t.teammateLoop(name, role, prompt)
	return fmt.Sprintf("Spawned '%s' (role: %s)", name, role)
}

func (t *TeammateManager) teammateLoop(name string, role string, prompt string) {
	sys_prompt := fmt.Sprintf("You are '%s', role: %s, at %s. \nUse send_message to communicate. Complete your task.", name, role, WORKDIR)
	messages := []api.Message{
		api.Message{Role: "user", Content: prompt},
	}
	deps := toolManager.Deps{
		Messenger:   t.bus,
		InboxReader: t.bus,
		// SubagentRunner 对 team agent 可先不传
	}
	member := t.findMember(name)
	agent := agent.NewTeamAgent(name, role, "working", t.client, os.Getenv("OLLAMA_MODELS"), context.Background(), 5000, sys_prompt, toolManager.NewTeamToolHandler(deps), t.bus)
	agent.AgentLoop(messages)
	if member != nil && member.Status != "shutdown" {
		member.Status = "idle"
		t.saveConfig()
	}
}

func (t *TeammateManager) ListAll() string {
	if t.config.Members == nil || len(t.config.Members) == 0 {
		return "No teammates."
	}
	lines := []string{fmt.Sprintf("Team: %s", t.config.TeamName)}
	for _, m := range t.config.Members {
		lines = append(lines, fmt.Sprintf("  %s (%s): %s", m.Name, m.Role, m.Status))
	}
	return strings.Join(lines, "\n")
}

func (t *TeammateManager) MemberNames() []string {
	list := make([]string, len(t.config.Members))
	for i, m := range t.config.Members {
		list[i] = m.Name
	}
	return list
}
