package scheduler

import "time"

// SchedulerService interface defines methods for notification scheduling
type SchedulerService interface {
	ScheduleJob(jobID string, scheduledTime time.Time, job func()) error
	CancelJob(jobID string) error
}
