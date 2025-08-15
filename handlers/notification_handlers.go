package handlers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gaurav2721/notification-service/models"
	"github.com/gin-gonic/gin"
)

// NotificationHandler handles HTTP requests for notifications
type NotificationHandler struct {
	notificationService interface {
		SendNotification(ctx context.Context, notification interface{}) (interface{}, error)
		ScheduleNotification(ctx context.Context, notification interface{}) (interface{}, error)
		GetNotificationStatus(ctx context.Context, notificationID string) (interface{}, error)
		CreateTemplate(ctx context.Context, template interface{}) (interface{}, error)
		GetTemplateVersion(ctx context.Context, templateID string, version int) (interface{}, error)
	}
}

// NewNotificationHandler creates a new notification handler
func NewNotificationHandler(notificationService interface {
	SendNotification(ctx context.Context, notification interface{}) (interface{}, error)
	ScheduleNotification(ctx context.Context, notification interface{}) (interface{}, error)
	GetNotificationStatus(ctx context.Context, notificationID string) (interface{}, error)
	CreateTemplate(ctx context.Context, template interface{}) (interface{}, error)
	GetTemplateVersion(ctx context.Context, templateID string, version int) (interface{}, error)
}) *NotificationHandler {
	return &NotificationHandler{
		notificationService: notificationService,
	}
}

// SendNotification handles POST /notifications
func (h *NotificationHandler) SendNotification(c *gin.Context) {
	var request struct {
		Type        string                 `json:"type" binding:"required"`
		Content     map[string]interface{} `json:"content"`
		Template    *models.TemplateData   `json:"template,omitempty"`
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
		Template    *models.TemplateData
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
	var request models.TemplateRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate template content
	if err := request.Content.ValidateTemplateContent(request.Type); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create template object
	template := &models.Template{
		Name:              request.Name,
		Type:              request.Type,
		Content:           request.Content,
		RequiredVariables: request.RequiredVariables,
		Description:       request.Description,
	}

	response, err := h.notificationService.CreateTemplate(c.Request.Context(), template)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, response)
}

// GetTemplateVersion handles GET /templates/:templateId/versions/:version
func (h *NotificationHandler) GetTemplateVersion(c *gin.Context) {
	templateID := c.Param("templateId")
	if templateID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "template ID is required"})
		return
	}

	versionStr := c.Param("version")
	if versionStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "version is required"})
		return
	}

	version, err := strconv.Atoi(versionStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid version number"})
		return
	}

	template, err := h.notificationService.GetTemplateVersion(c.Request.Context(), templateID, version)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, template)
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
