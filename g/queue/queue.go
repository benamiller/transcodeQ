package queue

import (
	"github.com/benamiller/transcodeQ/g/models"
	"sync"
)

type JobQueue struct {
	jobs map[string]models.TranscodeJob
	mu sync.Mutex
}

func NewJobQueue() *JobQueue {
	return &JobQueue{
		jobs: make(map[string]models.TranscodeJob),
	}
}

func (q *JobQueue) AddJob(job models.TranscodeJob) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.jobs[job.ID] = job
}

func (q *JobQueue) GetJob(id string) (models.TranscodeJob, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()
	job, ok := q.jobs[id]
	return job, ok
}

func (q *JobQueue) ListJobs() []models.TranscodeJob {
	q.mu.Lock()
	defer q.mu.Unlock()
	var list []models.TranscodeJob
	for _, job := range q.jobs {
		list = append(list, job)
	}
	return list
}
