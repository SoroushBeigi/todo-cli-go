package memorystore

import (
	"github.com/SoroushBeigi/todo-cli-go/entity"
)

type Task struct {
	tasks []entity.Task
}

func NewTaskStore() *Task {
	return &Task{
		tasks: make([]entity.Task, 0),
	}
}

func (t *Task) CreateNewTask(task entity.Task) (entity.Task, error) {
	task.ID = len(t.tasks) + 1

	t.tasks = append(t.tasks, task)

	return task, nil
}

func (t *Task) List(userID int) ([]entity.Task, error) {

	var tasks []entity.Task
	for _, tsk := range t.tasks {
		if tsk.UserID == userID {
			tasks = append(tasks, tsk)
		}
	}

	return tasks, nil
}
