package scheduler

import (
	"errors"
	"time"
)

// SchedulerService interface defines methods for notification scheduling
type SchedulerService interface {
	ScheduleJob(jobID string, scheduledTime time.Time, job func()) error
	CancelJob(jobID string) error
}

// SchedulerConfig holds scheduler service configuration
type SchedulerConfig struct {
	MaxConcurrentJobs int
	JobTimeout        int
	RetentionDays     int
}

// DefaultSchedulerConfig returns default scheduler configuration
func DefaultSchedulerConfig() *SchedulerConfig {
	return &SchedulerConfig{
		MaxConcurrentJobs: 10,
		JobTimeout:        300,
		RetentionDays:     30,
	}
}

// Scheduler service errors
var (
	ErrSchedulingFailed = errors.New("failed to schedule notification")
	ErrJobNotFound      = errors.New("scheduled job not found")
	ErrJobTimeout       = errors.New("job execution timeout")
)
