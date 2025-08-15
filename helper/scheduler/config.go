package scheduler

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
