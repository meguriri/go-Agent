package task

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

type Task struct {
	ID          int    `json:"id"`
	Subject     string `json:"subject"`
	Description string `json:"description"`
	Status      string `json:"status"`
	BlockedBy   []int  `json:"blockedBy"`
	Blocks      []int  `json:"blocks"`
	Owner       string `json:"owner"`
}

type TaskManager struct {
	taskDir string
	nextID  int
}

var MyTaskManager *TaskManager = nil

func NewTaskManager(dir string) *TaskManager {
	os.MkdirAll(dir, 0755)
	t := &TaskManager{
		taskDir: dir,
		nextID:  1,
	}
	t.nextID = t.max_ID() + 1
	return t
}

func (t *TaskManager) max_ID() int {
	maxID := 0
	matches, err := filepath.Glob(filepath.Join(t.taskDir, "task_*.json"))
	if err != nil {
		return 0
	}
	for _, path := range matches {
		base := filepath.Base(path)
		stem := strings.TrimSuffix(base, ".json")
		parts := strings.Split(stem, "_")
		if len(parts) != 2 {
			continue
		}
		ID, err := strconv.Atoi(parts[1])
		if err != nil {
			continue
		}
		if ID > maxID {
			maxID = ID
		}
	}
	return maxID
}

func (t *TaskManager) load(ID int) Task {
	path := filepath.Join(t.taskDir, "task_"+strconv.Itoa(ID)+".json")
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatalln("read file error: " + err.Error())
		return Task{}
	}
	task := Task{}
	err = json.Unmarshal(data, &task)
	if err != nil {
		log.Fatalln("json unmarshal error: " + err.Error())
		return Task{}
	}
	return task
}

func (t *TaskManager) save(task Task) {
	path := filepath.Join(t.taskDir, "task_"+strconv.Itoa(task.ID)+".json")
	data, err := json.Marshal(task)
	if err != nil {
		log.Fatalln("json marshal error: " + err.Error())
	}
	err = os.WriteFile(path, data, 0644)
	if err != nil {
		log.Fatalln("write file error: " + err.Error())
	}
}

func (t *TaskManager) Create(Subject string, Description string) string {
	task := Task{
		ID:          t.nextID,
		Subject:     Subject,
		Description: Description,
		Status:      "pending",
		BlockedBy:   make([]int, 0),
		Blocks:      make([]int, 0),
		Owner:       "",
	}
	t.save(task)
	t.nextID++
	data, err := json.Marshal(task)
	if err != nil {
		log.Fatalln("json marshal error: " + err.Error())
	}
	return string(data)
}

func (t *TaskManager) Get(taskID int) string {
	task := t.load(taskID)
	data, err := json.Marshal(task)
	if err != nil {
		log.Fatalln("json marshal error: " + err.Error())
	}
	return string(data)
}

func (t *TaskManager) Update(ID int, Status string, addBlockBy []int, addBlock []int) string {
	task := t.load(ID)
	if Status != "" {
		if Status != "pending" && Status != "in_progress" && Status != "completed" {
			log.Fatalf("Item %d: invalID Status '%s'", ID, Status)
		}
		task.Status = Status
		if task.Status == "completed" {
			t.clearDependency(ID)
		}
	}
	if addBlockBy != nil || len(addBlockBy) != 0 {
		task.BlockedBy = append(task.BlockedBy, addBlockBy...)
	}
	if addBlock != nil || len(addBlock) != 0 {
		task.Blocks = append(task.Blocks, addBlock...)
		for _, blockId := range addBlock {
			blocked := t.load(blockId)
			find := false
			for _, ID := range blocked.BlockedBy {
				if ID == task.ID {
					find = true
					break
				}
			}
			if !find {
				blocked.BlockedBy = append(blocked.BlockedBy, ID)
				t.save(blocked)
			}
		}
	}
	t.save(task)
	data, err := json.Marshal(task)
	if err != nil {
		log.Fatalln("json marshal error: " + err.Error())
	}
	return string(data)
}

func (t *TaskManager) clearDependency(completed_ID int) {
	// Remove completed_ID from all other tasks' BlockedBy lists.
	matches, err := filepath.Glob(filepath.Join(t.taskDir, "task_*.json"))
	if err != nil {
		log.Fatalln("filepath glob error: " + err.Error())
	}
	for _, path := range matches {
		var task Task
		data, err := os.ReadFile(path)
		if err != nil {
			log.Fatalln("read file error: " + err.Error())
		}
		err = json.Unmarshal(data, &task)
		if err != nil {
			log.Fatalln("json unmarshal error: " + err.Error())
		}
		for i, ID := range task.BlockedBy {
			if ID == completed_ID {
				task.BlockedBy = append(task.BlockedBy[:i], task.BlockedBy[i+1:]...)
				t.save(task)
				break
			}
		}
	}
}

func (t *TaskManager) ListAll() string {
	tasks := make([]Task, 0)
	lines := make([]string, 0)
	matches, err := filepath.Glob(filepath.Join(t.taskDir, "task_*.json"))
	if err != nil {
		log.Fatalln("filepath glob error: " + err.Error())
		return ""
	}
	sort.Slice(matches, func(i, j int) bool {
		return matches[i] < matches[j]
	})
	for _, path := range matches {
		var task Task
		data, err := os.ReadFile(path)
		if err != nil {
			log.Fatalln("read file error: " + err.Error())
		}
		err = json.Unmarshal(data, &task)
		if err != nil {
			log.Fatalln("json unmarshal error: " + err.Error())
		}
		tasks = append(tasks, task)
	}
	if len(tasks) == 0 {
		return "No tasks."
	}
	for _, task := range tasks {
		var marker string
		switch task.Status {
		case "pending":
			marker = "[ ]"
		case "in_progress":
			marker = "[>]"
		case "completed":
			marker = "[x]"
		default:
			marker = "[?]"
		}
		var blocked string = ""
		if len(task.BlockedBy) != 0 {
			blocked = fmt.Sprintf(" (blocked by: %v)", task.BlockedBy)
		}
		lines = append(lines, fmt.Sprintf("%s #%d: %s%s", marker, task.ID, task.Subject, blocked))
	}
	return strings.Join(lines, "\n")
}
