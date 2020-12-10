package app

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strconv"

	// "strings"

	"github.com/AthanatiusC/TaskManager/models"
	jwt "github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Message(status bool, message string) map[string]interface{} {
	return map[string]interface{}{"status": status, "message": message}
}

func Respond(w http.ResponseWriter, data map[string]interface{}) {
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func Paginate(r *http.Request) (int, int, int) {
	currentPage, errC := strconv.Atoi(r.URL.Query().Get("page"))
	perPage, errP := strconv.Atoi(r.URL.Query().Get("per-page"))
	var offset int

	if errC != nil {
		offset = 0
		currentPage = 1
	}
	if errP != nil {
		perPage = 10
	}

	if currentPage > 1 {
		offset = currentPage * perPage
	}

	return currentPage, perPage, offset
}

func JwtAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		notAuth := []string{"/", "/api/v1/user/auth/", "/api/v1/user/"}
		requestPath := r.URL.Path
		response := make(map[string]interface{})

		// password := r.FormValue("password")
		// fmt.Println(password)

		for _, value := range notAuth {
			if value == requestPath {
				next.ServeHTTP(w, r)
				return
			}
		}

		tokenHeader := r.Header.Get("auth_key")
		UserID := r.Header.Get("user_id")
		tk := &models.Token{}

		var user models.User
		collection := models.GetDB("main").Collection("users")
		objid, _ := primitive.ObjectIDFromHex(UserID)
		err := collection.FindOne(context.TODO(), bson.M{"_id": objid}).Decode(&user)
		if err != nil {
			response = Message(false, "Auth Error")
			w.Header().Add("Content-Type", "application/json")
			Respond(w, response)
			return
		}

		token, err := jwt.ParseWithClaims(tokenHeader, tk, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("token_password")), nil
		})

		if err != nil { //Malformed token, returns with http code 403 as usual
			response = Message(false, "Malformed Auth!")
			w.Header().Add("Content-Type", "application/json")
			Respond(w, response)
			return
		}

		if !token.Valid { //Token is invalid, maybe not signed on this server
			response = Message(false, "Invalid Auth!")
			w.Header().Add("Content-Type", "application/json")
			Respond(w, response)
			return
		}

		if user.Token != tokenHeader {
			response = Message(false, "Incorrect Auth!")
			w.Header().Add("Content-Type", "application/json")
			Respond(w, response)
			return
		}

		//Everything went well, proceed with the request and set the caller to the user retrieved from the parsed token
		// ctx := context.WithValue(r.Context(), "user", tk.UserID)
		// r = r.WithContext(ctx)
		next.ServeHTTP(w, r) //proceed in the middleware chain!

	})
}
