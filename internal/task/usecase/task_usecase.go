package usecase

import (
	"mymodule/internal/task/model"

)

type TaskRepository interface {
	Save(task model.Task) error
	FindByID(taskID uint) (*model.Task, error)
	FindByUser(userID uint) (*[]model.Task, error)
	FindByIDAndUser(taskID , userID uint) (*model.Task, error)
	Update(task model.Task) error
	Delete(taskID uint) error
}

type TaskUsecase interface {
	Create(task model.Task) error
	GetByID(taskID uint) (*model.Task, error)
	GetByUser(userID uint) (*[]model.Task, error)
	GetByIDAndUser(taskID , userID uint) (*model.Task, error)
	UpdateTask(task model.Task) error
	DeleteTask(taskID uint) error
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
func (uc *TaskusecaseImpl) GetByID(taskID uint) (*model.Task, error) {
		return  &model.Task{},nil
}
func (uc *TaskusecaseImpl) GetByUser(userID uint) (*[]model.Task, error) {
		return  &[]model.Task{},nil
}
func (uc *TaskusecaseImpl) GetByIDAndUser(taskID , userID uint) (*model.Task, error) {
		return  &model.Task{},nil
}
func (uc *TaskusecaseImpl) UpdateTask(task model.Task)  error {
		return  nil
}
func (uc *TaskusecaseImpl) DeleteTask(taskID uint) error {
	return nil
}
