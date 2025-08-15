package scheduler

import (
	"sync"
	"time"

	"github.com/go-co-op/gocron"
)

// SchedulerImpl implements the Scheduler interface
type SchedulerImpl struct {
	scheduler *gocron.Scheduler
	jobs      map[string]*gocron.Job
	mutex     sync.RWMutex
}

// NewScheduler creates a new scheduler instance
func NewScheduler() *SchedulerImpl {
	scheduler := gocron.NewScheduler(time.UTC)
	scheduler.StartAsync()

	return &SchedulerImpl{
		scheduler: scheduler,
		jobs:      make(map[string]*gocron.Job),
	}
}

// ScheduleJob schedules a job to run at a specific time
func (ss *SchedulerImpl) ScheduleJob(jobID string, scheduledTime time.Time, job func()) error {
	ss.mutex.Lock()
	defer ss.mutex.Unlock()

	// Calculate delay until scheduled time
	delay := scheduledTime.Sub(time.Now())
	if delay < 0 {
		delay = 0 // If scheduled time is in the past, run immediately
	}

	// Schedule the job
	scheduledJob, err := ss.scheduler.Every(delay).Do(job)
	if err != nil {
		return err
	}

	// Store the job reference
	ss.jobs[jobID] = scheduledJob

	return nil
}

// CancelJob cancels a scheduled job
func (ss *SchedulerImpl) CancelJob(jobID string) error {
	ss.mutex.Lock()
	defer ss.mutex.Unlock()

	if job, exists := ss.jobs[jobID]; exists {
		ss.scheduler.RemoveByReference(job)
		delete(ss.jobs, jobID)
		return nil
	}

	return nil // Job not found, consider it already cancelled
}

// GetScheduledJobs returns all scheduled job IDs
func (ss *SchedulerImpl) GetScheduledJobs() []string {
	ss.mutex.RLock()
	defer ss.mutex.RUnlock()

	jobIDs := make([]string, 0, len(ss.jobs))
	for jobID := range ss.jobs {
		jobIDs = append(jobIDs, jobID)
	}

	return jobIDs
}

// Stop stops the scheduler
func (ss *SchedulerImpl) Stop() {
	ss.scheduler.Stop()
}
