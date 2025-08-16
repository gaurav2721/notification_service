package scheduler

import (
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// SchedulerImpl implements the Scheduler interface
type SchedulerImpl struct {
	jobs  map[string]*time.Timer
	mutex sync.RWMutex
}

// NewScheduler creates a new scheduler instance
func NewScheduler() *SchedulerImpl {
	return &SchedulerImpl{
		jobs: make(map[string]*time.Timer),
	}
}

// ScheduleJob schedules a job to run at a specific time
func (ss *SchedulerImpl) ScheduleJob(jobID string, scheduledTime time.Time, job func()) error {
	ss.mutex.Lock()
	defer ss.mutex.Unlock()

	// Log scheduling information for debugging
	now := time.Now()
	delay := scheduledTime.Sub(now)

	logrus.WithFields(logrus.Fields{
		"job_id":         jobID,
		"scheduled_time": scheduledTime.Format("2006-01-02T15:04:05Z"),
		"current_time":   now.Format("2006-01-02T15:04:05Z"),
		"delay_seconds":  int(delay.Seconds()),
	}).Debug("Scheduling job")

	// If scheduled time is in the past, run immediately
	if delay <= 0 {
		logrus.WithField("job_id", jobID).Warn("Scheduled time is in the past, running immediately")
		go func() {
			job()
			logrus.WithField("job_id", jobID).Debug("Immediate job completed")
		}()
		return nil
	}

	// Create a timer for the scheduled job
	timer := time.AfterFunc(delay, func() {
		logrus.WithField("job_id", jobID).Info("Executing scheduled job")

		// Execute the original job
		job()

		// Remove the job from the jobs map after execution
		ss.mutex.Lock()
		delete(ss.jobs, jobID)
		ss.mutex.Unlock()

		logrus.WithField("job_id", jobID).Debug("Job removed from scheduler after execution")
	})

	// Store the timer reference
	ss.jobs[jobID] = timer

	return nil
}

// CancelJob cancels a scheduled job
func (ss *SchedulerImpl) CancelJob(jobID string) error {
	ss.mutex.Lock()
	defer ss.mutex.Unlock()

	if timer, exists := ss.jobs[jobID]; exists {
		timer.Stop()
		delete(ss.jobs, jobID)
		logrus.WithField("job_id", jobID).Debug("Job cancelled")
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
	ss.mutex.Lock()
	defer ss.mutex.Unlock()

	// Stop all timers
	for jobID, timer := range ss.jobs {
		timer.Stop()
		logrus.WithField("job_id", jobID).Debug("Job stopped during scheduler shutdown")
	}

	// Clear the jobs map
	ss.jobs = make(map[string]*time.Timer)
}
