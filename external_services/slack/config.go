package slack

// SlackConfig holds Slack service configuration
type SlackConfig struct {
	BotToken       string
	DefaultChannel string
}

// DefaultSlackConfig returns default Slack configuration
func DefaultSlackConfig() *SlackConfig {
	return &SlackConfig{
		BotToken:       "",
		DefaultChannel: "#general",
	}
}
