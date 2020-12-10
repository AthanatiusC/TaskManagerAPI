package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/AthanatiusC/TaskManager/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func UserCreate(w http.ResponseWriter, r *http.Request) {
	// username := r.FormValue("username")
	// password := r.FormValue("password")
	var user models.User
	var dbuser models.User
	json.NewDecoder(r.Body).Decode(&user)
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	models.GetDB("main").Collection("users").FindOne(context.TODO(), bson.M{"username": user.Username}).Decode(&dbuser)
	if dbuser.Username == "" {
		user.ID = primitive.NewObjectID()
		user.Password = string(hashedPassword)
		claims := jwt.MapClaims{}
		claims["authorized"] = true
		claims["user_id"] = user.ID
		claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
		token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), claims)
		tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))

		user.Token = tokenString
		// fmt.Println(user.Username)
		models.GetDB("main").Collection("users").InsertOne(context.TODO(), &user)
		respondJSON(w, 200, "User successfully created!", user)
		return
	} else {
		fmt.Println(dbuser.Username)
		respondJSON(w, 409, "Username already exist!", user)
		return
	}
}

func UserGetOne(w http.ResponseWriter, r *http.Request) {
	var users models.User
	raw_param := mux.Vars(r)["id"]
	id := raw_param

	objid, _ := primitive.ObjectIDFromHex(id)
	fmt.Println(objid)

	err := models.GetDB("main").Collection("users").FindOne(context.TODO(), bson.M{"_id": objid}).Decode(&users)
	if err != nil {
		respondJSON(w, 404, "User not found!", map[string]interface{}{})
		return
	}
	respondJSON(w, 200, "Returned user detail", users)

}

func UserGetAll(w http.ResponseWriter, r *http.Request) {
	var user models.User
	var users []models.User
	// json.NewDecoder(r.Body).Decode(user)

	coll, err := models.GetDB("main").Collection("users").Find(context.TODO(), bson.M{})
	if err != nil {
		fmt.Println(err)
	}

	for coll.Next(context.TODO()) {
		coll.Decode(&user)
		users = append(users, user)

		user = models.User{}
	}
	respondJSON(w, 200, "Returned All user", users)
}

func UserUpdate(w http.ResponseWriter, r *http.Request) {
	var user models.User
	json.NewDecoder(r.Body).Decode(&user)
	// userid, _ := primitive.ObjectIDFromHex(mux.Vars(r)["id"])
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	user.Password = string(hashedPassword)

	data := bson.D{{Key: "$set", Value: user}}
	result, err := models.GetDB("main").Collection("users").UpdateOne(context.TODO(), bson.M{"_id": user.ID}, data)
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
	json.NewDecoder(r.Body).Decode(&user)
	collection := models.GetDB("main").Collection("users")
	err := collection.FindOne(context.TODO(), bson.M{"username": user.Username}).Decode(&dbuser)
	if err != nil {
		respondJSON(w, 404, "Error!", err)
		return
	}
	fmt.Println("User ", dbuser.ID, " requesting access!")
	// password:= hashAndSalt()
	ismatch := comparePasswords(dbuser.Password, []byte(user.Password))
	if ismatch {
		claims := jwt.MapClaims{}
		claims["authorized"] = true
		claims["user_id"] = user.ID
		claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
		token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), claims)
		tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))

		dbuser.Token = tokenString

		models.GetDB("main").Collection("users").UpdateOne(context.TODO(), bson.M{"_id": dbuser.ID}, bson.D{{Key: "$set", Value: dbuser}})

		respondJSON(w, http.StatusOK, "Successfully Logged In!", map[string]interface{}{"token": tokenString, "id": dbuser.ID})
		return
	}
	respondJSON(w, 200, "Wrong password or email!", map[string]interface{}{})
}

func comparePasswords(hashedPwd string, plainPwd []byte) bool { // Since we'll be getting the hashed password from the DB it
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	if err != nil {
		return false
	}
	return true
}
