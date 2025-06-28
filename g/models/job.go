package models

type JobStatus string

const (
	StatusQueued	JobStatus = "queued"
	StatusProcessing	JobStatus = "processing"
	StatusCompleted JobStatus = "completed"
	StatusFailed	JobStatus = "failed"
)

type CreateTranscodeJobRequest struct {
	Title	string	`json:"title"`
	Formats	[]string	`json:"formats"`
}
	

type TranscodeJob struct {
	ID	string	`json:"id"`
	Title	string	`json:"title"`
	Formats	[]string	`json:"formats"`
	StatusMap map[string]JobStatus `json:"status_map"`
}
