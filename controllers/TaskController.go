package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

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
}

func TaskCreate(w http.ResponseWriter, r *http.Request) {
	var task models.Task
	json.NewDecoder(r.Body).Decode(&task)
	task.ID = primitive.NewObjectID()
	task.Status = false
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
	var task models.Task
	json.NewDecoder(r.Body).Decode(&task)
	res, err := models.GetDB("main").Collection("tasks").UpdateOne(context.TODO(), bson.M{"_id": task.ID, "uid": task.UserID}, bson.D{{Key: "$set", Value: task}})
	if err != nil {
		respondJSON(w, 404, "Error occured", err)
		return
	}
	respondJSON(w, 200, "Task Updated!", res)
}
