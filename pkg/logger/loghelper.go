package logger

import (
	userModel "mymodule/internal/user/model"
	taskModel "mymodule/internal/task/model"

	"github.com/sirupsen/logrus"
)

func LogUser(user userModel.User) *logrus.Entry {
	return Log.WithFields(map[string]interface{}{
		"userID": user.ID,
		"email":  user.Email,
	})
}


func LogTask(task taskModel.Task) *logrus.Entry {
	return Log.WithFields(map[string]interface{}{
		"taskID":   task.ID,
		"userID":   task.UserID,
		"title":    task.Title,
		"status":   task.Status,
		"due_date": task.DueDate,
	})
}

func LogFields(taskID, userID uint) map[string]interface{} {
	return map[string]interface{}{
		"taskID": taskID,
		"userID": userID,
	}
}

