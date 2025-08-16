package fcm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

// FCMServiceImpl implements the FCMService interface
type FCMServiceImpl struct {
	config *FCMConfig
	client *http.Client
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
// It checks environment variables and returns mock service if config is incomplete
func NewFCMService() FCMService {
	serverKey := os.Getenv("FCM_SERVER_KEY")
	timeoutStr := os.Getenv("FCM_TIMEOUT")
	batchSizeStr := os.Getenv("FCM_BATCH_SIZE")

	// Check if all required environment variables are present and non-empty
	if serverKey == "" {
		return NewMockFCMService()
	}

	timeout := 30 // default timeout
	if timeoutStr != "" {
		if t, err := strconv.Atoi(timeoutStr); err == nil && t > 0 {
			timeout = t
		}
	}

	batchSize := 1000 // default batch size
	if batchSizeStr != "" {
		if b, err := strconv.Atoi(batchSizeStr); err == nil && b > 0 {
			batchSize = b
		}
	}

	return &FCMServiceImpl{
		config: &FCMConfig{
			ServerKey: serverKey,
			Timeout:   timeout,
			BatchSize: batchSize,
		},
		client: &http.Client{
			Timeout: time.Duration(timeout) * time.Second,
		},
	}
}

// NewFCMServiceWithConfig creates a new FCM service with custom configuration
func NewFCMServiceWithConfig(config *FCMConfig) (FCMService, error) {
	if config == nil || config.ServerKey == "" {
		return nil, ErrInvalidConfiguration
	}

	return &FCMServiceImpl{
		config: config,
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

	// Extract device tokens from recipients
	// Recipients should contain Android device tokens
	deviceTokens := notif.Recipients

	if len(deviceTokens) == 0 {
		return &struct {
			ID      string    `json:"id"`
			Status  string    `json:"status"`
			Message string    `json:"message"`
			SentAt  time.Time `json:"sent_at"`
			Channel string    `json:"channel"`
		}{
			ID:      notif.ID,
			Status:  "no_devices",
			Message: "No device tokens provided in recipients",
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
			SuccessCount: len(deviceTokens),
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

	for i := 0; i < len(deviceTokens); i += batchSize {
		end := i + batchSize
		if end > len(deviceTokens) {
			end = len(deviceTokens)
		}

		batchTokens := deviceTokens[i:end]
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
