package grpc

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/todo-project/models"
	"github.com/todo-project/pb"
	"github.com/todo-project/services"
	"github.com/todo-project/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
)

type MockTodoServiceImpl struct {
	todoCollection *mongo.Collection
	ctx            context.Context
}

var (
	id1    = primitive.ObjectID{}
	id2, _ = primitive.ObjectIDFromHex("update_user")
	id3, _ = primitive.ObjectIDFromHex("update_title")
	id4, _ = primitive.ObjectIDFromHex("get_todo_id")
)

type mockGrpc_TodoServer struct {
	grpc.ServerStream
	Results []*pb.ToDo
}

func (_m *mockGrpc_TodoServer) Send(todo *pb.ToDo) error {
	_m.Results = append(_m.Results, todo)
	return nil
}

func (m MockTodoServiceImpl) CreateTodo(request *models.CreateTodoRequest) (*models.Todo, error) {
	if request.Title == "internal error" {
		return nil, errors.New("error creating todo")
	}
	return &models.Todo{
		Id:          id1,
		Title:       request.Title,
		Description: request.Description,
		User:        request.User,
		Done:        false,
	}, nil
}

func (m MockTodoServiceImpl) UpdateTodo(s string, todo *models.UpdateTodo) (*models.Todo, error) {
	if s == "internal error" {
		return nil, errors.New("error updating todo")
	}
	response := &models.Todo{
		Title:       todo.Title,
		Description: todo.Description,
		User:        todo.User,
		Done:        todo.Done,
	}
	if s == "update_user" {
		response.Id = id2
	}
	if s == "update_title" {
		response.Id = id3
	}
	return response, nil
}

func (m MockTodoServiceImpl) GetTodoById(s string) (*models.Todo, error) {
	if s == "internal error" {
		return nil, errors.New("error fetching todo")
	}
	return &models.Todo{
		Id:          id3,
		Title:       "title",
		Description: "",
		User:        "1",
		Done:        true,
	}, nil
}

func (m MockTodoServiceImpl) GetAllTodos(status pb.GetItemsRequest_TodoStatus, user string) ([]*models.Todo, error) {
	if user == "internal error" {
		return nil, errors.New("error updating todo")
	}
	if user == "1" {
		return []*models.Todo{
			{Id: primitive.ObjectID{}, Title: "one", User: "1"},
		}, nil
	}
	return []*models.Todo{
		{Id: primitive.ObjectID{}},
		{Id: primitive.ObjectID{}},
	}, nil
}

func (m MockTodoServiceImpl) DeleteTodo(s string) error {
	if s == "internal error" {
		return errors.New("error deleting todo")
	}
	if s == "nothing to delete" {
		return errors.New("todo does not exist")
	}
	return nil
}

