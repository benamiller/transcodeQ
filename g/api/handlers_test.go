package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/benamiller/transcodeQ/g/models"
	"github.com/benamiller/transcodeQ/g/queue"
)

func TestCreateJobHandler(t *testing.T) {
	q := queue.NewJobQueue()
	apiHandler := &API{
		Queue: q,
	}

	createReq := models.CreateTranscodeJobRequest{
		Title:	"My video",
		Formats: []string{"720p", "1080p"},
	}
	bodyBytes, _ := json.Marshal(createReq)

	req := httptest.NewRequest("POST", "/jobs", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	apiHandler.JobsHandler(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected 201, get %d", rr.Code)
	}

	var job models.TranscodeJob
	err := json.NewDecoder(rr.Body).Decode(&job)
	if err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if job.Title != "My video" {
		t.Errorf("expected title 'My video', got %s", job.Title)
	}

	if len(job.Formats) != 2 {
		t.Errorf("expected 2 formats, got %d", len(job.Formats))
	}

	if job.StatusMap["720p"] != models.StatusQueued {
		t.Errorf("expected 720p queue, got %s", job.StatusMap["720p"])
	}
}


