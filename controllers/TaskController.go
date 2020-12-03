package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/AthanatiusC/TaskManager/models"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//TaskGetAll for return index api
func TaskGetAll(w http.ResponseWriter, r *http.Request) {
	raw := mux.Vars(r)["id"]
	userid, _ := primitive.ObjectIDFromHex(raw)
	tasks := []models.Task{}
	task := models.Task{}

	ctx := context.TODO() // Options to the database.
	coll, err := models.GetDB("main").Collection("tasks").Find(ctx, bson.M{"uid": userid})
	if err != nil {
		fmt.Println(err)
	}

	for coll.Next(ctx) {
		coll.Decode(&task)
		task.Time = task.Time
		tasks = append(tasks, task)
		task = models.Task{}
	}
	respondJSON(w, 200, "Success get all task for current users!", tasks)
}

// TaskGetOne for returning single item
func TaskGetOne(w http.ResponseWriter, r *http.Request) {
	raw := mux.Vars(r)["id"]
	taskID, _ := primitive.ObjectIDFromHex(raw)

	var task models.Task
	err := models.GetDB("main").Collection("tasks").FindOne(context.TODO(), bson.M{"_id": taskID}).Decode(&task)

	if err != nil {
		fmt.Println(err)
		respondJSON(w, 200, "Task not found!", map[string]interface{}{})
		return
	}

	respondJSON(w, 200, "Get Task Detail", task)
	return
}

func TaskCreate(w http.ResponseWriter, r *http.Request) {
	var task models.Task
	json.NewDecoder(r.Body).Decode(&task)
	task.ID = primitive.NewObjectID()
	task.Status = false
	// userid, _ := primitive.ObjectIDFromHex(task.UserID)
	// time, _ := time.Parse("2006-01-02", task.Time)
	// newTask := models.Task{
	// 	ID:          primitive.NewObjectID(),
	// 	UserID:      userid,
	// 	Name:        r.FormValue("name"),
	// 	Time:        time,
	// 	Place:       r.FormValue("place"),
	// 	Description: r.FormValue("description"),
	// 	Status:      false,
	// }
	models.GetDB("main").Collection("tasks").InsertOne(context.TODO(), &task)
	respondJSON(w, 200, "Success Create New Task!", task)
}

func TaskDelete(w http.ResponseWriter, r *http.Request) {
	raw := mux.Vars(r)["id"]
	taskid, _ := primitive.ObjectIDFromHex(raw)
	deleteResult, err := models.GetDB("main").Collection("tasks").DeleteOne(context.TODO(), bson.M{"_id": taskid})
	if err != nil {
		respondJSON(w, 404, "Error!", err)
		return
	}
	respondJSON(w, 200, "Task deleted", deleteResult)
}

func TaskUpdate(w http.ResponseWriter, r *http.Request) {
	taskid, _ := primitive.ObjectIDFromHex(mux.Vars(r)["id"])
	uid, _ := primitive.ObjectIDFromHex(r.FormValue("uid"))
	time, _ := time.Parse("2006-01-02", r.FormValue("time"))
	val, _ := strconv.ParseBool(r.FormValue("status"))

	var task models.Task
	task.ID = taskid
	task.UserID = uid
	task.Name = r.FormValue("name")
	task.Time = time
	task.Place = r.FormValue("place")
	task.Description = r.FormValue("description")
	task.Status = val

	data := bson.D{{Key: "$set", Value: task}}
	res, err := models.GetDB("main").Collection("users").UpdateOne(context.TODO(), bson.M{"_id": taskid, "uid": uid}, data)
	if err != nil {
		respondJSON(w, 500, "Error occured", err)
		return
	}
	respondJSON(w, 200, "Successfully updated", res)
}
