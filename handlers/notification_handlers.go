package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gaurav2721/notification-service/models"
	"github.com/gaurav2721/notification-service/notification_manager"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// NotificationHandler handles HTTP requests for notifications
type NotificationHandler struct {
	notificationService notification_manager.NotificationManager
}

// NewNotificationHandler creates a new notification handler
func NewNotificationHandler(
	notificationService notification_manager.NotificationManager,
) *NotificationHandler {
	return &NotificationHandler{
		notificationService: notificationService,
	}
}

// SendNotification handles POST /notifications
func (h *NotificationHandler) SendNotification(c *gin.Context) {
	logrus.Debug("Received notification send request")

	// Get validated request from middleware
	validatedRequestInterface, exists := c.Get("validated_request")
	if !exists {
		logrus.Error("Validated request not found in context")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	requestPtr, ok := validatedRequestInterface.(*models.NotificationRequest)
	if !ok {
		logrus.Error("Failed to cast validated request to NotificationRequest")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// Dereference the pointer to get the actual request
	request := *requestPtr

	logrus.WithFields(logrus.Fields{
		"type":        request.Type,
		"recipients":  len(request.Recipients),
		"scheduled":   request.ScheduledAt != nil,
		"hasTemplate": request.Template != nil,
		"hasFrom":     request.From != nil,
	}).Debug("Processing notification request")

	// Process the notification request through the notification manager
	response, err := h.notificationService.ProcessNotificationRequest(c.Request.Context(), &request)
	if err != nil {
		logrus.WithError(err).Error("Failed to process notification request")
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
	// Get validated request from middleware
	validatedRequestInterface, exists := c.Get("validated_template_request")
	if !exists {
		logrus.Error("Validated template request not found in context")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	requestPtr, ok := validatedRequestInterface.(*models.TemplateRequest)
	if !ok {
		logrus.Error("Failed to cast validated request to TemplateRequest")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// Dereference the pointer to get the actual request
	request := *requestPtr

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

// GetPredefinedTemplates handles GET /templates/predefined
func (h *NotificationHandler) GetPredefinedTemplates(c *gin.Context) {
	templates := h.notificationService.GetPredefinedTemplates()

	// Convert to response format
	var response []map[string]interface{}
	for _, template := range templates {
		response = append(response, map[string]interface{}{
			"id":                 template.ID,
			"name":               template.Name,
			"type":               string(template.Type),
			"version":            template.Version,
			"content":            template.Content,
			"description":        template.Description,
			"required_variables": template.RequiredVariables,
			"status":             template.Status,
			"created_at":         template.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"templates": response,
		"count":     len(response),
	})
}

// GetTemplateVersion handles GET /templates/:templateId/versions/:version
func (h *NotificationHandler) GetTemplateVersion(c *gin.Context) {
	templateID := c.Param("templateId")
	versionStr := c.Param("version")

	// Parameters are already validated by middleware
	version, err := strconv.Atoi(versionStr)
	if err != nil {
		// This should not happen as middleware validates this
		logrus.WithError(err).Error("Failed to parse version parameter")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
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
