package repository

import (
	"mymodule/internal/task/model"
	"mymodule/internal/task/usecase"
	"mymodule/pkg/logger"
	"time"

	"gorm.io/gorm"
)

type GormTaskRepository struct {
	db *gorm.DB
}

func NewGormTaskRepository(db *gorm.DB) usecase.TaskRepository {
	return &GormTaskRepository{db: db}
}

func (r *GormTaskRepository) Save(task model.Task) error {
	if err := r.db.Create(&task).Error; err != nil {
		logger.LogTask(task).Error("Failed to save task")
		return err
	}
	logger.LogTask(task).Info("Task saved successfully")
	return nil
}

func (r *GormTaskRepository) FindByID(taskID uint) (*model.Task, error) {
	var task model.Task
	if err := r.db.First(&task, taskID).Error; err != nil {
		logger.Log.WithField("taskID", taskID).Error("Failed to find task by ID")
		return nil, err
	}
	logger.Log.WithField("taskID", taskID).Info("Task found by ID")
	return &task, nil
}

func (r *GormTaskRepository) FindByUser(userID uint) (*[]model.Task, error) {
	var tasks []model.Task
	if err := r.db.Where("user_id = ?", userID).Find(&tasks).Error; err != nil {
		logger.Log.WithField("userID", userID).Error("Failed to find tasks by user ID")
		return nil, err
	}
	logger.Log.WithField("userID", userID).Info("Tasks found by user ID")
	return &tasks, nil
}

func (r *GormTaskRepository) FindByIDAndUser(taskID, userID uint) (*model.Task, error) {
	var task model.Task
	if err := r.db.Where("id = ? AND user_id = ?", taskID, userID).First(&task).Error; err != nil {
		logger.Log.WithFields(logger.LogFields(taskID, userID)).Error("Failed to find task by ID and user ID")
		return nil, err
	}
	logger.Log.WithFields(logger.LogFields(taskID, userID)).Info("Task found by ID and user ID")
	return &task, nil
}

func (r *GormTaskRepository) Update(task *model.Task) error {
	if err := r.db.Save(task).Error; err != nil {
		logger.LogTask(*task).Error("Failed to update task")
		return err
	}
	logger.LogTask(*task).Info("Task updated successfully")
	return nil
}

func (r *GormTaskRepository) Delete(taskID uint) error {
	if err := r.db.Delete(&model.Task{}, taskID).Error; err != nil {
		logger.Log.WithField("taskID", taskID).Error("Failed to delete task")
		return err
	}
	logger.Log.WithField("taskID", taskID).Info("Task deleted successfully")
	return nil
}

func (r *GormTaskRepository) UpdateOverdueTasks(userID uint) error {
	now := time.Now().UTC()
	return r.db.Model(&model.Task{}).
         Where("user_id = ? AND due_date <= ? AND status NOT IN ?", userID, now, []string{"completed","overdue"}).
		Update("status", "overdue").Error

}