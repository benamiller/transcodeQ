package queue

import (
	"github.com/benamiller/transcodeQ/g/models"
	"sync"
)

type JobQueue struct {
	jobs map[string]models.TranscodeJob
	mu sync.Mutex
}
