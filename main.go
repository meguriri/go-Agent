package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"s01/model"
	"s01/task"
	"s01/tools"

	"github.com/joho/godotenv"
	"github.com/ollama/ollama/api"
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
	// skillLoader := skill.NewSkillLoader("./skills")
	task.MyTaskManager = task.NewTaskManager(".tasks")
	agent := model.NewAgent(c, modelID, ctx,
		// fmt.Sprintf("你是一名位于 %s 的编程agent。在处理不熟悉的主题之前，请使用 load_skill 这个工具来获取专业知识。可用技能：%s", dir, skillLoader.GetDescriptions()),
		//fmt.Sprintf("你是目录：%s 的一名编程agent。你必须使用任务工具(工具名：task)来把任务委派给subagent。", dir),
		fmt.Sprintf("You are a coding agent at %s. Use task tools to plan and track work. Please note that task tools are intended solely for planning and recording tasks; for the actual execution of a task, you must utilize other tools to the specific nature of the task—to carry it out to completion.", dir),
		// fmt.Sprintf("你是目录：%s 的一名编程agent。你必须使用工具来解决问题", dir),
		tools.NewToolHandler(),
	)
	agent.Chat()
}
