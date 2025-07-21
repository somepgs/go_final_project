package api

import (
	"encoding/json"
	"github.com/somepgs/go_final_project/pkg/db"
	"net/http"
	"time"
)

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		writeJson(w, http.StatusNotFound, map[string]any{"error": "Не указан ID задачи"})
		return
	}
	task, err := db.GetTask(id)
	if err != nil {
		writeJson(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
		return
	}
	if task == nil {
		writeJson(w, http.StatusNotFound, map[string]any{"error": "Задача не найдена"})
		return
	}
	writeJson(w, http.StatusOK, task)
}

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJson(w, http.StatusBadRequest, map[string]any{"error": "Неверный формат данных"})
		return
	}
	// Validate the task
	if task.Title == "" {
		writeJson(w, http.StatusBadRequest, map[string]any{"error": "Заголовок задачи не может быть пустым"})
		return
	}
	// Check if the date is valid
	if err := checkDate(&task); err != nil {
		writeJson(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		return
	}
	// Update the task in the database
	err := db.UpdateTask(&task)
	if err != nil {
		writeJson(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
		return
	}
	writeJson(w, http.StatusOK, map[string]any{})
}

func doneTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJson(w, http.StatusMethodNotAllowed, map[string]any{"error": "Метод не поддерживается"})
		return
	}
	id := r.FormValue("id")
	if id == "" {
		writeJson(w, http.StatusBadRequest, map[string]any{"error": "Не указан ID задачи"})
		return
	}
	// Retrieve the task from the database
	task, err := db.GetTask(id)
	if err != nil {
		writeJson(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
		return
	}
	if task == nil {
		writeJson(w, http.StatusNotFound, map[string]any{"error": "Задача не найдена"})
		return
	}
	// If the task has no repeat, delete it; otherwise, update the date
	if len(task.Repeat) == 0 {
		err = db.DeleteTask(id)
		if err != nil {
			writeJson(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
			return
		}
	}
	if len(task.Repeat) > 0 {
		// Update the task's date to the next occurrence based on the repeat pattern
		now := time.Now()
		next, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			writeJson(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
			return
		}
		err = db.UpdateDate(next, id)
		if err != nil {
			writeJson(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
			return
		}
	}
	writeJson(w, http.StatusOK, map[string]any{})
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		writeJson(w, http.StatusBadRequest, map[string]any{"error": "Не указан ID задачи"})
		return
	}
	err := db.DeleteTask(id)
	if err != nil {
		writeJson(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
		return
	}
	writeJson(w, http.StatusOK, map[string]any{})
}
