package workers

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"time"
	"github.com/benamiller/transcodeQ/g/models"
	"github.com/benamiller/transcodeQ/g/queue"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func getIterationsForFormat(format string) int {
	var formatIterations = map[string]int{
		"240p": 1000,
		"360p": 5000,
		"480p": 20_000,
		"720p": 50_000,
		"1080p": 100_000,
	}

	if iters, ok := formatIterations[format]; ok {
		return iters
	}
	return 50_000
}

func performWorkForFormat(format string) {
	iterations := getIterationsForFormat(format)
	
	var sum float64
	for i := 0; i < iterations; i++ {
		sum += math.Sqrt(float64(i)) * math.Sin(float64(i))
	}

	time.Sleep(500 * time.Millisecond)
}

func shouldFail() bool {
	return rand.Float64() < 0.05
}

func retryFormat(job models.TranscodeJob, format string) (models.TranscodeJob, error) {
	fmt.Printf("Job %s, format %s FAILED. Retrying %d more time(s)\n", job.ID, format, job.Retries)
	job.Retries = job.Retries - 1
	if job.Retries < 0 {
		return job, errors.New("Exhausted all retries")
	} else {
		performWorkForFormat(format)

		if (shouldFail()) {
			retryFormat(job, format)
		} else {
			job.StatusMap[format] = models.StatusCompleted
			return job, nil
		}
	}
	return job, errors.New("Exhausted all retries")
}

func ProcessJob(jobID string, q *queue.JobQueue) {
	job, ok := q.GetJob(jobID)
	if !ok {
		fmt.Printf("Job %s not found\n", jobID)
		return
	}

	for _, format := range job.Formats {
		fmt.Printf("Job %s: processing format %s\n", jobID, format)

		job.StatusMap[format] = models.StatusProcessing
		q.AddJob(job)

		performWorkForFormat(format)

		if (shouldFail()) {
			job.StatusMap[format] = models.StatusRetrying
			job, err := retryFormat(job, format)
			if err != nil {
				job.StatusMap[format] = models.StatusCompleted
			} else {
				job.StatusMap[format] = models.StatusFailed
				fmt.Printf("Job %s, format %s FAILED\n", jobID, format)
			}
		} else {
			fmt.Printf("Job %s, format %s succeeded\n", jobID, format)
			job.StatusMap[format] = models.StatusCompleted
		}

		q.AddJob(job)
	}
}
