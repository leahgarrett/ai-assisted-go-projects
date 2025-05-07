package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
)

type Task struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
	// ...other fields...
}

var dataFile = "data.json"

func GetTasks(w http.ResponseWriter, r *http.Request) {
	file, _ := os.Open(dataFile)
	defer file.Close()

	var tasks []Task
	_ = json.NewDecoder(file).Decode(&tasks)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func GetTaskByID(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	file, _ := os.Open(dataFile)
	defer file.Close()

	var tasks []Task
	_ = json.NewDecoder(file).Decode(&tasks)

	for _, task := range tasks {
		if task.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(task)
			return
		}
	}

	http.NotFound(w, r)
}

func CreateTask(w http.ResponseWriter, r *http.Request) {
	var newTask Task
	_ = json.NewDecoder(r.Body).Decode(&newTask)

	file, _ := os.Open(dataFile)
	defer file.Close()

	var tasks []Task
	_ = json.NewDecoder(file).Decode(&tasks)

	newTask.ID = len(tasks) + 1
	tasks = append(tasks, newTask)

	saveTasks(tasks)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newTask)
}

func UpdateTask(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	file, _ := os.Open(dataFile)
	defer file.Close()

	var tasks []Task
	_ = json.NewDecoder(file).Decode(&tasks)

	for i, task := range tasks {
		if task.ID == id {
			_ = json.NewDecoder(r.Body).Decode(&tasks[i])
			saveTasks(tasks)

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(tasks[i])
			return
		}
	}

	http.NotFound(w, r)
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	file, _ := os.Open(dataFile)
	defer file.Close()

	var tasks []Task
	_ = json.NewDecoder(file).Decode(&tasks)

	for i, task := range tasks {
		if task.ID == id {
			tasks = append(tasks[:i], tasks[i+1:]...)
			saveTasks(tasks)

			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	http.NotFound(w, r)
}

func saveTasks(tasks []Task) {
	file, _ := os.Create(dataFile)
	defer file.Close()

	_ = json.NewEncoder(file).Encode(tasks)
}
