package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Task ...
type Task struct {
    ID           string   `json:"id"`
    Description  string   `json:"description"`
    Note         string   `json:"note"`
    Applications []string `json:"applications"`
}


var tasks = map[string]Task{
    "1": {
        ID:          "1",
        Description: "Сделать финальное задание темы REST API",
        Note:        "Если сегодня сделаю, то завтра будет свободный день. Ура!",
        Applications: []string{
            "VS Code",
            "Terminal",
            "git",
        },
    },
    "2": {
        ID:          "2",
        Description: "Протестировать финальное задание с помощью Postmen",
        Note:        "Лучше это делать в процессе разработки, каждый раз, когда запускаешь сервер и проверяешь хендлер",
        Applications: []string{
            "VS Code",
            "Terminal",
            "git",
            "Postman",
        },
    },
}

// handlerFn for receiving all saved tasks
func getTasks(w http.ResponseWriter, r *http.Request) {
    resp, err := json.Marshal(tasks) 
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write(resp)
}

// handlerFn for adding new task
func postTask(w http.ResponseWriter, r *http.Request) {
    var task Task
    var buf bytes.Buffer

    _, err := buf.ReadFrom(r.Body)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    tasks[task.ID] = task

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
}

// handlerFn for receiving specified task
func getTask(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")

    task, ok := tasks[id]
    if !ok {
        http.Error(w, "Task not found", http.StatusNoContent)
        return
    }

    resp, err := json.Marshal(task)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write(resp)
}

// handlerFn for removing specified task
func deleteTask(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")

    _, ok := tasks[id]
    if !ok {
        http.Error(w, "task not found", http.StatusBadRequest)
        return
    }
    delete(tasks, id)
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)

}

func main() {
    r := chi.NewRouter()

    // endpoint `/tasks` - method GET
    r.Get("/tasks", getTasks)

    // endpoint `/tasks` - method POST
    r.Post("/tasks", postTask)

    // endpoint `/tasks/{id}` - method GET
    r.Get("/tasks/{id}", getTask)

    // endpoint `/tasks/{id}` - method DELETE
    r.Delete("/tasks/{id}", deleteTask)

    if err := http.ListenAndServe(":8080", r); err != nil {
        fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
        return
    }
}

