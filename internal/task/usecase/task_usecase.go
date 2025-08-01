package usecase

import (
	"mymodule/internal/task/model"

)

type TaskRepository interface {
	Save(task model.Task) error
	FindByID(taskID uint) (*model.Task, error)
	FindByUser(userName string) (*model.Task, error)
	Update(task model.Task) error
	Delete(taskID uint) error
}

type TaskUsecase interface {
	Create(task model.Task) error
	GetByID(taskID uint) (string, error)
	GetByUser(userName string) (string, error)
	Updatetask(task model.Task) error
	Deletetask(taskID uint) error
}

type TaskusecaseImpl struct {
	repo  TaskRepository
}

func NewtaskUsecase(repo TaskRepository) TaskUsecase {
	return &TaskusecaseImpl{
		repo:  repo,
	}
}


func (uc *TaskusecaseImpl) Create(task model.Task) error {
	return  nil
}
func (uc *TaskusecaseImpl) GetByID(taskID uint) (string, error) {
		return  "",nil
}
func (uc *TaskusecaseImpl) GetByUser(userName string) (string, error) {
		return  "",nil
}
func (uc *TaskusecaseImpl) Updatetask(task model.Task)  error {
		return  nil
}
func (uc *TaskusecaseImpl) Deletetask(taskID uint) error {
	return nil
}
