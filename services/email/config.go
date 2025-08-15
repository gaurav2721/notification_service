package email

// EmailConfig holds email service configuration
type EmailConfig struct {
	SMTPHost     string
	SMTPPort     int
	SMTPUsername string
	SMTPPassword string
	FromEmail    string
	FromName     string
}

// DefaultEmailConfig returns default email configuration
func DefaultEmailConfig() *EmailConfig {
	return &EmailConfig{
		SMTPHost:     "localhost",
		SMTPPort:     587,
		SMTPUsername: "",
		SMTPPassword: "",
		FromEmail:    "noreply@company.com",
		FromName:     "Notification Service",
	}
}
