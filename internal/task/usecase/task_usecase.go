package usecase

import (
	"errors"
	"fmt"
	"mymodule/internal/task/model"
	"mymodule/pkg/logger"
	"time"

	"gorm.io/gorm"
)

type TaskRepository interface {
	Save(task model.Task) error
	FindByID(taskID uint) (*model.Task, error)
	FindByUser(userID uint) (*[]model.Task, error)
	FindByIDAndUser(taskID, userID uint) (*model.Task, error)
	Update(task *model.Task) error
	Delete(taskID uint) error
	UpdateOverdueTasks(userID uint) error
}

type TaskUsecase interface {
	Create(task model.Task) error
	GetByID(taskID uint) (*model.Task, error)
	GetByUser(userID uint) (*[]model.Task, error)
	GetByIDAndUser(taskID, userID uint) (*model.Task, error)
	UpdateTask(task *model.UpdateTaskInput, taskID, userID uint) error
	DeleteTask(taskID, userID uint) error
}

type TaskusecaseImpl struct {
	repo TaskRepository
}

func NewTaskUsecase(repo TaskRepository) TaskUsecase {
	return &TaskusecaseImpl{
		repo: repo,
	}
}

func (uc *TaskusecaseImpl) SetStatusBasedOnDueDate(task *model.Task) {
	if task.DueDate != nil && task.DueDate.Before(time.Now()) {
		task.Status = "overdue"
	} else if task.Status == "" {
		task.Status = "pending"
	}
}

func (uc *TaskusecaseImpl) Create(task model.Task) error {
	uc.SetStatusBasedOnDueDate(&task)

	// Valid duedate
	if task.DueDate != nil && task.DueDate.Before(time.Now()) {
		return errors.New("invalid due date")
	}

	if err := uc.repo.Save(task); err != nil {
		logger.Log.WithField("userID", task.UserID).Error("Failed to create task")
		return err
	}

	logger.Log.WithField("userID", task.UserID).Info("Task created successfully")
	return nil

}

// For admin
func (uc *TaskusecaseImpl) GetByID(taskID uint) (*model.Task, error) {
	// task, err := uc.repo.FindByID(taskID)
	// if err != nil {
	// 	if errors.Is(err, gorm.ErrRecordNotFound) {
	// 		logger.Log.WithField("taskID", taskID).Warn("Task not found")
	// 		return nil, fmt.Errorf("task not found")
	// 	}
	// 	logger.Log.WithField("taskID", taskID).Error("Failed to get task by ID")
	// 	return nil, err
	// }
	// logger.Log.WithField("taskID", taskID).Info("Task retrieved by ID")
	// return task, nil
	return &model.Task{}, nil
}

func (uc *TaskusecaseImpl) GetByUser(userID uint) (*[]model.Task, error) {

	if err := uc.repo.UpdateOverdueTasks(userID); err != nil {
		logger.Log.WithField("userID", userID).Error("Failed to update overdue tasks")
		return nil, err
	}

	tasks, err := uc.repo.FindByUser(userID)

	if err != nil {
		logger.Log.WithField("userID", userID).Error("Failed to get tasks by user")
		return nil, err
	}

	logger.Log.WithField("userID", userID).Info("Tasks retrieved for user")
	return tasks, nil
}

func (uc *TaskusecaseImpl) GetByIDAndUser(taskID, userID uint) (*model.Task, error) {
	task, err := uc.repo.FindByIDAndUser(taskID, userID)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {

			logger.Log.WithFields(logger.LogFields(taskID, userID)).Warn("Task not found for this user")
			return nil, fmt.Errorf("task not found for user")
		}
		logger.Log.WithFields(logger.LogFields(taskID, userID)).Error("Failed to get task by ID and user")
		return nil, err
	}

	logger.Log.WithFields(logger.LogFields(taskID, userID)).Info("Task retrieved by ID and user")
	return task, nil
}

func (uc *TaskusecaseImpl) UpdateTask(input *model.UpdateTaskInput, taskID, userID uint) error {
    existingTask, err := uc.repo.FindByIDAndUser(taskID, userID)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            logger.Log.WithField("taskID", taskID).Warn("Update failed: task not found")
            return fmt.Errorf("task not found")
        }
        logger.Log.WithField("taskID", taskID).Error("Database error when checking task existence")
        return err
    }

    model.ApplyUpdate(existingTask, *input)
    uc.SetStatusBasedOnDueDate(existingTask)

    if err := uc.repo.Update(existingTask); err != nil {
        logger.Log.WithField("taskID", existingTask.ID).Error("Failed to update task")
        return err
    }

    logger.Log.WithField("taskID", existingTask.ID).Info("Task updated successfully")
    return nil
}


func (uc *TaskusecaseImpl) DeleteTask(taskID, userID uint) error {
	task, err := uc.repo.FindByIDAndUser(taskID, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Log.Warn("Delete failed: task not found or unauthorized")
			return fmt.Errorf("task not found or unauthorized")
		}
		logger.Log.Error("DB error when finding task by ID and userID: ", err)
		return err
	}

	if err := uc.repo.Delete(task.ID); err != nil {
		logger.Log.Error("Delete failed : ", err)
		return err
	}

	logger.Log.Info("Task deleted : ", taskID)
	return nil
}