func TestTodoServer_Create(t *testing.T) {
	type fields struct {
		UnimplementedToDoServiceServer pb.UnimplementedToDoServiceServer
		todoCollection                 *mongo.Collection
		todoService                    services.TodoService
	}
	type args struct {
		ctx context.Context
		req *pb.CreateItemRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *pb.TodoResponse
		wantErr bool
	}{
		{
			name: "create todo success",
			fields: fields{
				todoService: MockTodoServiceImpl{},
			},
			args: args{
				req: &pb.CreateItemRequest{
					Title:       "this one",
					Description: "desc 1",
					User:        "user 1",
				},
			},
			want: &pb.TodoResponse{
				ToDo: &pb.ToDo{
					Id:          id1.Hex(),
					Title:       "this one",
					Description: "desc 1",
					User:        "user 1",
					Done:        false,
				},
			},
			wantErr: false,
		},
		{
			name: "create todo failure",
			fields: fields{
				todoService: MockTodoServiceImpl{},
			},
			args: args{
				req: &pb.CreateItemRequest{Title: "internal error"},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := &TodoServer{
				UnimplementedToDoServiceServer: tt.fields.UnimplementedToDoServiceServer,
				todoCollection:                 tt.fields.todoCollection,
				todoService:                    tt.fields.todoService,
			}
			got, err := ts.Create(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Create() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTodoServer_Delete(t *testing.T) {
	type fields struct {
		UnimplementedToDoServiceServer pb.UnimplementedToDoServiceServer
		todoCollection                 *mongo.Collection
		todoService                    services.TodoService
	}
	type args struct {
		ctx context.Context
		req *pb.DeleteItemRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *pb.DeleteItemResponse
		wantErr bool
	}{
		{
			name: "delete todo success",
			fields: fields{
				todoService: MockTodoServiceImpl{},
			},
			args: args{
				req: &pb.DeleteItemRequest{
					Id: "valid_id",
				},
			},
			want:    &pb.DeleteItemResponse{Deleted: true},
			wantErr: false,
		},
		{
			name: "delete todo failure",
			fields: fields{
				todoService: MockTodoServiceImpl{},
			},
			args: args{
				req: &pb.DeleteItemRequest{Id: "nothing to delete"},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := &TodoServer{
				UnimplementedToDoServiceServer: tt.fields.UnimplementedToDoServiceServer,
				todoCollection:                 tt.fields.todoCollection,
				todoService:                    tt.fields.todoService,
			}
			got, err := ts.Delete(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Delete() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTodoServer_Get(t *testing.T) {
	type fields struct {
		UnimplementedToDoServiceServer pb.UnimplementedToDoServiceServer
		todoCollection                 *mongo.Collection
		todoService                    services.TodoService
	}
	type args struct {
		ctx context.Context
		req *pb.GetItemByID
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *pb.TodoResponse
		wantErr bool
	}{
		{
			name: "get todo success",
			fields: fields{
				todoService: MockTodoServiceImpl{},
			},
			args: args{
				req: &pb.GetItemByID{Id: "get_todo_id"},
			},
			want: &pb.TodoResponse{ToDo: &pb.ToDo{
				Id:          id4.Hex(),
				Title:       "title",
				Description: "",
				User:        "1",
				Done:        true,
			}},
			wantErr: false,
		},
		{
			name: "get todo failure",
			fields: fields{
				todoService: MockTodoServiceImpl{},
			},
			args: args{
				req: &pb.GetItemByID{Id: "internal error"},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := &TodoServer{
				UnimplementedToDoServiceServer: tt.fields.UnimplementedToDoServiceServer,
				todoCollection:                 tt.fields.todoCollection,
				todoService:                    tt.fields.todoService,
			}
			got, err := ts.Get(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTodoServer_GetAll(t *testing.T) {
	type fields struct {
		UnimplementedToDoServiceServer pb.UnimplementedToDoServiceServer
		todoCollection                 *mongo.Collection
		todoService                    services.TodoService
	}
	type args struct {
		req    *pb.GetItemsRequest
		stream pb.ToDoService_GetAllServer
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "get all todo success, stream response has one value",
			fields: fields{
				todoService: MockTodoServiceImpl{},
			},
			args: args{
				req: &pb.GetItemsRequest{
					User: utils.Pointer("1"),
				},
				stream: &mockGrpc_TodoServer{},
			},
			want:    1,
			wantErr: false,
		},
		{
			name: "get all todo success, stream response has multiple values",
			fields: fields{
				todoService: MockTodoServiceImpl{},
			},
			args: args{
				req: &pb.GetItemsRequest{
					User: utils.Pointer("2"),
				},
				stream: &mockGrpc_TodoServer{},
			},
			want:    2,
			wantErr: false,
		},
		{
			name: "get all todo failure",
			fields: fields{
				todoService: MockTodoServiceImpl{},
			},
			args: args{
				req: &pb.GetItemsRequest{
					User: utils.Pointer("internal error"),
				},
			},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := &TodoServer{
				UnimplementedToDoServiceServer: tt.fields.UnimplementedToDoServiceServer,
				todoCollection:                 tt.fields.todoCollection,
				todoService:                    tt.fields.todoService,
			}
			err := ts.GetAll(tt.args.req, tt.args.stream)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAll() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				if allTodos, ok := tt.args.stream.(*mockGrpc_TodoServer); ok {
					fmt.Println(len(allTodos.Results))
					if tt.want != len(allTodos.Results) {
						t.Errorf("GetAll() wanted = %d, got number of Todos %d", tt.want, len(allTodos.Results))
					}
				}
			}
		})
	}
}

func TestTodoServer_Update(t *testing.T) {
	type fields struct {
		UnimplementedToDoServiceServer pb.UnimplementedToDoServiceServer
		todoCollection                 *mongo.Collection
		todoService                    services.TodoService
	}
	type args struct {
		ctx context.Context
		req *pb.UpdateItemRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *pb.TodoResponse
		wantErr bool
	}{
		{
			name: "update USER in a Todo success",
			fields: fields{
				todoService: MockTodoServiceImpl{},
			},
			args: args{
				req: &pb.UpdateItemRequest{
					Id:   "update_user",
					User: utils.Pointer("new_user"),
				},
			},
			want: &pb.TodoResponse{ToDo: &pb.ToDo{
				Id:   id2.Hex(),
				User: "new_user",
			}},
			wantErr: false,
		},
		{
			name: "update TITLE in a Todo success",
			fields: fields{
				todoService: MockTodoServiceImpl{},
			},
			args: args{
				req: &pb.UpdateItemRequest{
					Id:    "update_title",
					Title: utils.Pointer("new_title"),
				},
			},
			want: &pb.TodoResponse{ToDo: &pb.ToDo{
				Id:    id3.Hex(),
				Title: "new_title",
			}},
			wantErr: false,
		},
		{
			name: "update todo failure",
			fields: fields{
				todoService: MockTodoServiceImpl{},
			},
			args: args{
				req: &pb.UpdateItemRequest{Id: "internal error"},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := &TodoServer{
				UnimplementedToDoServiceServer: tt.fields.UnimplementedToDoServiceServer,
				todoCollection:                 tt.fields.todoCollection,
				todoService:                    tt.fields.todoService,
			}
			got, err := ts.Update(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Update() got = %v, want %v", got, tt.want)
			}
		})
	}
}
