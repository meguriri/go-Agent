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
	id          int
	subject     string
	description string
	status      string
	blockedBy   []int
	blocks      []int
	owner       string
}

type TaskManager struct {
	taskDir string
	nextID  int
}

func NewTaskManager(dir string) *TaskManager {
	os.MkdirAll(dir, 0755)
	t := &TaskManager{
		taskDir: dir,
		nextID:  1,
	}
	t.nextID = t.max_id() + 1
	return t
}

func (t *TaskManager) max_id() int {
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
		id, err := strconv.Atoi(parts[1])
		if err != nil {
			continue
		}
		if id > maxID {
			maxID = id
		}
	}
	return maxID
}

func (t *TaskManager) load(id int) Task {
	path := filepath.Join(t.taskDir, "task_"+strconv.Itoa(id)+".json")
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
	path := filepath.Join(t.taskDir, "task_"+strconv.Itoa(task.id)+".json")
	data, err := json.Marshal(task)
	if err != nil {
		log.Fatalln("json marshal error: " + err.Error())
	}
	err = os.WriteFile(path, data, 0644)
	if err != nil {
		log.Fatalln("write file error: " + err.Error())
	}
}

func (t *TaskManager) Create(subject string, description string) string {
	task := Task{
		id:          t.nextID,
		subject:     subject,
		description: description,
		status:      "pending",
		blockedBy:   make([]int, 0),
		blocks:      make([]int, 0),
		owner:       "",
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

func (t *TaskManager) Update(id int, status string, addBlockBy []int, addBlock []int) string {
	task := t.load(id)
	if status != "" {
		if status != "pending" && status != "in_progress" && status != "completed" {
			log.Fatalf("Item %d: invalid status '%s'", id, status)
		}
		task.status = status
		if task.status == "completed" {
			t.clearDependency(id)
		}
	}
	if addBlockBy != nil || len(addBlockBy) != 0 {
		task.blockedBy = append(task.blockedBy, addBlockBy...)
	}
	if addBlock != nil || len(addBlock) != 0 {
		task.blocks = append(task.blocks, addBlock...)
		for _, blockId := range addBlock {
			blocked := t.load(blockId)
			find := false
			for _, id := range blocked.blockedBy {
				if id == task.id {
					find = true
					break
				}
			}
			if !find {
				blocked.blockedBy = append(blocked.blockedBy, id)
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

func (t *TaskManager) clearDependency(completed_id int) {
	// Remove completed_id from all other tasks' blockedBy lists.
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
		for i, id := range task.blockedBy {
			if id == completed_id {
				task.blockedBy = append(task.blockedBy[:i], task.blockedBy[i+1:]...)
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
		switch task.status {
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
		if len(task.blockedBy) != 0 {
			blocked = fmt.Sprintf(" (blocked by: %v)", task.blockedBy)
		}
		lines = append(lines, fmt.Sprintf("%s #%d: %s%s", marker, task.id, task.subject, blocked))
	}
	return strings.Join(lines, "\n")
}
