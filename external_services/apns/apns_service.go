package apns

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// APNSServiceImpl implements the APNSService interface
type APNSServiceImpl struct {
	config     *APNSConfig
	privateKey *ecdsa.PrivateKey
	client     *http.Client
}

// NewAPNSService creates a new APNS service instance
// It checks environment variables and returns mock service if config is incomplete
func NewAPNSService() APNSService {
	bundleID := os.Getenv("APNS_BUNDLE_ID")
	keyID := os.Getenv("APNS_KEY_ID")
	teamID := os.Getenv("APNS_TEAM_ID")
	privateKeyPath := os.Getenv("APNS_PRIVATE_KEY_PATH")
	timeoutStr := os.Getenv("APNS_TIMEOUT")

	// Check if all required environment variables are present and non-empty
	if bundleID == "" || keyID == "" || teamID == "" || privateKeyPath == "" {
		return NewMockAPNSService()
	}

	timeout := 30 // default timeout
	if timeoutStr != "" {
		if t, err := strconv.Atoi(timeoutStr); err == nil && t > 0 {
			timeout = t
		}
	}

	return &APNSServiceImpl{
		client: &http.Client{
			Timeout: time.Duration(timeout) * time.Second,
		},
	}
}

// SendPushNotification sends a push notification to Apple devices
func (aps *APNSServiceImpl) SendPushNotification(ctx context.Context, notification interface{}) (interface{}, error) {
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
	// Recipients should contain iOS device tokens
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
			Channel: "apns",
		}, nil
	}

	// Check if config is available for JWT authentication
	if aps.config == nil || aps.privateKey == nil {
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
			Message:      "APNS notification simulated (no config provided)",
			SentAt:       time.Now(),
			Channel:      "apns",
			SuccessCount: len(deviceTokens),
			FailureCount: 0,
		}, nil
	}

	// Create JWT token for authentication
	token := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		"iss": aps.config.TeamID,
		"iat": time.Now().Unix(),
	})

	token.Header["kid"] = aps.config.KeyID
	tokenString, err := token.SignedString(aps.privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to sign JWT token: %w", err)
	}

	// Prepare notification payload
	payload := map[string]interface{}{
		"aps": map[string]interface{}{
			"alert": map[string]interface{}{
				"title": notif.Content["title"],
				"body":  notif.Content["body"],
			},
			"sound": "default",
			"badge": 1,
		},
	}

	if notif.Content["data"] != nil {
		payload["aps"].(map[string]interface{})["custom-data"] = notif.Content["data"]
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	// Send to each device token
	successCount := 0
	failureCount := 0

	for _, deviceToken := range deviceTokens {
		url := fmt.Sprintf("https://api.push.apple.com/3/device/%s", deviceToken)

		req, err := http.NewRequestWithContext(ctx, "POST", url, nil)
		if err != nil {
			failureCount++
			continue
		}

		req.Header.Set("Authorization", "bearer "+tokenString)
		req.Header.Set("apns-topic", aps.config.BundleID)
		req.Header.Set("Content-Type", "application/json")
		req.Body = ioutil.NopCloser(bytes.NewReader(payloadBytes))

		resp, err := aps.client.Do(req)
		if err != nil {
			failureCount++
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			successCount++
		} else {
			failureCount++
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
		Message:      fmt.Sprintf("APNS notification sent successfully. Success: %d, Failed: %d", successCount, failureCount),
		SentAt:       time.Now(),
		Channel:      "apns",
		SuccessCount: successCount,
		FailureCount: failureCount,
	}, nil
}
