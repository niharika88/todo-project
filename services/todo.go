package services

import (
	"github.com/todo-project/models"
	"github.com/todo-project/pb"
)

type TodoService interface {
	CreateTodo(request *models.CreateTodoRequest) (*models.Todo, error)
	UpdateTodo(string, *models.UpdateTodo) (*models.Todo, error)
	GetTodoById(string) (*models.Todo, error)
	GetAllTodos(status pb.GetItemsRequest_TodoStatus, user string) ([]*models.Todo, error)
	DeleteTodo(string) error
}
