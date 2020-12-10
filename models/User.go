package models

import (
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Payload struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type Token struct {
	UserID uint
	jwt.StandardClaims
}

//User structure
type User struct {
	ID         primitive.ObjectID `json:"id" bson:"_id" `
	Username   string             `json:"username" bson:"username"`
	Password   string             `json:"password" bson:"password" `
	Token      string             `json:"token" bson:"token" `
	Occupation string             `json:"occupation" bson:"occupation"`
	FirstName  string             `json:"firstname" bson:"firstname"`
	LastName   string             `json:"lastname" bson:"lastname"`
	Address    string             `json:"address" bson:"address"`
	Zip        string             `json:"zip" bson:"zip"`
	About      string             `json:"about" bson:"about"`
}
