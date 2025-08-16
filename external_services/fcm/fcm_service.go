package fcm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

// FCMServiceImpl implements the FCMService interface
type FCMServiceImpl struct {
	config       *FCMConfig
	client       *http.Client
	deviceTokens map[string][]string // userID -> device tokens
	mutex        sync.RWMutex
}

// FCMRequest represents the FCM API request structure
type FCMRequest struct {
	To              string                 `json:"to,omitempty"`
	RegistrationIDs []string               `json:"registration_ids,omitempty"`
	Data            map[string]interface{} `json:"data,omitempty"`
	Notification    *FCMNotification       `json:"notification,omitempty"`
	Priority        string                 `json:"priority,omitempty"`
	TTL             int                    `json:"time_to_live,omitempty"`
}

// FCMNotification represents the notification payload
type FCMNotification struct {
	Title string `json:"title,omitempty"`
	Body  string `json:"body,omitempty"`
	Sound string `json:"sound,omitempty"`
	Badge string `json:"badge,omitempty"`
}

// FCMResponse represents the FCM API response structure
type FCMResponse struct {
	MulticastID  int64    `json:"multicast_id"`
	Success      int      `json:"success"`
	Failure      int      `json:"failure"`
	CanonicalIDs int      `json:"canonical_ids"`
	Results      []Result `json:"results"`
}

// Result represents individual result for each token
type Result struct {
	MessageID      string `json:"message_id,omitempty"`
	RegistrationID string `json:"registration_id,omitempty"`
	Error          string `json:"error,omitempty"`
}

