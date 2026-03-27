package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"s01/model"
	"s01/tool"

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

	client := model.NewClient(c, modelID, ctx,
		fmt.Sprintf(`You are a coding agent at %s.You have to use the todo tool to plan multi-step tasks.
		 Mark in_progress before starting, completed when done.Prefer tools over prose.`, dir),
		tool.NewToolHandler(),
	)
	client.Chat()
}
