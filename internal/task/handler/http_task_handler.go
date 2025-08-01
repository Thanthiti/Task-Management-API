package handler

import (
	"mymodule/internal/task/usecase"
	"mymodule/internal/task/model"
	"mymodule/pkg/logger"

	"github.com/gofiber/fiber/v2"
)

type HttpTaskhandler struct {
	usecase usecase.TaskUsecase
}

func NewUserHandler(app *fiber.App, usecase usecase.TaskUsecase) {
	handler := &HttpTaskhandler{
		usecase: usecase,
	}
	logger.Log.Debug(handler)
	
}


	func (h *HttpTaskhandler)Create(task model.Task) error{
		return nil
	}
	func (h *HttpTaskhandler)GetByID(taskID uint) (string, error){
		return  "",nil
	}
	func (h *HttpTaskhandler)GetByUser(userName string) (string, error){
		return "",nil
	}
	func (h *HttpTaskhandler)Updatetask(task model.Task) error{
		return nil
	}
	func (h *HttpTaskhandler)Deletetask(taskID uint) error{
		return nil
	}