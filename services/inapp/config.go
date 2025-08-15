package inapp

// InAppConfig holds in-app service configuration
type InAppConfig struct {
	DatabaseURL string
	MaxRetries  int
	RetryDelay  int
}

// DefaultInAppConfig returns default in-app configuration
func DefaultInAppConfig() *InAppConfig {
	return &InAppConfig{
		DatabaseURL: "memory://",
		MaxRetries:  3,
		RetryDelay:  1000,
	}
}
