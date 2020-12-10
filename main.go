/* Copyright (C) Ahmad Saugi & Lexi Anugrah - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Lexi Anugrah <athanatius@4save.me>, November 2019
 */

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	// "github.com/jinzhu/gorm"

	"github.com/AthanatiusC/TaskManager/app"
	"github.com/AthanatiusC/TaskManager/controllers"
)

func main() {
	router := mux.NewRouter()
	router.Methods("OPTIONS").HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Headers", "*")
			w.Header().Set("Access-Control-Allow-Methods", "*")
			w.WriteHeader(200)
			json.NewEncoder(w).Encode("PREFLIGHT OK")
		})
	APP_PORT := os.Getenv("APP_PORT")
	if APP_PORT == "" {
		APP_PORT = "8088"
	}

	apiV1 := router.PathPrefix("/api/v1").Subrouter()
	apiV1.Use(app.JwtAuthentication)

	v1Task := apiV1.PathPrefix("/task").Subrouter()
	v1Task.HandleFunc("/personal/{id}/", controllers.TaskGetAll).Methods("GET", "OPTIONS") // View All
	v1Task.HandleFunc("/", controllers.TaskCreate).Methods("POST", "OPTIONS")              // Store
	v1Task.HandleFunc("/update/", controllers.TaskUpdate).Methods("POST", "OPTIONS")       // Update
	v1Task.HandleFunc("/{id}/", controllers.TaskGetOne).Methods("GET", "OPTIONS")          // Get Detail
	v1Task.HandleFunc("/delete/{id}/", controllers.TaskDelete).Methods("GET", "OPTIONS")   // Delete

	v1User := apiV1.PathPrefix("/user").Subrouter()
	// v1User.HandleFunc("/", controllers.UserGetOne).Methods("GET", "OPTIONS")        // View All
	v1User.HandleFunc("/", controllers.UserCreate).Methods("POST", "OPTIONS")           // Store
	v1User.HandleFunc("/", controllers.UserGetAll).Methods("GET", "OPTIONS")            // Store
	v1User.HandleFunc("/auth/", controllers.Auth).Methods("POST", "OPTIONS")            // Store
	v1User.HandleFunc("/{id}", controllers.UserGetOne).Methods("GET", "OPTIONS")        // Get Detail
	v1User.HandleFunc("/update", controllers.UserUpdate).Methods("PUT", "OPTIONS")      // Update
	v1User.HandleFunc("/delete/{id}", controllers.UserDelete).Methods("GET", "OPTIONS") // Delete
	// router.Use(mux.CORSMethodMiddleware(router))

	fmt.Println("App running on port " + APP_PORT)
	log.Fatal(http.ListenAndServe(":"+APP_PORT, router))
}
