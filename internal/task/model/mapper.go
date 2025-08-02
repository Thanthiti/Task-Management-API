package model

import "github.com/jinzhu/copier"

func ToTask(req CreateTaskRequest, userID uint) Task {
	return Task{
		Title:       req.Title,
		Description: req.Description,
		DueDate:     req.DueDate,
		Status:      "pending", // default
		UserID:      userID,
	}
}

func ToTaskResponse(task Task) TaskResponse {
	return TaskResponse{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		DueDate:     task.DueDate,
		Status:      task.Status,
		UserID:      task.UserID,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
	}
}

func ApplyUpdate(task *Task, req UpdateTaskRequest) {
	_ = copier.CopyWithOption(task, &req, copier.Option{IgnoreEmpty: true, DeepCopy: true})
}
