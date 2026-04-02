package background

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

var MyBackgroundManager *BackgroundManager = nil

type Task struct {
	ID      string
	Command string
	Status  string
	Result  string
}

type BackgroundManager struct {
	Tasks   map[string]*Task
	notifCh chan *Task
	mu      sync.Mutex
}

func NewBackgroundManager() *BackgroundManager {
	return &BackgroundManager{
		Tasks:   make(map[string]*Task),
		notifCh: make(chan *Task, 100),
	}
}

func (b *BackgroundManager) Run(command string) string {
	// Start a background thread, return task_id immediately.
	taskID := uuid.New().String()[:8]
	task := &Task{
		Status:  "running",
		Result:  "",
		Command: command,
	}
	b.mu.Lock()
	b.Tasks[taskID] = task
	b.mu.Unlock()

	go b.execute(taskID, command)
	if len(command) > 80 {
		command = command[:80]
	}
	return fmt.Sprintf("Background task {task_id} started: %s", command)
}

func (b *BackgroundManager) execute(taskID string, command string) {
	// Thread target: run subprocess, capture output, push to queue.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	fmt.Println("start a goroutine!!!!!!!!")
	cmd := exec.CommandContext(ctx, "sh", "-c", command)
	output, err := cmd.CombinedOutput()

	b.mu.Lock()
	defer b.mu.Unlock()

	task := b.Tasks[taskID]

	if ctx.Err() == context.DeadlineExceeded {
		task.Status = "timeout"
		task.Result = "Error: Timeout (300s)"
	} else if err != nil {
		task.Status = "error"
		task.Result = fmt.Sprintf("Error: %v\nOutput: %s", err, string(output))
	} else {
		task.Status = "completed"
		task.Result = "(no output)"
		if output != nil && len(output) != 0 {
			task.Result = string(output)
		}
	}

	select {
	case b.notifCh <- task:
		fmt.Println("task is complete", *task)
	default: // 如果管道满了，防止阻塞后台协程
	}
}

func (b *BackgroundManager) Check(taskID string) string {
	// Check status of one task or list all.
	if taskID != "" {
		if task, ok := b.Tasks[taskID]; ok {
			var result string = "(running)"
			if task.Result != "" {
				result = task.Result
			}
			if len(task.Command) > 60 {
				task.Command = task.Command[:60]
			}
			return fmt.Sprintf("[%s] %s\n%s", task.Status, task.Command, result)
		} else {
			return fmt.Sprintf("Error: Unknown task %s", taskID)
		}
	}
	lines := make([]string, 0)
	for tid, task := range b.Tasks {
		var command string = task.Command
		if len(command) > 60 {
			command = command[:60]
		}
		lines = append(lines, fmt.Sprintf("%s: [%s] %s", tid, task.Status, command))
	}
	if len(lines) == 0 {
		return "No background tasks."
	}
	return strings.Join(lines, "\n")
}

func (b *BackgroundManager) DrainNotifications() []*Task {
	// Return and clear all pending completion notifications.
	var results []*Task
	for {
		select {
		case t := <-b.notifCh:
			fmt.Println("task is complete send signal", t.ID)
			results = append(results, t)
		default:
			// 管道空了，立即返回已收集到的结果
			return results
		}
	}
}
