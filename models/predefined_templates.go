package models

import (
	"time"
)

// PredefinedTemplates returns a list of templates that should be loaded on service startup
func PredefinedTemplates() []*Template {
	return []*Template{
		// Email Templates
		welcomeEmailTemplate(),
		passwordResetTemplate(),
		orderConfirmationTemplate(),

		// Slack Templates
		systemAlertTemplate(),
		deploymentNotificationTemplate(),

		// In-App Templates
		orderStatusUpdateTemplate(),
		paymentReminderTemplate(),
	}
}

// welcomeEmailTemplate creates the welcome email template
func welcomeEmailTemplate() *Template {
	return &Template{
		ID:   "550e8400-e29b-41d4-a716-446655440000", // Fixed UUID for consistency
		Name: "Welcome Email Template",
		Type: EmailNotification,
		Content: TemplateContent{
			Subject:   "Welcome to {{platform}}, {{name}}!",
			EmailBody: "Hello {{name}},\n\nWelcome to {{platform}}! We are excited to have you on board.\n\nYour account has been successfully created with the following details:\n- Username: {{username}}\n- Email: {{email}}\n- Account Type: {{account_type}}\n\nPlease click the following link to activate your account:\n{{activation_link}}\n\nIf you have any questions, please contact our support team.\n\nBest regards,\nThe {{platform}} Team",
		},
		RequiredVariables: []string{"name", "platform", "username", "email", "account_type", "activation_link"},
		Description:       "Welcome email template for new user onboarding",
		Version:           1,
		Status:            "active",
		CreatedAt:         time.Now(),
	}
}

// passwordResetTemplate creates the password reset template
func passwordResetTemplate() *Template {
	return &Template{
		ID:   "550e8400-e29b-41d4-a716-446655440001", // Fixed UUID for consistency
		Name: "Password Reset Template",
		Type: EmailNotification,
		Content: TemplateContent{
			Subject:   "Password Reset Request - {{platform}}",
			EmailBody: "Hello {{name}},\n\nWe received a request to reset your password for your {{platform}} account.\n\nIf you made this request, please click the link below to reset your password:\n{{reset_link}}\n\nThis link will expire in {{expiry_hours}} hours.\n\nIf you did not request a password reset, please ignore this email.\n\nBest regards,\nThe {{platform}} Security Team",
		},
		RequiredVariables: []string{"name", "platform", "reset_link", "expiry_hours"},
		Description:       "Password reset email template",
		Version:           1,
		Status:            "active",
		CreatedAt:         time.Now(),
	}
}

// orderConfirmationTemplate creates the order confirmation template
func orderConfirmationTemplate() *Template {
	return &Template{
		ID:   "550e8400-e29b-41d4-a716-446655440002", // Fixed UUID for consistency
		Name: "Order Confirmation Template",
		Type: EmailNotification,
		Content: TemplateContent{
			Subject:   "Order Confirmed - #{{order_id}}",
			EmailBody: "Dear {{customer_name}},\n\nThank you for your order! Your order has been confirmed and is being processed.\n\n**Order Details:**\n- Order ID: #{{order_id}}\n- Order Date: {{order_date}}\n- Total Amount: ${{total_amount}}\n- Payment Method: {{payment_method}}\n\n**Items Ordered:**\n{{items_list}}\n\n**Shipping Information:**\n{{shipping_address}}\n\n**Estimated Delivery:** {{delivery_date}}\n\nTrack your order: {{tracking_link}}\n\nIf you have any questions, please contact our support team.\n\nBest regards,\nThe {{platform}} Team",
		},
		RequiredVariables: []string{"customer_name", "order_id", "order_date", "total_amount", "payment_method", "items_list", "shipping_address", "delivery_date", "tracking_link", "platform"},
		Description:       "Order confirmation email template",
		Version:           1,
		Status:            "active",
		CreatedAt:         time.Now(),
	}
}

