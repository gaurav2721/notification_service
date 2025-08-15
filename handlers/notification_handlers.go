package handlers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// NotificationHandler handles HTTP requests for notifications
type NotificationHandler struct {
	notificationService interface {
		SendNotification(ctx context.Context, notification interface{}) (interface{}, error)
		ScheduleNotification(ctx context.Context, notification interface{}) (interface{}, error)
		GetNotificationStatus(ctx context.Context, notificationID string) (interface{}, error)
		CreateTemplate(ctx context.Context, template interface{}) error
		GetTemplate(ctx context.Context, templateID string) (interface{}, error)
		UpdateTemplate(ctx context.Context, template interface{}) error
		DeleteTemplate(ctx context.Context, templateID string) error
	}
}

// NewNotificationHandler creates a new notification handler
func NewNotificationHandler(notificationService interface {
	SendNotification(ctx context.Context, notification interface{}) (interface{}, error)
	ScheduleNotification(ctx context.Context, notification interface{}) (interface{}, error)
	GetNotificationStatus(ctx context.Context, notificationID string) (interface{}, error)
	CreateTemplate(ctx context.Context, template interface{}) error
	GetTemplate(ctx context.Context, templateID string) (interface{}, error)
	UpdateTemplate(ctx context.Context, template interface{}) error
	DeleteTemplate(ctx context.Context, templateID string) error
}) *NotificationHandler {
	return &NotificationHandler{
		notificationService: notificationService,
	}
}

// SendNotification handles POST /notifications
func (h *NotificationHandler) SendNotification(c *gin.Context) {
	type TemplateData struct {
		ID   string                 `json:"id"`
		Data map[string]interface{} `json:"data"`
	}

	var request struct {
		Type        string                 `json:"type" binding:"required"`
		Content     map[string]interface{} `json:"content"`
		Template    *TemplateData          `json:"template,omitempty"`
		Recipients  []string               `json:"recipients" binding:"required"`
		Metadata    map[string]interface{} `json:"metadata"`
		ScheduledAt *time.Time             `json:"scheduled_at"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create notification object
	notification := &struct {
		ID          string
		Type        string
		Content     map[string]interface{}
		Template    *TemplateData
		Recipients  []string
		Metadata    map[string]interface{}
		ScheduledAt *time.Time
	}{
		ID:          generateID(),
		Type:        request.Type,
		Content:     request.Content,
		Template:    request.Template,
		Recipients:  request.Recipients,
		Metadata:    request.Metadata,
		ScheduledAt: request.ScheduledAt,
	}

	// Send notification
	response, err := h.notificationService.SendNotification(c.Request.Context(), notification)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetNotificationStatus handles GET /notifications/:id
func (h *NotificationHandler) GetNotificationStatus(c *gin.Context) {
	notificationID := c.Param("id")
	if notificationID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "notification ID is required"})
		return
	}

	response, err := h.notificationService.GetNotificationStatus(c.Request.Context(), notificationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// CreateTemplate handles POST /templates
func (h *NotificationHandler) CreateTemplate(c *gin.Context) {
	var request struct {
		Name      string   `json:"name" binding:"required"`
		Type      string   `json:"type" binding:"required"`
		Subject   string   `json:"subject" binding:"required"`
		Body      string   `json:"body" binding:"required"`
		Variables []string `json:"variables"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	template := &struct {
		ID        string
		Name      string
		Type      string
		Subject   string
		Body      string
		Variables []string
		CreatedAt time.Time
		UpdatedAt time.Time
	}{
		ID:        generateID(),
		Name:      request.Name,
		Type:      request.Type,
		Subject:   request.Subject,
		Body:      request.Body,
		Variables: request.Variables,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := h.notificationService.CreateTemplate(c.Request.Context(), template)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, template)
}

// GetTemplate handles GET /templates/:id
func (h *NotificationHandler) GetTemplate(c *gin.Context) {
	templateID := c.Param("id")
	if templateID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "template ID is required"})
		return
	}

	template, err := h.notificationService.GetTemplate(c.Request.Context(), templateID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, template)
}

// UpdateTemplate handles PUT /templates/:id
func (h *NotificationHandler) UpdateTemplate(c *gin.Context) {
	templateID := c.Param("id")
	if templateID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "template ID is required"})
		return
	}

	var request struct {
		Name      string   `json:"name"`
		Type      string   `json:"type"`
		Subject   string   `json:"subject"`
		Body      string   `json:"body"`
		Variables []string `json:"variables"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	template := &struct {
		ID        string
		Name      string
		Type      string
		Subject   string
		Body      string
		Variables []string
		UpdatedAt time.Time
	}{
		ID:        templateID,
		Name:      request.Name,
		Type:      request.Type,
		Subject:   request.Subject,
		Body:      request.Body,
		Variables: request.Variables,
		UpdatedAt: time.Now(),
	}

	err := h.notificationService.UpdateTemplate(c.Request.Context(), template)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, template)
}

// DeleteTemplate handles DELETE /templates/:id
func (h *NotificationHandler) DeleteTemplate(c *gin.Context) {
	templateID := c.Param("id")
	if templateID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "template ID is required"})
		return
	}

	err := h.notificationService.DeleteTemplate(c.Request.Context(), templateID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Template deleted successfully"})
}

// HealthCheck handles GET /health
func (h *NotificationHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"timestamp": time.Now(),
		"service":   "notification-service",
	})
}

// Helper function to generate a simple ID
func generateID() string {
	return strconv.FormatInt(time.Now().UnixNano(), 10)
}
