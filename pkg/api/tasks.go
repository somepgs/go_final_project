package api

import (
	"github.com/somepgs/go_final_project/pkg/db"
	"net/http"
)

const limitTasks = 50 // limitTasks defines the maximum number of tasks to return in a single request.

type tasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	search := r.FormValue("search")
	if search != "" {
		tasks, err := db.SearchTasks(search, limitTasks)
		if err != nil {
			writeJson(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
			return
		}
		if len(tasks) == 0 {
			writeJson(w, http.StatusOK, tasksResp{Tasks: []*db.Task{}})
			return
		}
		writeJson(w, http.StatusOK, tasksResp{Tasks: tasks})
		return
	}
	tasks, err := db.Tasks(limitTasks)
	if err != nil {
		writeJson(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
		return
	}
	if len(tasks) == 0 {
		writeJson(w, http.StatusOK, tasksResp{Tasks: []*db.Task{}})
		return
	}
	writeJson(w, http.StatusOK, tasksResp{Tasks: tasks})
}
