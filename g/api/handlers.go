package api

import (
	"encoding/json"
	"net/http"
	"sync"
	"github.com/benamiller/transcodeQ/g/models"
	"github.com/benamiller/transcodeQ/g/queue"
)

type API struct {
	Queue *queue.JobQueue
	mu	sync.Mutex
	nextID int
}

func (api *API) NextID() {
	api.mu.Lock()
	defer api.mu.Unlock()

	api.nextID++
	return fmt.Sprintf("%d", api.nextID)
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
	var req models.CreateTranscodeJobRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newID := api.NextID()
	newStatusMap := make(map[string]models.JobStatus)
	for _, format := range req.formats {
		newStatusMap[format] = models.StatusQueued
	}

	job := models.TranscodeJob {
		ID: newID,
		Title: req.title,
		Formats: req.formats,
		StatusMap: newStatusMap,
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
	jobs := api.Queue.ListJobs()

	json.NewEncoder(w).Encode(jobs)
}
	
