package task

import (
	"fmt"

	"github.com/SoroushBeigi/todo-cli-go/entity"
)

type TaskRepository interface {
	// CanCreateTaskInCategory(userID, categoryID int) (bool, error)
	CreateNewTask(t entity.Task) (entity.Task, error)
	List(userID int) ([]entity.Task, error)
}

type Service struct {
	repository TaskRepository
}

func NewService(repo TaskRepository) Service{
	return Service{repository: repo}
}

type CreateRequest struct {
	UserID     int
	Title      string
	DueDate    string
	CategoryID int
}

type CreateResponse struct {
	Task entity.Task
}

func (s Service) CreateNewTask(req CreateRequest) (CreateResponse, error) {

	// ok, err := s.repository.CanCreateTaskInCategory(req.UserID, req.CategoryID)
	// if !ok {

	// 	return CreateResponse{}, fmt.Errorf("user %d does not have category %d", req.UserID, req.CategoryID)
	// }


	createdTask, err := s.repository.CreateNewTask(entity.Task{
		ID:         0,
		Title:      req.Title,
		DueDate:    req.DueDate,
		CategoryID: req.CategoryID,
		UserID:     req.UserID,
		IsDone:     false,
	})

	if err != nil {

		return CreateResponse{}, fmt.Errorf("can't create new task: %v", err)
	}

	return CreateResponse{Task: createdTask}, nil

}

type ListResponse struct {
	Tasks []entity.Task
}

type ListRequest struct {
	UserID int
}

func (s Service) List(req ListRequest) (ListResponse, error) {
	tasks, err := s.repository.List(req.UserID)
	if err != nil {
		return ListResponse{}, fmt.Errorf("Cannot link task: %v", err)
	}
	return ListResponse{Tasks: tasks}, nil
}
