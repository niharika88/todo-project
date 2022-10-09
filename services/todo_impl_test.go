package services

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/todo-project/models"
	"github.com/todo-project/pb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func TestTodoServiceImpl_CreateTodo(t1 *testing.T) {
	mt := mtest.New(t1, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()
	todoImpl := &TodoServiceImpl{
		ctx: context.TODO(),
	}
	createReq := &models.CreateTodoRequest{
		Title:       "title",
		Description: "desc",
		User:        "1",
	}
	id := "new_id"
	_id, _ := primitive.ObjectIDFromHex(id)
	expectedTodo := &models.Todo{
		Id:          _id,
		Title:       "title",
		Description: "desc",
		User:        "1",
		Done:        false,
	}

	mt.Run("success", func(mt *mtest.T) {
		todoImpl.todoCollection = mt.Coll
		mt.AddMockResponses(bson.D{
			{Key: "ok", Value: 1},
			{Key: "value", Value: bson.D{
				{Key: "_id", Value: expectedTodo.Id},
				{Key: "title", Value: expectedTodo.Title},
				{Key: "description", Value: expectedTodo.Description},
				{Key: "user", Value: expectedTodo.User},
				{Key: "done", Value: expectedTodo.Done},
			}},
		})

		newTodo, err := todoImpl.CreateTodo(createReq)
		assert.Nil(t1, err)
		assert.Equal(t1, expectedTodo.Title, newTodo.Title)
	})

	mt.Run("simple error", func(mt *mtest.T) {
		todoImpl.todoCollection = mt.Coll
		mt.AddMockResponses(bson.D{{Key: "ok", Value: 0}})

		newTodo, err := todoImpl.CreateTodo(&models.CreateTodoRequest{})

		assert.Nil(t1, newTodo)
		assert.NotNil(t1, err)
	})
}

