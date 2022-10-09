package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateTodoRequest struct {
	Title       string `json:"title" bson:"title" binding:"required"`
	Description string `json:"description" bson:"description" binding:"required"`
	User        string `json:"user" bson:"user" binding:"required"`
}

type Todo struct {
	Id          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Title       string             `json:"title,omitempty" bson:"title,omitempty"`
	Description string             `json:"description,omitempty" bson:"description,omitempty"`
	User        string             `json:"user,omitempty" bson:"user,omitempty"`
	Done        bool               `json:"done,omitempty" bson:"done,omitempty"`
}

type UpdateTodo struct {
	Title       string `json:"title,omitempty" bson:"title,omitempty"`
	Description string `json:"description,omitempty" bson:"description,omitempty"`
	User        string `json:"user,omitempty" bson:"user,omitempty"`
	Done        bool   `json:"done,omitempty" bson:"done,omitempty"`
}