// NewFCMService creates a new FCM service instance
func NewFCMService() FCMService {
	return &FCMServiceImpl{
		deviceTokens: make(map[string][]string),
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// NewFCMServiceWithConfig creates a new FCM service with custom configuration
func NewFCMServiceWithConfig(config *FCMConfig) (FCMService, error) {
	if config == nil || config.ServerKey == "" {
		return nil, ErrInvalidConfiguration
	}

	return &FCMServiceImpl{
		config:       config,
		deviceTokens: make(map[string][]string),
		client: &http.Client{
			Timeout: time.Duration(config.Timeout) * time.Second,
		},
	}, nil
}

// SendPushNotification sends a push notification to Android devices via FCM
func (fcm *FCMServiceImpl) SendPushNotification(ctx context.Context, notification interface{}) (interface{}, error) {
	// Type assertion to get the notification
	notif, ok := notification.(*struct {
		ID       string
		Type     string
		Content  map[string]interface{}
		Template *struct {
			ID   string
			Data map[string]interface{}
		}
		Recipients  []string
		ScheduledAt *time.Time
	})
	if !ok {
		return nil, ErrInvalidNotificationPayload
	}

	// Get device tokens for all recipients
	fcm.mutex.RLock()
	var allDeviceTokens []string
	for _, recipient := range notif.Recipients {
		if tokens, exists := fcm.deviceTokens[recipient]; exists {
			allDeviceTokens = append(allDeviceTokens, tokens...)
		}
	}
	fcm.mutex.RUnlock()

	if len(allDeviceTokens) == 0 {
		return &struct {
			ID      string    `json:"id"`
			Status  string    `json:"status"`
			Message string    `json:"message"`
			SentAt  time.Time `json:"sent_at"`
			Channel string    `json:"channel"`
		}{
			ID:      notif.ID,
			Status:  "no_devices",
			Message: "No device tokens found for recipients",
			SentAt:  time.Now(),
			Channel: "fcm",
		}, nil
	}

	// Prepare notification payload
	fcmNotification := &FCMNotification{
		Title: fmt.Sprintf("%v", notif.Content["title"]),
		Body:  fmt.Sprintf("%v", notif.Content["body"]),
		Sound: "default",
	}

	// Prepare data payload
	data := make(map[string]interface{})
	if notif.Content["data"] != nil {
		if dataMap, ok := notif.Content["data"].(map[string]interface{}); ok {
			data = dataMap
		}
	}
	data["notification_id"] = notif.ID
	data["type"] = notif.Type

	// Check if config is available
	if fcm.config == nil {
		// Return mock response for demo purposes when config is not available
		return &struct {
			ID           string    `json:"id"`
			Status       string    `json:"status"`
			Message      string    `json:"message"`
			SentAt       time.Time `json:"sent_at"`
			Channel      string    `json:"channel"`
			SuccessCount int       `json:"success_count"`
			FailureCount int       `json:"failure_count"`
		}{
			ID:           notif.ID,
			Status:       "demo_mode",
			Message:      "FCM notification simulated (no config provided)",
			SentAt:       time.Now(),
			Channel:      "fcm",
			SuccessCount: len(allDeviceTokens),
			FailureCount: 0,
		}, nil
	}

	// Send notifications in batches
	batchSize := fcm.config.BatchSize
	if batchSize <= 0 {
		batchSize = 1000 // Default batch size
	}

	totalSuccess := 0
	totalFailure := 0

	for i := 0; i < len(allDeviceTokens); i += batchSize {
		end := i + batchSize
		if end > len(allDeviceTokens) {
			end = len(allDeviceTokens)
		}

		batchTokens := allDeviceTokens[i:end]
		success, failure, err := fcm.sendBatch(ctx, batchTokens, fcmNotification, data)
		if err != nil {
			failure += len(batchTokens)
		} else {
			totalSuccess += success
			totalFailure += failure
		}
	}

	// Return success response
	return &struct {
		ID           string    `json:"id"`
		Status       string    `json:"status"`
		Message      string    `json:"message"`
		SentAt       time.Time `json:"sent_at"`
		Channel      string    `json:"channel"`
		SuccessCount int       `json:"success_count"`
		FailureCount int       `json:"failure_count"`
	}{
		ID:           notif.ID,
		Status:       "sent",
		Message:      fmt.Sprintf("FCM notification sent successfully. Success: %d, Failed: %d", totalSuccess, totalFailure),
		SentAt:       time.Now(),
		Channel:      "fcm",
		SuccessCount: totalSuccess,
		FailureCount: totalFailure,
	}, nil
}

// sendBatch sends a batch of notifications to FCM
func (fcm *FCMServiceImpl) sendBatch(ctx context.Context, tokens []string, notification *FCMNotification, data map[string]interface{}) (success, failure int, err error) {
	request := &FCMRequest{
		RegistrationIDs: tokens,
		Notification:    notification,
		Data:            data,
		Priority:        "high",
		TTL:             86400, // 24 hours
	}

	requestBytes, err := json.Marshal(request)
	if err != nil {
		return 0, len(tokens), fmt.Errorf("failed to marshal FCM request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://fcm.googleapis.com/fcm/send", bytes.NewReader(requestBytes))
	if err != nil {
		return 0, len(tokens), fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Authorization", "key="+fcm.config.ServerKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := fcm.client.Do(req)
	if err != nil {
		return 0, len(tokens), fmt.Errorf("failed to send FCM request: %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, len(tokens), fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return 0, len(tokens), fmt.Errorf("FCM API error: %s", string(body))
	}

	var fcmResp FCMResponse
	if err := json.Unmarshal(body, &fcmResp); err != nil {
		return 0, len(tokens), fmt.Errorf("failed to unmarshal FCM response: %w", err)
	}

	return fcmResp.Success, fcmResp.Failure, nil
}

// RegisterDeviceToken registers a device token for a user
func (fcm *FCMServiceImpl) RegisterDeviceToken(userID, deviceToken string) error {
	if userID == "" || deviceToken == "" {
		return ErrInvalidDeviceToken
	}

	fcm.mutex.Lock()
	defer fcm.mutex.Unlock()

	// Check if token already exists
	if tokens, exists := fcm.deviceTokens[userID]; exists {
		for _, token := range tokens {
			if token == deviceToken {
				return nil // Token already registered
			}
		}
		fcm.deviceTokens[userID] = append(tokens, deviceToken)
	} else {
		fcm.deviceTokens[userID] = []string{deviceToken}
	}

	return nil
}

// UnregisterDeviceToken removes a device token for a user
func (fcm *FCMServiceImpl) UnregisterDeviceToken(userID, deviceToken string) error {
	if userID == "" || deviceToken == "" {
		return ErrInvalidDeviceToken
	}

	fcm.mutex.Lock()
	defer fcm.mutex.Unlock()

	if tokens, exists := fcm.deviceTokens[userID]; exists {
		for i, token := range tokens {
			if token == deviceToken {
				fcm.deviceTokens[userID] = append(tokens[:i], tokens[i+1:]...)
				return nil
			}
		}
		return ErrDeviceTokenNotFound
	}

	return ErrUserNotFound
}

// GetDeviceTokensForUser retrieves all device tokens for a user
func (fcm *FCMServiceImpl) GetDeviceTokensForUser(userID string) ([]string, error) {
	if userID == "" {
		return nil, ErrUserNotFound
	}

	fcm.mutex.RLock()
	defer fcm.mutex.RUnlock()

	if tokens, exists := fcm.deviceTokens[userID]; exists {
		return tokens, nil
	}

	return []string{}, nil
}
