package scan

import "time"

type JobStatus string

const (
	JobRunning JobStatus = "running"
	JobDone    JobStatus = "done"
	JobFailed  JobStatus = "failed"
)

type Job struct {
	ID         string     `json:"id"`
	LibraryID  int64      `json:"library_id"`
	StartedAt  time.Time  `json:"started_at"`
	FinishedAt *time.Time `json:"finished_at,omitempty"`
	Status     JobStatus  `json:"status"`
	Error      string     `json:"error,omitempty"`
}
