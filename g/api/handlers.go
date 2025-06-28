package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"github.com/benamiller/transcodeQ/g/models"
	"github.com/benamiller/transcodeQ/g/queue"
	"github.com/benamiller/transcodeQ/g/workers"
)

type API struct {
	Queue *queue.JobQueue
	mu	sync.Mutex
	nextID int
}

func (api *API) NextID() string {
	api.mu.Lock()
	defer api.mu.Unlock()

	api.nextID++
	return fmt.Sprintf("%d", api.nextID)
}

func (api *API) JobsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		api.CreateJobHandler(w, r)
		return
	case "GET":
		id := r.URL.Query().Get("id")
		format := r.URL.Query().Get("format")

		if id != "" && format != "" {
			api.GetJobFormatStatusHandler(w, r)
			return
		}
		if id != "" {
			api.GetJobHandler(w, r)
			return
		}
		api.ListJobsHandler(w, r)
		return
	default:
		http.Error(w, "Method not supported", http.StatusMethodNotAllowed)
		return
	}

}

func (api *API) CreateJobHandler(w http.ResponseWriter, r *http.Request) {
	var req models.CreateTranscodeJobRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newID := api.NextID()
	newStatusMap := make(map[string]models.JobStatus)
	for _, format := range req.Formats {
		newStatusMap[format] = models.StatusQueued
	}

	job := models.TranscodeJob {
		ID: newID,
		Title: req.Title,
		Formats: req.Formats,
		StatusMap: newStatusMap,
		Retries: 3,
	}

	api.Queue.AddJob(job)

	go workers.ProcessJob(newID, api.Queue)

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

func (api *API) GetJobFormatStatusHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	format := r.URL.Query().Get("format")

	job, ok := api.Queue.GetJob(id)
	if !ok {
		http.Error(w, "Job not found", http.StatusNotFound)
		return
	}

	status := job.StatusMap[format]

	json.NewEncoder(w).Encode(status)
}

func (api *API) ListJobsHandler(w http.ResponseWriter, r *http.Request) {
	jobs := api.Queue.ListJobs()

	json.NewEncoder(w).Encode(jobs)
}
	
