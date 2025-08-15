package scheduler

import "errors"

// Scheduler service errors
var (
	ErrSchedulingFailed = errors.New("failed to schedule notification")
	ErrJobNotFound      = errors.New("scheduled job not found")
	ErrJobTimeout       = errors.New("job execution timeout")
)
