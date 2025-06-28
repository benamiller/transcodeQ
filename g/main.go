package main

import (
	"net/http"
	"github.com/benamiller/transcodeQ/g/queue"
	"github.com/benamiller/transcodeQ/g/api"
)

func main() {
	jobQueue := queue.NewJobQueue()
	apiHandler := &api.API{Queue: jobQueue}

	http.HandleFunc("/jobs", apiHandler.JobsHandler)

	http.ListenAndServe(":8080", nil)
}
