package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

//Task structure
type Task struct {
	ID          primitive.ObjectID `json:"id" bson:"_id" `
	UserID      primitive.ObjectID `json:"userid" bson:"uid"`
	Name        string             `json:"name" bson:"name"`
	Time        time.Time          `json:"time" bson:"time"`
	Place       string             `json:"place" bson:"place"`
	Description string             `json:"description" bson:"description"`
	Status      bool               `json:"status" bson:"status"`
}
