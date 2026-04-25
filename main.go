package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"s01/internal/agent"
	"s01/internal/background"
	"s01/internal/team"
	"s01/internal/toolManager"

	"github.com/joho/godotenv"
	"github.com/ollama/ollama/api"
)

const (
	THRESHOLD      = 50000
	TRANSCRIPT_DIR = ".transcripts"
	TASKS_DIR      = ".tasks"
	TEAM_DIR       = ".team"
	INBOX_DIR      = ".inbox"
	SKILL_DIR      = "./skills"
	KEEP_RECENT    = 3
)

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
	bus := team.NewMessageBus(INBOX_DIR)
	tm := team.NewTeammateManager(TEAM_DIR, bus)
	subRunner := &agent.SubagentService{Client: c, Model: modelID}
	deps := toolManager.Deps{
		Messenger:      bus,
		InboxReader:    bus,
		SubagentRunner: subRunner,
		TeamManager:    tm,
		Broadcaster:    bus,
	}
	// skillLoader := skill.NewSkillLoader(SKILL_DIR)
	// task.MyTaskManager = task.NewTaskManager(TASKS_DIR)
	background.MyBackgroundManager = background.NewBackgroundManager()
	toolHandler := toolManager.NewToolHandler(deps)
	agent := agent.NewMainAgent(c, modelID, ctx, THRESHOLD,
		// fmt.Sprintf("你是一名位于 %s 的编程agent。在处理不熟悉的主题之前，请使用 load_skill 这个工具来获取专业知识。可用技能：%s", dir, skillLoader.GetDescriptions()),
		//fmt.Sprintf("你是目录：%s 的一名编程agent。你必须使用任务工具(工具名：task)来把任务委派给subagent。", dir),
		// fmt.Sprintf("You are a coding agent at %s. Use task tools to plan and track work. Please note that task tools are intended solely for planning and recording tasks; for the actual execution of a task, you must utilize other tools to the specific nature of the task—to carry it out to completion.", dir),
		// fmt.Sprintf("你是目录：%s 的一名编程agent。你必须使用工具来解决问题", dir),s
		// fmt.Sprintf("你是目录：%s 的一名编程agent。对长时间运行的命令使用 `background_run`。", dir+"/sandbox"),
		fmt.Sprintf("你是位于 %s 的团队负责人,你的名字是lead。生成teammate，并通过inbox进行沟通。你可以用spawn_teammate创建teamagent，然后别忘了写prompt，让他们去完成一个任务。别忘了teammate通过inbox进行沟通", dir),
		toolHandler,
		bus,
		tm,
	)
	agent.Chat()
}
