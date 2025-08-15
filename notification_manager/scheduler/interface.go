package scheduler

import "time"

// Scheduler interface defines methods for notification scheduling
type Scheduler interface {
	ScheduleJob(jobID string, scheduledTime time.Time, job func()) error
	CancelJob(jobID string) error
}
