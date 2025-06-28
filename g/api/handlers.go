package api

import (
	"encoding/json"
	"net/http"
	"github.com/benamiller/transcodeQ/g/models"
	"github.com/benamiller/transcodeQ/g/queue"
)

type API struct {
	Queue *queue.JobQueue
}

func (api *API) JobsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		api.CreateJobHandler(w, r)
	case "GET":
		id := r.URL.Query().Get("id")
		if id != "" {
			api.GetJobHandler(w, r)
		} else {
			api.ListJobsHandler(w, r)
		}
	default:
		http.Error(w, "Method not supported", http.StatusMethodNotAllowed)
	}

}

func (api *API) CreateJobHandler(w http.ResponseWriter, r *http.Request) {
	var job models.TranscodeJob
	err := json.NewDecoder(r.Body).Decode(&job)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	api.Queue.AddJob(job)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(job)
}
	
func (api *API) GetJobHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	job, ok := api.Queue.GetJob(id)
	if !ok {
		http.Error(w, "Job not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(job)
}

func (api *API) ListJobsHandler(w http.ResponseWriter, r *http.Request) {
	jobs, ok := api.Queue.ListJobs()
	if !ok {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}

	json.NewEncoder(w).Encode(jobs)
}
	