func TestTodoServiceImpl_UpdateTodo(t1 *testing.T) {
	mt := mtest.New(t1, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()
	todoImpl := &TodoServiceImpl{
		ctx: context.TODO(),
	}
	updateReq := &models.UpdateTodo{
		Title:       "dummy",
		Description: "desc",
		User:        "1",
		Done:        true,
	}
	id := "dummy_id"
	_id, _ := primitive.ObjectIDFromHex(id)

	mt.Run("success", func(mt *mtest.T) {
		todoImpl.todoCollection = mt.Coll
		mt.AddMockResponses(bson.D{
			{Key: "ok", Value: 1},
			{Key: "value", Value: bson.D{
				{Key: "_id", Value: _id},
				{Key: "title", Value: updateReq.Title},
				{Key: "description", Value: updateReq.Description},
				{Key: "user", Value: updateReq.User},
				{Key: "done", Value: updateReq.Done},
			}},
		})

		updatedTodo, err := todoImpl.UpdateTodo(id, updateReq)

		assert.Nil(t1, err)
		assert.Equal(t1, updateReq.User, updatedTodo.User)
	})
}

func TestTodoServiceImpl_GetTodoById(t1 *testing.T) {
	mt := mtest.New(t1, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()
	todoImpl := &TodoServiceImpl{
		ctx: context.TODO(),
	}
	id := "dummy_id"
	_id, _ := primitive.ObjectIDFromHex(id)
	expectedTodo := &models.Todo{
		Id:          _id,
		Title:       "dummy",
		Description: "desc",
		User:        "1",
		Done:        true,
	}

	mt.Run("success", func(mt *mtest.T) {
		todoImpl.todoCollection = mt.Coll
		mt.AddMockResponses(mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch, bson.D{
			{"_id", expectedTodo.Id},
			{"title", expectedTodo.Title},
			{"description", expectedTodo.Description},
			{"user", expectedTodo.User},
			{"done", expectedTodo.Done},
		}))
		todoResponse, err := todoImpl.GetTodoById(id)
		assert.Nil(t1, err)
		assert.Equal(t1, expectedTodo, todoResponse)
	})
}

func TestTodoServiceImpl_GetAllTodos(t1 *testing.T) {
	mt := mtest.New(t1, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()
	todoImpl := &TodoServiceImpl{
		ctx: context.TODO(),
	}
	id1 := "dummy_id1"
	_id1, _ := primitive.ObjectIDFromHex(id1)
	expectedTodo1 := &models.Todo{
		Id:          _id1,
		Title:       "dummy",
		Description: "desc",
		User:        "1",
		Done:        false,
	}
	id2 := "dummy_id2"
	_id2, _ := primitive.ObjectIDFromHex(id2)
	expectedTodo2 := &models.Todo{
		Id:          _id2,
		Title:       "dummy",
		Description: "desc",
		User:        "2",
		Done:        false,
	}
	id3 := "dummy_id3"
	_id3, _ := primitive.ObjectIDFromHex(id3)
	expectedTodo3 := &models.Todo{
		Id:          _id3,
		Title:       "dummy",
		Description: "desc",
		User:        "2",
		Done:        true,
	}

	// Test to return ALL todos.
	mt.Run("success", func(mt *mtest.T) {
		todoImpl.todoCollection = mt.Coll
		first := mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch, bson.D{
			{"_id", expectedTodo1.Id},
			{"title", expectedTodo1.Title},
			{"description", expectedTodo1.Description},
			{"user", expectedTodo1.User},
			{"done", expectedTodo1.Done},
		})
		second := mtest.CreateCursorResponse(1, "foo.bar", mtest.NextBatch, bson.D{
			{"_id", expectedTodo2.Id},
			{"title", expectedTodo2.Title},
			{"description", expectedTodo2.Description},
			{"user", expectedTodo2.User},
			{"done", expectedTodo2.Done},
		})
		third := mtest.CreateCursorResponse(1, "foo.bar", mtest.NextBatch, bson.D{
			{"_id", expectedTodo3.Id},
			{"title", expectedTodo3.Title},
			{"description", expectedTodo3.Description},
			{"user", expectedTodo3.User},
			{"done", expectedTodo3.Done},
		})
		killCursors := mtest.CreateCursorResponse(0, "foo.bar", mtest.NextBatch)
		mt.AddMockResponses(first, second, third, killCursors)

		todos, err := todoImpl.GetAllTodos(pb.GetItemsRequest_ALL, "")
		assert.Nil(t1, err)
		assert.Equal(t1, []*models.Todo{
			expectedTodo1,
			expectedTodo2,
			expectedTodo3,
		}, todos)
	})

	// Test to return todo for a single user with a specific status.
	mt.Run("success", func(mt *mtest.T) {
		todoImpl.todoCollection = mt.Coll
		first := mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch, bson.D{
			{"_id", expectedTodo3.Id},
			{"title", expectedTodo3.Title},
			{"description", expectedTodo3.Description},
			{"user", expectedTodo3.User},
			{"done", expectedTodo3.Done},
		})
		killCursors := mtest.CreateCursorResponse(0, "foo.bar", mtest.NextBatch)
		mt.AddMockResponses(first, killCursors)

		todos, err := todoImpl.GetAllTodos(pb.GetItemsRequest_DONE, "2")
		assert.Nil(t1, err)
		assert.Equal(t1, []*models.Todo{
			expectedTodo3,
		}, todos)
	})
}

func TestTodoServiceImpl_DeleteTodo(t1 *testing.T) {
	mt := mtest.New(t1, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()
	todoImpl := &TodoServiceImpl{
		ctx: context.TODO(),
	}
	id := "dummy_id"

	mt.Run("success", func(mt *mtest.T) {
		todoImpl.todoCollection = mt.Coll
		mt.AddMockResponses(bson.D{{"ok", 1}, {"acknowledged", true}, {"n", 1}})
		err := todoImpl.DeleteTodo(id)
		assert.Nil(t1, err)
	})

	mt.Run("no document deleted", func(mt *mtest.T) {
		todoImpl.todoCollection = mt.Coll
		mt.AddMockResponses(bson.D{{"ok", 1}, {"acknowledged", true}, {"n", 0}})
		err := todoImpl.DeleteTodo(id)
		assert.NotNil(t1, err)
	})
}
