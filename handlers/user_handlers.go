package handlers

import (
	"net/http"

	"github.com/gaurav2721/notification-service/external_services/user"
	"github.com/gaurav2721/notification-service/models"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// UserHandler handles HTTP requests for user management
type UserHandler struct {
	userService user.UserService
}

// NewUserHandler creates a new user handler
func NewUserHandler(userService user.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// GetUsers handles GET /api/v1/users
func (h *UserHandler) GetUsers(c *gin.Context) {
	logrus.Debug("Received get users request")
	ctx := c.Request.Context()

	users, err := h.userService.GetAllUsers(ctx)
	if err != nil {
		logrus.WithError(err).Error("Failed to get all users")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logrus.WithField("user_count", len(users)).Debug("Retrieved users successfully")
	c.JSON(http.StatusOK, gin.H{
		"users": users,
		"count": len(users),
	})
}

// GetUser handles GET /api/v1/users/:id
func (h *UserHandler) GetUser(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		logrus.Warn("Get user request missing user ID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "user ID is required"})
		return
	}

	logrus.WithField("user_id", userID).Debug("Received get user request")
	ctx := c.Request.Context()
	user, err := h.userService.GetUserByID(ctx, userID)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"user_id": userID,
			"error":   err.Error(),
		}).Error("Failed to get user by ID")
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	logrus.WithField("user_id", userID).Debug("User retrieved successfully")
	c.JSON(http.StatusOK, user)
}

// CreateUser handles POST /api/v1/users
func (h *UserHandler) CreateUser(c *gin.Context) {
	logrus.Debug("Received create user request")

	var request struct {
		Email    string `json:"email" binding:"required,email"`
		FullName string `json:"full_name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		logrus.WithError(err).Warn("Invalid request body for create user")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	logrus.WithFields(logrus.Fields{
		"email":     request.Email,
		"full_name": request.FullName,
	}).Debug("Creating new user")

	// Create new user using the models.NewUser function
	newUser := models.NewUser(request.Email, request.FullName)

	ctx := c.Request.Context()
	err := h.userService.CreateUser(ctx, newUser)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"email": request.Email,
			"error": err.Error(),
		}).Error("Failed to create user")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logrus.WithFields(logrus.Fields{
		"user_id": newUser.ID,
		"email":   request.Email,
	}).Debug("User created successfully")
	c.JSON(http.StatusCreated, newUser)
}

// UpdateUser handles PUT /api/v1/users/:id
func (h *UserHandler) UpdateUser(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		logrus.Warn("Update user request missing user ID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "user ID is required"})
		return
	}

	logrus.WithField("user_id", userID).Debug("Received update user request")

	var request struct {
		Email        string `json:"email"`
		FullName     string `json:"full_name"`
		SlackUserID  string `json:"slack_user_id"`
		SlackChannel string `json:"slack_channel"`
		PhoneNumber  string `json:"phone_number"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		logrus.WithError(err).Warn("Invalid request body for update user")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get existing user first
	ctx := c.Request.Context()
	existingUser, err := h.userService.GetUserByID(ctx, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Update fields if provided
	if request.Email != "" {
		existingUser.Email = request.Email
	}
	if request.FullName != "" {
		existingUser.FullName = request.FullName
	}
	if request.SlackUserID != "" {
		existingUser.SlackUserID = request.SlackUserID
	}
	if request.SlackChannel != "" {
		existingUser.SlackChannel = request.SlackChannel
	}
	if request.PhoneNumber != "" {
		existingUser.PhoneNumber = request.PhoneNumber
	}

	// Update user
	err = h.userService.UpdateUser(ctx, existingUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, existingUser)
}

// DeleteUser handles DELETE /api/v1/users/:id
func (h *UserHandler) DeleteUser(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user ID is required"})
		return
	}

	ctx := c.Request.Context()
	err := h.userService.DeleteUser(ctx, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

// GetUserNotificationInfo handles GET /api/v1/users/:id/notification-info
func (h *UserHandler) GetUserNotificationInfo(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user ID is required"})
		return
	}

	ctx := c.Request.Context()
	notificationInfo, err := h.userService.GetUserNotificationInfo(ctx, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, notificationInfo)
}

// RegisterDevice handles POST /api/v1/users/:id/devices
func (h *UserHandler) RegisterDevice(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user ID is required"})
		return
	}

	var request struct {
		DeviceToken string `json:"device_token" binding:"required"`
		DeviceType  string `json:"device_type" binding:"required"`
		AppVersion  string `json:"app_version"`
		OSVersion   string `json:"os_version"`
		DeviceModel string `json:"device_model"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()
	device, err := h.userService.RegisterDevice(ctx, userID, request.DeviceToken, request.DeviceType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Update additional device info if provided
	if request.AppVersion != "" || request.OSVersion != "" || request.DeviceModel != "" {
		err = h.userService.UpdateDeviceInfo(ctx, device.ID, request.AppVersion, request.OSVersion, request.DeviceModel)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		// Get updated device
		devices, err := h.userService.GetUserDevices(ctx, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		for _, d := range devices {
			if d.ID == device.ID {
				device = d
				break
			}
		}
	}

	c.JSON(http.StatusCreated, device)
}

// GetUserDevices handles GET /api/v1/users/:id/devices
func (h *UserHandler) GetUserDevices(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user ID is required"})
		return
	}

	ctx := c.Request.Context()
	devices, err := h.userService.GetUserDevices(ctx, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id": userID,
		"devices": devices,
		"count":   len(devices),
	})
}

// GetActiveUserDevices handles GET /api/v1/users/:id/devices/active
func (h *UserHandler) GetActiveUserDevices(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user ID is required"})
		return
	}

	ctx := c.Request.Context()
	devices, err := h.userService.GetActiveUserDevices(ctx, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id": userID,
		"devices": devices,
		"count":   len(devices),
	})
}

// UpdateDeviceInfo handles PUT /api/v1/devices/:deviceId
func (h *UserHandler) UpdateDeviceInfo(c *gin.Context) {
	deviceID := c.Param("deviceId")
	if deviceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "device ID is required"})
		return
	}

	var request struct {
		AppVersion  string `json:"app_version"`
		OSVersion   string `json:"os_version"`
		DeviceModel string `json:"device_model"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()
	err := h.userService.UpdateDeviceInfo(ctx, deviceID, request.AppVersion, request.OSVersion, request.DeviceModel)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Device information updated successfully"})
}

// RemoveDevice handles DELETE /api/v1/devices/:deviceId
func (h *UserHandler) RemoveDevice(c *gin.Context) {
	deviceID := c.Param("deviceId")
	if deviceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "device ID is required"})
		return
	}

	ctx := c.Request.Context()
	err := h.userService.RemoveDevice(ctx, deviceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Device removed successfully"})
}

// DeactivateDevice handles PATCH /api/v1/devices/:deviceId/deactivate
func (h *UserHandler) DeactivateDevice(c *gin.Context) {
	deviceID := c.Param("deviceId")
	if deviceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "device ID is required"})
		return
	}

	ctx := c.Request.Context()
	err := h.userService.DeactivateDevice(ctx, deviceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Device deactivated successfully"})
}

// UpdateDeviceLastUsed handles PATCH /api/v1/devices/:deviceId/last-used
func (h *UserHandler) UpdateDeviceLastUsed(c *gin.Context) {
	deviceID := c.Param("deviceId")
	if deviceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "device ID is required"})
		return
	}

	ctx := c.Request.Context()
	err := h.userService.UpdateDeviceLastUsed(ctx, deviceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Device last used timestamp updated successfully"})
}