// systemAlertTemplate creates the system alert template
func systemAlertTemplate() *Template {
	return &Template{
		ID:   "550e8400-e29b-41d4-a716-446655440003", // Fixed UUID for consistency
		Name: "System Alert Template",
		Type: SlackNotification,
		Content: TemplateContent{
			Text: "🚨 *{{alert_type}} Alert*\n\n*System:* {{system_name}}\n*Severity:* {{severity}}\n*Environment:* {{environment}}\n*Message:* {{message}}\n*Timestamp:* {{timestamp}}\n*Action Required:* {{action_required}}\n\n*Affected Services:* {{affected_services}}\n*Dashboard:* {{dashboard_link}}\n\nPlease take immediate action if this is a critical alert.",
		},
		RequiredVariables: []string{"alert_type", "system_name", "severity", "environment", "message", "timestamp", "action_required", "affected_services", "dashboard_link"},
		Description:       "Slack alert template for system monitoring",
		Version:           1,
		Status:            "active",
		CreatedAt:         time.Now(),
	}
}

// deploymentNotificationTemplate creates the deployment notification template
func deploymentNotificationTemplate() *Template {
	return &Template{
		ID:   "550e8400-e29b-41d4-a716-446655440004", // Fixed UUID for consistency
		Name: "Deployment Notification Template",
		Type: SlackNotification,
		Content: TemplateContent{
			Text: "🚀 *Deployment {{status}}*\n\n*Service:* {{service_name}}\n*Environment:* {{environment}}\n*Version:* {{version}}\n*Deployed By:* {{deployed_by}}\n*Duration:* {{duration}}\n\n*Changes:*\n{{changes_summary}}\n\n*Rollback:* {{rollback_command}}\n*Monitoring:* {{monitoring_link}}",
		},
		RequiredVariables: []string{"status", "service_name", "environment", "version", "deployed_by", "duration", "changes_summary", "rollback_command", "monitoring_link"},
		Description:       "Slack notification template for deployment events",
		Version:           1,
		Status:            "active",
		CreatedAt:         time.Now(),
	}
}

// orderStatusUpdateTemplate creates the order status update template
func orderStatusUpdateTemplate() *Template {
	return &Template{
		ID:   "550e8400-e29b-41d4-a716-446655440005", // Fixed UUID for consistency
		Name: "Order Status Update Template",
		Type: InAppNotification,
		Content: TemplateContent{
			Title: "Order #{{order_id}} - {{status}}",
			Body:  "Your order has been {{status}}.\n\n*Order Details:*\n- Items: {{item_count}} items\n- Total: ${{total_amount}}\n- Status: {{status}}\n\n{{status_message}}\n\n{{action_button}}",
		},
		RequiredVariables: []string{"order_id", "status", "item_count", "total_amount", "status_message", "action_button"},
		Description:       "In-app notification template for order status updates",
		Version:           1,
		Status:            "active",
		CreatedAt:         time.Now(),
	}
}

// paymentReminderTemplate creates the payment reminder template
func paymentReminderTemplate() *Template {
	return &Template{
		ID:   "550e8400-e29b-41d4-a716-446655440006", // Fixed UUID for consistency
		Name: "Payment Reminder Template",
		Type: InAppNotification,
		Content: TemplateContent{
			Title: "Payment Due - ${{amount}}",
			Body:  "Your payment of ${{amount}} is due on {{due_date}}.\n\n*Invoice Details:*\n- Invoice #: {{invoice_id}}\n- Due Date: {{due_date}}\n- Amount: ${{amount}}\n\nPlease update your payment method or contact support if you have any questions.",
		},
		RequiredVariables: []string{"amount", "due_date", "invoice_id"},
		Description:       "In-app notification template for payment reminders",
		Version:           1,
		Status:            "active",
		CreatedAt:         time.Now(),
	}
}

// GetTemplateByID returns a predefined template by ID
func GetTemplateByID(templateID string) *Template {
	templates := PredefinedTemplates()
	for _, template := range templates {
		if template.ID == templateID {
			return template
		}
	}
	return nil
}

// GetTemplateByName returns a predefined template by name
func GetTemplateByName(name string) *Template {
	templates := PredefinedTemplates()
	for _, template := range templates {
		if template.Name == name {
			return template
		}
	}
	return nil
}

// GetTemplatesByType returns all predefined templates of a specific type
func GetTemplatesByType(templateType NotificationType) []*Template {
	templates := PredefinedTemplates()
	var filteredTemplates []*Template

	for _, template := range templates {
		if template.Type == templateType {
			filteredTemplates = append(filteredTemplates, template)
		}
	}

	return filteredTemplates
}
