package tool

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/ollama/ollama/api"
)

type bashTool api.Tool

func (b bashTool) GetTool() api.Tool {
	props := api.NewToolPropertiesMap()
	props.Set("command", api.ToolProperty{
		Type:        api.PropertyType{"string"},
		Description: "shell command",
	})

	bashToolFunction := api.ToolFunction{
		Name:        "bash",
		Description: "Run a shell command.",
		Parameters: api.ToolFunctionParameters{
			Type:       "object",
			Required:   []string{"command"},
			Properties: props,
		},
	}
	return api.Tool{
		Type:     "function",
		Function: bashToolFunction,
	}
}

func (b *bashTool) Run(args api.ToolCallFunctionArguments) string {
	dangerous := []string{
		"rm -rf .", "sudo",
		"shutdown", "reboot", "> /dev/",
	}
	lowerCmd := strings.ToLower(command)
	for _, d := range dangerous {
		if strings.Contains(lowerCmd, d) {
			return "Error: Dangerous command blocked"
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "bash", "-c", command)
	dir, _ := os.Getwd()
	cmd.Dir = dir

	out, err := cmd.CombinedOutput()
	if ctx.Err() == context.DeadlineExceeded {
		return "Error: Timeout (120s)"
	}
	result := strings.TrimSpace(string(out))
	if result == "" {
		if err != nil {
			return fmt.Sprintf("exec error: %v", err)
		}
		return "(no output)"
	}

	if len(result) > 50000 {
		return result[:50000]
	}
	return result
}
