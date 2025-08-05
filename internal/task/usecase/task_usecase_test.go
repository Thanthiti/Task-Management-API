package usecase_test

import (
	"errors"
	"mymodule/internal/task/model"
	"mymodule/internal/task/usecase"
	"mymodule/pkg/logger"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockTaskRepository struct {
	mock.Mock
}

func (m *MockTaskRepository) Save(task model.Task) error {
	args := m.Called(task)
	return args.Error(0)
}

func (m *MockTaskRepository) FindByID(taskID uint) (*model.Task, error) {
	args := m.Called(taskID)
	return args.Get(0).(*model.Task), args.Error(1)
}

func (m *MockTaskRepository) FindByUser(userID uint) (*[]model.Task, error) {
	args := m.Called(userID)
	return args.Get(0).(*[]model.Task), args.Error(1)
}

func (m *MockTaskRepository) FindByIDAndUser(taskID, userID uint) (*model.Task, error) {
	args := m.Called(taskID, userID)
	return args.Get(0).(*model.Task), args.Error(1)
}

func (m *MockTaskRepository) Update(task *model.Task) error {
	args := m.Called(task)
	return args.Error(0)
}

func (m *MockTaskRepository) Delete(taskID uint) error {
	args := m.Called(taskID)
	return args.Error(0)
}
func Testlog(m *testing.M) {
	logger.InitLogger()
	os.Exit(m.Run())
}

func TestCreateTask(t *testing.T) {
	logger.InitLogger()
	task := model.Task{
		Title:       "Test Task",
		Description: "Testing create task",
		UserID:      1,
	}

	t.Run("Success", func(t *testing.T) {
		mockRepo := new(MockTaskRepository)
		taskUC := usecase.NewTaskUsecase(mockRepo)

		mockRepo.On("Save", mock.AnythingOfType("model.Task")).Return(nil)

		err := taskUC.Create(task)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("SaveError", func(t *testing.T) {
		mockRepo := new(MockTaskRepository)
		taskUC := usecase.NewTaskUsecase(mockRepo)

		mockRepo.On("Save", mock.AnythingOfType("model.Task")).Return(errors.New("db error"))
		err := taskUC.Create(task)

		assert.EqualError(t, err, "db error")
		mockRepo.AssertExpectations(t)
	})

	t.Run("OverdueStatusSet", func(t *testing.T) {
		mockRepo := new(MockTaskRepository)
		taskUC := usecase.NewTaskUsecase(mockRepo)

		past := time.Now().Add(-24 * time.Hour)

		task := model.Task{
			Title:   "Overdue Task",
			UserID:  1,
			DueDate: &past,
		}

		mockRepo.On("Save", mock.MatchedBy(func(t model.Task) bool {
			return t.Status == "overdue"
		})).Return(nil)

		err := taskUC.Create(task)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestGetByUser(t *testing.T) {
	logger.InitLogger()

	t.Run("Success", func(t *testing.T) {
		mockRepo := new(MockTaskRepository)
		taskUC := usecase.NewTaskUsecase(mockRepo)

		tasks := []model.Task{
			{ID: 1, Title: "Task 1", UserID: 1},
			{ID: 2, Title: "Task 2", UserID: 1},
		}

		mockRepo.On("FindByUser", uint(1)).Return(&tasks, nil)

		result, err := taskUC.GetByUser(1)

		assert.NoError(t, err)
		assert.Equal(t, &tasks, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		mockRepo := new(MockTaskRepository)
		taskUC := usecase.NewTaskUsecase(mockRepo)

		mockRepo.On("FindByUser", uint(1)).Return((*[]model.Task)(nil), errors.New("db error"))

		result, err := taskUC.GetByUser(1)

		assert.Nil(t, result)
		assert.EqualError(t, err, "db error")
		mockRepo.AssertExpectations(t)
	})
}

func TestGetByIDAndUser(t *testing.T) {
	t.Run("Success", func(t *testing.T) {

		mockRepo := new(MockTaskRepository)
		taskUC := usecase.NewTaskUsecase(mockRepo)

		taskID := uint(1)
		userID := uint(100)
		expectedTask := &model.Task{
			ID:     taskID,
			Title:  "Test Task",
			UserID: userID,
		}
		mockRepo.On("FindByIDAndUser", taskID, userID).Return(expectedTask, nil)

		result, err := taskUC.GetByIDAndUser(taskID, userID)
		assert.NoError(t, err)
		assert.Equal(t, expectedTask, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		mockRepo := new(MockTaskRepository)
		taskUC := usecase.NewTaskUsecase(mockRepo)

		taskID := uint(1)
		userID := uint(100)

		mockRepo.On("FindByIDAndUser", taskID, userID).Return((*model.Task)(nil), errors.New("task not found"), errors.New("Task not found"))

		result, err := taskUC.GetByIDAndUser(taskID, userID)
		assert.Error(t, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}
func TestUpdateTask(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockRepo := new(MockTaskRepository)
		taskUC := usecase.NewTaskUsecase(mockRepo)

		taskID := uint(1)
		userID := uint(100)
		existingTask := &model.Task{
			ID:     taskID,
			UserID: userID,
			Title:  "Old Title",
		}
		title := "Updated title"
		description := "Updated description"
		status := "completed"
		dueDate := time.Now().Add(24 * time.Hour)

		input := &model.UpdateTaskInput{
			Title:       &title,
			Description: &description,
			Status:      &status,
			DueDate:     &dueDate,
		}
		
		mockRepo.On("FindByIDAndUser", taskID, userID).Return(existingTask, nil)
		mockRepo.On("Update", mock.AnythingOfType("*model.Task")).Return(nil)
		
		
		err := taskUC.UpdateTask(input,taskID, userID) 
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Task Not Found", func(t *testing.T) {
		mockRepo := new(MockTaskRepository)
		taskUC := usecase.NewTaskUsecase(mockRepo)

		taskID := uint(1)
		userID := uint(100)

		mockRepo.On("FindByIDAndUser", taskID, userID).Return((*model.Task)(nil), errors.New("Task not found"), errors.New("Task not found"))


		err := taskUC.UpdateTask( &model.UpdateTaskInput{},taskID, userID)
		assert.EqualError(t, err, "Task not found")
		mockRepo.AssertExpectations(t)
	})
	
	t.Run("Update Error", func(t *testing.T) {
		mockRepo := new(MockTaskRepository)
		taskUC := usecase.NewTaskUsecase(mockRepo)
		
		taskID := uint(1)
		userID := uint(100)
		existingTask := &model.Task{ID: taskID, UserID: userID}
		
		mockRepo.On("FindByIDAndUser", taskID, userID).Return(existingTask, nil)

		mockRepo.On("Update", mock.AnythingOfType("*model.Task")).Return(errors.New("Failed to update task"))
		
		err := taskUC.UpdateTask( &model.UpdateTaskInput{},taskID, userID)
		assert.EqualError(t, err, "Failed to update task")
		mockRepo.AssertExpectations(t)
	})
}

func TestDeleteTask(t *testing.T) {
	t.Run("Success", func(t *testing.T) {

		mockRepo := new(MockTaskRepository)
		taskUC := usecase.NewTaskUsecase(mockRepo)

		taskID := uint(1)
		userID := uint(100)

		expectedTask := &model.Task{
			ID:     taskID,
			Title:  "Test Task",
			UserID: userID,
		}
		mockRepo.On("FindByIDAndUser", taskID, userID).Return(expectedTask, nil)
		mockRepo.On("Delete", taskID).Return(nil)

		err := taskUC.DeleteTask(taskID, userID)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		mockRepo := new(MockTaskRepository)
		taskUC := usecase.NewTaskUsecase(mockRepo)

		taskID := uint(1)
		userID := uint(100)
		mockRepo.On("FindByIDAndUser", taskID, userID).Return((*model.Task)(nil), errors.New("Task not found"), errors.New("Task not found"))

		err := taskUC.DeleteTask(taskID, userID)
		assert.Error(t, err)
		assert.EqualError(t, err, "Task not found")
		mockRepo.AssertExpectations(t)
	})
}
