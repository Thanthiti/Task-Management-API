package repository

import (
	"mymodule/internal/task/model"
	"mymodule/internal/task/usecase"

	"gorm.io/gorm"
)

type GormTaskRepository struct {
	db *gorm.DB
}

func NewGormTaskRepository(db *gorm.DB) usecase.TaskRepository {
	return &GormTaskRepository{db: db}
}

func (r *GormTaskRepository) Save(Task model.Task) error {
return  nil
}

func (r *GormTaskRepository) FindByID(taskID uint) (*model.Task, error) {
return  &model.Task{},nil
}


func (r *GormTaskRepository) FindByUser(userName string) (*model.Task, error) {
	return  &model.Task{},nil
}

func (r *GormTaskRepository) Update(Task model.Task) error {
return  nil
}

func (r *GormTaskRepository) Delete(TaskID uint) error {
	return nil
}
