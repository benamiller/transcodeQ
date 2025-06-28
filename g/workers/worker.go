package workers

import (
	"fmt"
	"math"
	"time"
	"github.com/benamiller/transcodeQ/g/models"
	"github.com/benamiller/transcodeQ/g/queue"
)

func getIterationsForFormat(format string) int {
	var formatIterations = map[string]int{
		"240p": 10_000,
		"360p": 30_000,
		"480p": 50_000,
		"720p": 100_000,
		"1080p": 200_000,
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

		job.StatusMap[format] = models.StatusCompleted

		q.AddJob(job)
	}
}
