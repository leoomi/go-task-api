package main

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

type Task struct {
	ID          int    `json:"id" gorm:"primarykey"`
	Description string `json:"description"`
	Done        bool   `json:"done"`
}

func getAllTasks(w http.ResponseWriter, r *http.Request) {
	var tasks []Task
	db.Find(&tasks)

	resBody, _ := json.Marshal(tasks)

	w.Write(resBody)
	w.WriteHeader(http.StatusOK)
}

func createTask(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)

	var task Task
	json.Unmarshal(body, &task)

	db.Create(&task)

	resBody, _ := json.Marshal(task)

	w.Write(resBody)
	w.WriteHeader(http.StatusCreated)
}

func updateTaskById(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)

	id := chi.URLParam(r, "id")
	parsedId, _ := strconv.Atoi(id)

	var task Task
	db.First(&task, parsedId)

	var newTask Task
	json.Unmarshal(body, &newTask)

	task.Description = newTask.Description
	task.Done = newTask.Done

	db.Save(&task)

	resBody, _ := json.Marshal(task)

	w.Write(resBody)
	w.WriteHeader(http.StatusOK)
}

func deleteTaskById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	parsedId, _ := strconv.Atoi(id)

	result := db.Delete(&Task{}, parsedId)
	if result.RowsAffected == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func main() {
	var err error
	db, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&Task{})

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Route("/tasks", func(r chi.Router) {
		r.Get("/", getAllTasks)
		r.Post("/", createTask)
		r.Put("/{id}", updateTaskById)
		r.Delete("/{id}", deleteTaskById)
	})

	http.ListenAndServe(":8000", r)
}
