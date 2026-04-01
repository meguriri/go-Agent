package task

type Task struct {
	status string
}

type TaskManager struct {
	taskDir string
	nextID  int
}

func NewTaskManager(dir string) *TaskManager {
	t := &TaskManager{
		taskDir: dir,
		nextID:  1,
	}
	return t
}

func (t *TaskManager) Create(description string) int {
	task := Task{}
	t.save(task)
	t.nextID++
	return t.nextID
}

func (t *TaskManager) save(task Task) {

}

func (t *TaskManager) load(id int) Task {
	return Task{}
}

func (t *TaskManager) ClearDependency(id int) {

}

func (t *TaskManager) Update(id int, status string, addBlockBy int, addBlock int) {
	task := t.load(id)
	if status != "" {
		task.status = status
		if task.status == "complete" {
			t.ClearDependency(id)
		}
	}
}
