/* Copyright (C) Ahmad Saugi & Lexi Anugrah - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Lexi Anugrah <athanatius@4save.me>, November 2019
 */

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	// "github.com/jinzhu/gorm"
	"encoding/json"

	"github.com/AthanatiusC/TaskManager/controllers"
)

func main() {
	router := mux.NewRouter()
	router.Methods("OPTIONS").HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Headers", "*")
			w.WriteHeader(200)
			json.NewEncoder(w).Encode("PREFLIGHT OK")
		})
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "socket-test.html")
		fmt.Println("home")
	})
	// http.Handle("/", http.FileServer(http.Dir("./assets")))
	APP_PORT := os.Getenv("APP_PORT")
	if APP_PORT == "" {
		APP_PORT = "8088"
	}

	apiV1 := router.PathPrefix("/api/v1").Subrouter()

	v1Task := apiV1.PathPrefix("/task").Subrouter()
	v1Task.HandleFunc("/personal/{id}", controllers.TaskGetAll).Methods("GET", "OPTIONS") // View All
	v1Task.HandleFunc("/", controllers.TaskCreate).Methods("POST", "OPTIONS")             // Store
	v1Task.HandleFunc("/{id}", controllers.TaskUpdate).Methods("PUT", "OPTIONS")          // Update
	v1Task.HandleFunc("/{id}", controllers.TaskGetOne).Methods("GET", "OPTIONS")          // Get Detail
	v1Task.HandleFunc("/delete/{id}", controllers.TaskDelete).Methods("GET", "OPTIONS")   // Delete

	v1User := apiV1.PathPrefix("/user").Subrouter()
	// v1User.HandleFunc("/", controllers.UserGetOne).Methods("GET", "OPTIONS")        // View All
	v1User.HandleFunc("/", controllers.UserCreate).Methods("POST", "OPTIONS")           // Store
	v1User.HandleFunc("/auth", controllers.Auth).Methods("POST", "OPTIONS")             // Store
	v1User.HandleFunc("/{id}", controllers.UserGetOne).Methods("GET", "OPTIONS")        // Get Detail
	v1User.HandleFunc("/{id}", controllers.UserUpdate).Methods("PUT", "OPTIONS")        // Update
	v1User.HandleFunc("/delete/{id}", controllers.UserDelete).Methods("GET", "OPTIONS") // Delete

	//  apiV1.Use(app.JwtAuthentication)
	//  apiV2 := router.PathPrefix("/api/v2").Subrouter()

	// v1Auth := apiV1.PathPrefix("/auth").Subrouter()
	// v1Auth.HandleFunc("/login", controllers.AuthLogin).Methods("POST", "OPTIONS")
	// v1Auth.HandleFunc("/register", controllers.AuthRegister).Methods("POST", "OPTIONS")

	//  v1Camera := apiV1.PathPrefix("/camera").Subrouter()
	//  v1Camera.HandleFunc("/", controllers.CameraIndex).Methods("GET","OPTIONS")         // View All
	//  v1Camera.HandleFunc("/", controllers.CameraStore).Methods("POST","OPTIONS")        // Store
	//  v1Camera.HandleFunc("/{id}", controllers.CameraDetail).Methods("GET","OPTIONS")    // Get Detail
	//  v1Camera.HandleFunc("/{id}", controllers.CameraUpdate).Methods("PUT","OPTIONS")    // Update
	//  v1Camera.HandleFunc("/{id}", controllers.CameraDelete).Methods("DELETE","OPTIONS") // Delete

	//  v1ClassRoom := apiV1.PathPrefix("/classroom").Subrouter()
	//  v1ClassRoom.HandleFunc("/", controllers.ClassroomIndex).Methods("GET","OPTIONS")         // View All
	//  v1ClassRoom.HandleFunc("/", controllers.ClassroomStore).Methods("POST","OPTIONS")        // Store
	//  v1ClassRoom.HandleFunc("/{id}", controllers.ClassroomDetail).Methods("GET","OPTIONS")    // Get Detail
	//  v1ClassRoom.HandleFunc("/{id}", controllers.ClassroomUpdate).Methods("PUT","OPTIONS")    // Update
	//  v1ClassRoom.HandleFunc("/{id}", controllers.ClassroomDelete).Methods("DELETE","OPTIONS") // Delete

	//  v1Schedule := apiV1.PathPrefix("/schedule").Subrouter()
	//  v1Schedule.HandleFunc("/", controllers.ScheduleIndex).Methods("GET","OPTIONS")         // View All
	//  v1Schedule.HandleFunc("/", controllers.ScheduleStore).Methods("POST","OPTIONS")        // Store
	//  v1Schedule.HandleFunc("/{id}", controllers.ScheduleDetail).Methods("GET","OPTIONS")    // Get Detail
	//  v1Schedule.HandleFunc("/{id}", controllers.ScheduleUpdate).Methods("PUT","OPTIONS")    // Update
	//  v1Schedule.HandleFunc("/{id}", controllers.ScheduleDelete).Methods("DELETE","OPTIONS") // Delete

	//  v1Logs := apiV1.PathPrefix("/logs").Subrouter()
	//  v1Logs.HandleFunc("/", controllers.LogIndex).Methods("GET","OPTIONS")      // View All
	//  v1Logs.HandleFunc("/", controllers.LogStore).Methods("POST","OPTIONS")     // Store
	//  v1Logs.HandleFunc("/{id}", controllers.LogDetail).Methods("GET","OPTIONS") // Get Detail

	//  v1Subject := apiV1.PathPrefix("/subject").Subrouter()
	//  v1Subject.HandleFunc("/", controllers.SubjectIndex).Methods("GET","OPTIONS")         // View All
	//  v1Subject.HandleFunc("/", controllers.SubjectStore).Methods("POST","OPTIONS")        // Store
	//  v1Subject.HandleFunc("/{id}", controllers.SubjectDetail).Methods("GET","OPTIONS")    // Get Detail
	//  v1Subject.HandleFunc("/{id}", controllers.SubjectUpdate).Methods("PUT","OPTIONS")    // Update
	//  v1Subject.HandleFunc("/{id}", controllers.SubjectDelete).Methods("DELETE","OPTIONS") // Delete

	//  v1Room := apiV1.PathPrefix("/room").Subrouter()
	//  v1Room.HandleFunc("/", controllers.RoomIndex).Methods("GET","OPTIONS")         // View All
	//  v1Room.HandleFunc("/", controllers.RoomStore).Methods("POST","OPTIONS")        // Store
	//  v1Room.HandleFunc("/{id}", controllers.RoomDetail).Methods("GET","OPTIONS")    // Get Detail
	//  v1Room.HandleFunc("/{id}", controllers.RoomUpdate).Methods("PUT","OPTIONS")    // Update
	//  v1Room.HandleFunc("/{id}", controllers.RoomDelete).Methods("DELETE","OPTIONS") // Delete

	//  v1RoomAccess := apiV1.PathPrefix("/room_access").Subrouter()
	//  v1RoomAccess.HandleFunc("/", controllers.RoomAccessIndex).Methods("GET","OPTIONS")         // View All
	//  v1RoomAccess.HandleFunc("/", controllers.RoomAccessStore).Methods("POST","OPTIONS")        // Store
	//  v1RoomAccess.HandleFunc("/{id}", controllers.RoomAccessDetail).Methods("GET","OPTIONS")    // Get Detail
	//  v1RoomAccess.HandleFunc("/{id}", controllers.RoomAccessUpdate).Methods("PUT","OPTIONS")    // Update
	//  v1RoomAccess.HandleFunc("/{id}", controllers.RoomAccessDelete).Methods("DELETE","OPTIONS") // Delete

	//  v2Attendance := apiV2.PathPrefix("/attendance").Subrouter()
	//  v2Attendance.HandleFunc("/", controllers.AttendanceV2).Methods("GET","OPTIONS")        // View
	//  v2Attendance.HandleFunc("/new", controllers.AttendanceV2New).Methods("POST","OPTIONS") // Store

	//  v2User := apiV2.PathPrefix("/user").Subrouter()
	//  v2User.HandleFunc("/", controllers.UserV2Index).Methods("GET","OPTIONS")                                           // View All
	//  v2User.HandleFunc("/embeddings", controllers.UserV2Embeddings).Methods("GET","OPTIONS")                            // Verify
	//  v2User.HandleFunc("/embeddings/clear", controllers.UserV2EmbeddingsClear).Methods("GET","OPTIONS")                 // Verify
	//  v2User.HandleFunc("/embeddings/clear/{user_id}", controllers.UserV2EmbeddingsClearOnUser).Methods("GET","OPTIONS") // Verify
	//  v2User.HandleFunc("/{id}", controllers.UserV2Detail).Methods("GET","OPTIONS")                                      // Detail
	//  v2User.HandleFunc("/verify", controllers.UserV2Verify).Methods("POST","OPTIONS")                                   // Verify
	//  v2User.HandleFunc("/register", controllers.UserRegister).Methods("POST","OPTIONS")                                 // Store
	//  v2User.HandleFunc("/recognize", controllers.UserRecognize).Methods("OPTIONS","POST")                     // Recognize

	//  apiV2.HandleFunc("/room_accesss/check", controllers.RoomAccessCheck).Methods("POST","OPTIONS")
	//  apiV2.HandleFunc("/classroom", controllers.ClassroomV2Index).Methods("GET","OPTIONS")
	//  apiV2.HandleFunc("/logs", controllers.LogIndex).Methods("GET","OPTIONS")
	//  apiV2.HandleFunc("/camera", controllers.CameraIndex).Methods("GET","OPTIONS")
	//  apiV2.HandleFunc("/test", controllers.ImportCsv).Methods("GET","OPTIONS")

	fmt.Println("App running on port " + APP_PORT)
	log.Fatal(http.ListenAndServe(":"+APP_PORT, router))
}
