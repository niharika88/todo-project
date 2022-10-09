package grpc

import (
	"context"

	"github.com/todo-project/models"
	"github.com/todo-project/pb"
	"github.com/todo-project/services"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TodoServer struct {
	pb.UnimplementedToDoServiceServer
	todoCollection *mongo.Collection
	todoService    services.TodoService
}

func NewGrpcTodoServer(todoCollection *mongo.Collection, todoService services.TodoService) (*TodoServer, error) {
	todoServer := &TodoServer{
		todoCollection: todoCollection,
		todoService:    todoService,
	}

	return todoServer, nil
}

func (ts *TodoServer) Create(_ context.Context, req *pb.CreateItemRequest) (*pb.TodoResponse, error) {
	post := &models.CreateTodoRequest{
		Title:       req.GetTitle(),
		Description: req.GetDescription(),
		User:        req.GetUser(),
	}

	newTodo, err := ts.todoService.CreateTodo(post)

	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	res := &pb.TodoResponse{
		ToDo: &pb.ToDo{
			Id:          newTodo.Id.Hex(),
			Title:       newTodo.Title,
			Description: newTodo.Description,
			User:        newTodo.User,
		},
	}
	return res, nil
}

func (ts *TodoServer) Update(_ context.Context, req *pb.UpdateItemRequest) (*pb.TodoResponse, error) {
	todo := &models.UpdateTodo{
		Title:       req.GetTitle(),
		Description: req.GetDescription(),
		Done:        req.GetDone(),
		User:        req.GetUser(),
	}

	updatedTodo, err := ts.todoService.UpdateTodo(req.GetId(), todo)

	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	res := &pb.TodoResponse{
		ToDo: &pb.ToDo{
			Id:          updatedTodo.Id.Hex(),
			Title:       updatedTodo.Title,
			Description: updatedTodo.Description,
			Done:        updatedTodo.Done,
			User:        updatedTodo.User,
		},
	}
	return res, nil
}

func (ts *TodoServer) Get(_ context.Context, req *pb.GetItemByID) (*pb.TodoResponse, error) {
	todo, err := ts.todoService.GetTodoById(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	res := &pb.TodoResponse{
		ToDo: &pb.ToDo{
			Id:          todo.Id.Hex(),
			Title:       todo.Title,
			Description: todo.Description,
			Done:        todo.Done,
			User:        todo.User,
		},
	}
	return res, nil
}

func (ts *TodoServer) GetAll(req *pb.GetItemsRequest, stream pb.ToDoService_GetAllServer) error {
	todos, err := ts.todoService.GetAllTodos(req.GetStatus(), req.GetUser())
	if err != nil {
		return status.Errorf(codes.Internal, err.Error())
	}
	for _, todo := range todos {
		err = stream.Send(&pb.ToDo{
			Id:          todo.Id.Hex(),
			Title:       todo.Title,
			Description: todo.Description,
			User:        todo.User,
			Done:        todo.Done,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (ts *TodoServer) Delete(_ context.Context, req *pb.DeleteItemRequest) (*pb.DeleteItemResponse, error) {
	if err := ts.todoService.DeleteTodo(req.GetId()); err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	res := &pb.DeleteItemResponse{
		Deleted: true,
	}
	return res, nil
}
