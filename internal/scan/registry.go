package scan

import (
	"sync"
	"time"

	"github.com/google/uuid"
)

type Registry struct {
	mu   sync.RWMutex
	jobs map[string]*Job
}

func NewRegistry() *Registry {
	return &Registry{
		jobs: make(map[string]*Job),
	}
}

func (r *Registry) Start(libraryID int64) *Job {
	job := &Job{
		ID:        uuid.NewString(),
		LibraryID: libraryID,
		StartedAt: time.Now(),
		Status:    JobRunning,
	}

	r.mu.Lock()
	r.jobs[job.ID] = job
	r.mu.Unlock()

	return job
}

func (r *Registry) Finish(jobID string) {
	now := time.Now()

	r.mu.Lock()
	if job, ok := r.jobs[jobID]; ok {
		job.Status = JobDone
		job.FinishedAt = &now
	}
	r.mu.Unlock()
}

func (r *Registry) Fail(jobID string, err error) {
	now := time.Now()

	r.mu.Lock()
	if job, ok := r.jobs[jobID]; ok {
		job.Status = JobFailed
		job.Error = err.Error()
		job.FinishedAt = &now
	}
	r.mu.Unlock()
}

func (r *Registry) Get(jobID string) (*Job, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	job, ok := r.jobs[jobID]
	return job, ok
}
