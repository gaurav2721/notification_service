package user

// UserConfig holds user service configuration
type UserConfig struct {
	DatabaseURL string
	CacheTTL    int
}

// DefaultUserConfig returns default user configuration
func DefaultUserConfig() *UserConfig {
	return &UserConfig{
		DatabaseURL: "memory://",
		CacheTTL:    3600,
	}
}
