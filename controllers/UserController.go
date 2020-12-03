package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/AthanatiusC/TaskManager/models"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func UserCreate(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 14)

	err := models.GetDB("main").Collection("users").FindOne(context.TODO(), bson.M{"username": username})
	if err != nil {
		respondJSON(w, 409, "Username already exist!", map[string]interface{}{})
		return
	}
	user := models.User{
		ID:       primitive.NewObjectID(),
		Username: username,
		Password: string(hashedPassword),
	}
	models.GetDB("main").Collection("users").InsertOne(context.TODO(), &user)
	respondJSON(w, 200, "User successfully created!", user)

}

func UserGetOne(w http.ResponseWriter, r *http.Request) {
	var users models.User
	raw_param := mux.Vars(r)

	id := raw_param["id"]

	objid, _ := primitive.ObjectIDFromHex(id)

	err := models.GetDB("main").Collection("users").FindOne(context.TODO(), bson.M{"_id": objid}).Decode(&users)
	if err != nil {
		respondJSON(w, 404, "User not found!", map[string]interface{}{})
		return
	}
	respondJSON(w, 200, "Returned user detail", users)

}

func UserUpdate(w http.ResponseWriter, r *http.Request) {
	var user models.User
	userid, _ := primitive.ObjectIDFromHex(mux.Vars(r)["id"])
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(r.FormValue("password")), 14)

	user.ID = userid
	user.Username = r.FormValue("username")
	user.Password = string(hashedPassword)

	data := bson.D{{Key: "$set", Value: user}}
	result, err := models.GetDB("main").Collection("users").UpdateOne(context.TODO(), bson.M{"_id": userid}, data)
	if err != nil {
		respondJSON(w, 500, "Error occured", err)
		return
	}
	respondJSON(w, 200, "Successfully updated", result)
}

func UserDelete(w http.ResponseWriter, r *http.Request) {
	raw := mux.Vars(r)["id"]
	userid, _ := primitive.ObjectIDFromHex(raw)
	deleteResult, err := models.GetDB("main").Collection("users").DeleteOne(context.TODO(), bson.M{"_id": userid})
	if err != nil {
		respondJSON(w, 404, "Error!", err)
		return
	}
	respondJSON(w, 200, "User deleted", deleteResult)
}

func Auth(w http.ResponseWriter, r *http.Request) {
	var user models.User
	var dbuser models.User
	user.Username = r.FormValue("username")
	user.Password = r.FormValue("password")
	collection := models.GetDB("main").Collection("users")
	err := collection.FindOne(context.TODO(), bson.M{"username": user.Username}).Decode(&dbuser)
	if err != nil {
		respondJSON(w, 404, "Error!", err)
		return
	}
	// password:= hashAndSalt()
	ismatch := comparePasswords(dbuser.Password, []byte(user.Password))
	if ismatch {
		expire := time.Now().Add(24 * time.Hour * 7)
		cookie := http.Cookie{
			Name:    "Auth",
			Value:   dbuser.ID.Hex(),
			Expires: expire,
		}
		http.SetCookie(w, &cookie)
		respondJSON(w, 200, "Authenticated", map[string]interface{}{})
	}
}

func comparePasswords(hashedPwd string, plainPwd []byte) bool { // Since we'll be getting the hashed password from the DB it
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	if err != nil {
		return false
	}
	return true
}
