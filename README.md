## go-Agent

一个基于 Ollama 的 Go 版编程 Agent 示例项目。  
项目采用「主 Agent + 子 Agent」架构：
- 主 Agent 负责理解用户目标，并通过 `task` 工具把具体任务委派给子 Agent。
- 子 Agent 具备读写文件、编辑文件、执行命令、管理待办等能力，完成后返回总结给主 Agent。

---

## 功能概览

- 基于 `github.com/ollama/ollama/api` 对接本地或远程 Ollama。
- 支持工具调用（Tool Calling）。
- 支持子代理隔离执行（新上下文，不共享主对话历史）。
- 支持从 `skills/` 加载技能文档（`SKILL.md` + frontmatter）。

---

## 环境要求

- Go 1.26.1（见 `go.mod`）
- 可用的 Ollama 服务
- 已下载可用模型（例如 `qwen3:8b`、`llama3.1` 等）

---

## 快速开始

### 1. 安装依赖

在项目根目录执行：

```bash
go mod tidy
```

### 2. 配置环境变量

在项目根目录创建 `.env` 文件：

```env
OLLAMA_HOST=http://127.0.0.1:11434
OLLAMA_MODELS=qwen3:8b
```

说明：
- `OLLAMA_HOST`：Ollama 服务地址。
- `OLLAMA_MODELS`：运行时使用的模型名。

### 3. 启动 Ollama（如未启动）

```bash
ollama serve
```

如模型尚未拉取，可执行：

```bash
ollama pull qwen3:8b
```

### 4. 运行 Agent

```bash
go run .
```

---


