package model

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

func ToTaskFromUpdate(req UpdateTaskInput) *Task {
	return &Task{
		Title:       *req.Title,
		Description: *req.Description,
		DueDate:     req.DueDate,
	}
}


func ToTaskResponseList(tasks []Task) []TaskResponse {
	res := make([]TaskResponse, 0, len(tasks))
	for _, t := range tasks {
		res = append(res, ToTaskResponse(t))
	}
	return res
}

func ApplyUpdate(existing *Task, input UpdateTaskInput) {
    if input.Title != nil {
        existing.Title = *input.Title
    }
    if input.Description != nil {
        existing.Description = *input.Description
    }
    if input.DueDate != nil {
        existing.DueDate = input.DueDate
    }
    if input.Status != nil {
        existing.Status = *input.Status
    }
}
