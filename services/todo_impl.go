package services

import (
	"context"
	"errors"

	"github.com/todo-project/models"
	"github.com/todo-project/pb"
	"github.com/todo-project/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TodoServiceImpl struct {
	todoCollection *mongo.Collection
	ctx            context.Context
}

func NewTodoService(todoCollection *mongo.Collection, ctx context.Context) TodoService {
	return &TodoServiceImpl{todoCollection, ctx}
}

func (t *TodoServiceImpl) CreateTodo(todo *models.CreateTodoRequest) (*models.Todo, error) {
	res, err := t.todoCollection.InsertOne(t.ctx, todo)
	if err != nil {
		return nil, err
	}

	createdTodo := &models.Todo{
		Id:          res.InsertedID.(primitive.ObjectID),
		Title:       todo.Title,
		Description: todo.Description,
		User:        todo.User,
	}

	return createdTodo, nil
}

func (t *TodoServiceImpl) UpdateTodo(id string, data *models.UpdateTodo) (*models.Todo, error) {
	doc, err := utils.ToMongoBson(data)
	if err != nil {
		return nil, err
	}

	obId, _ := primitive.ObjectIDFromHex(id)
	query := bson.D{{Key: "_id", Value: obId}}
	update := bson.D{{Key: "$set", Value: doc}}
	res := t.todoCollection.FindOneAndUpdate(t.ctx, query, update, options.FindOneAndUpdate().SetReturnDocument(1))

	var updatedPost *models.Todo
	if err := res.Decode(&updatedPost); err != nil {
		return nil, errors.New("no Todo document found for given Id")
	}

	return updatedPost, nil
}

func (t *TodoServiceImpl) GetTodoById(id string) (*models.Todo, error) {
	objectId, _ := primitive.ObjectIDFromHex(id)

	query := bson.M{"_id": objectId}

	var todo *models.Todo
	if err := t.todoCollection.FindOne(t.ctx, query).Decode(&todo); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("no Todo document found for given Id")
		}
		return nil, err
	}
	return todo, nil
}

func (t *TodoServiceImpl) GetAllTodos(status pb.GetItemsRequest_TodoStatus, user string) ([]*models.Todo, error) {

	query := bson.M{}
	switch status {
	case pb.GetItemsRequest_DONE:
		query["done"] = true
	case pb.GetItemsRequest_PENDING:
		query["done"] = false
	default:
		// do nothing, since we need all the Todos
	}
	if len(user) != 0 {
		query["user"] = user
	}

	cursor, err := t.todoCollection.Find(t.ctx, query)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(t.ctx)

	var todoList []*models.Todo
	for cursor.Next(t.ctx) {
		todo := &models.Todo{}
		if err = cursor.Decode(todo); err != nil {
			return nil, err
		}
		todoList = append(todoList, todo)
	}

	if err = cursor.Err(); err != nil {
		return nil, err
	}

	if len(todoList) == 0 {
		return []*models.Todo{}, nil
	}
	return todoList, nil
}

func (t *TodoServiceImpl) DeleteTodo(id string) error {
	objectId, _ := primitive.ObjectIDFromHex(id)
	query := bson.M{"_id": objectId}

	res, err := t.todoCollection.DeleteOne(t.ctx, query)
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return errors.New("no Todo document found for given Id")
	}
	return nil
}
