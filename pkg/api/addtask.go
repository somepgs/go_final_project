package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/somepgs/go_final_project/pkg/db"
)

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJson(w, http.StatusBadRequest, map[string]any{"error": "Invalid JSON format"})
		return
	}
	// Validate the task
	if task.Title == "" {
		writeJson(w, http.StatusBadRequest, map[string]any{"error": "Title cannot be empty"})
		return
	}

	if err := checkDate(&task); err != nil {
		writeJson(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		return
	}

	// Add the task to the database
	id, err := db.AddTask(&task)
	if err != nil {
		writeJson(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		return
	}
	// Write the response
	writeJson(w, http.StatusCreated, map[string]any{"id": id})
}

func checkDate(task *db.Task) error {
	now := time.Now()
	if task.Date == "" {
		task.Date = now.Format(formatDate)
	}

	t, err := time.Parse(formatDate, task.Date)
	if err != nil {
		return err
	}

	next := now.Format(formatDate)
	if len(task.Repeat) != 0 {
		next, err = NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return err
		}
	}

	if afterNow(now, t) {
		if len(task.Repeat) == 0 {
			task.Date = now.Format(formatDate)
		}
		if len(task.Repeat) > 0 {
			task.Date = next
		}
	}
	return nil
}

// writeJson writes the given data as JSON to the response writer.
func writeJson(w http.ResponseWriter, status int, data any) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(status)
	// Set the Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	_, err = w.Write(jsonData)
	if err != nil {
		log.Printf("Error writing response: %v", err)
	}
}
