package apns

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// APNSServiceImpl implements the APNSService interface
type APNSServiceImpl struct {
	config       *APNSConfig
	privateKey   *ecdsa.PrivateKey
	client       *http.Client
	deviceTokens map[string][]string // userID -> device tokens
	mutex        sync.RWMutex
}

// NewAPNSService creates a new APNS service instance
func NewAPNSService() APNSService {
	return &APNSServiceImpl{
		deviceTokens: make(map[string][]string),
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// NewAPNSServiceWithConfig creates a new APNS service with custom configuration
func NewAPNSServiceWithConfig(config *APNSConfig) (APNSService, error) {
	if config == nil {
		return nil, ErrInvalidConfiguration
	}

	// Load private key
	privateKeyBytes, err := ioutil.ReadFile(config.PrivateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read private key: %w", err)
	}

	privateKey, err := jwt.ParseECPrivateKeyFromPEM(privateKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	// Note: APNS endpoint is determined by environment
	// Sandbox: https://api.sandbox.push.apple.com
	// Production: https://api.push.apple.com

	return &APNSServiceImpl{
		config:       config,
		privateKey:   privateKey,
		deviceTokens: make(map[string][]string),
		client: &http.Client{
			Timeout: time.Duration(config.Timeout) * time.Second,
		},
	}, nil
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

	// Get device tokens for all recipients
	aps.mutex.RLock()
	var allDeviceTokens []string
	for _, recipient := range notif.Recipients {
		if tokens, exists := aps.deviceTokens[recipient]; exists {
			allDeviceTokens = append(allDeviceTokens, tokens...)
		}
	}
	aps.mutex.RUnlock()

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
			SuccessCount: len(allDeviceTokens),
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

	for _, deviceToken := range allDeviceTokens {
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

// RegisterDeviceToken registers a device token for a user
func (aps *APNSServiceImpl) RegisterDeviceToken(userID, deviceToken string) error {
	if userID == "" || deviceToken == "" {
		return ErrInvalidDeviceToken
	}

	aps.mutex.Lock()
	defer aps.mutex.Unlock()

	// Check if token already exists
	if tokens, exists := aps.deviceTokens[userID]; exists {
		for _, token := range tokens {
			if token == deviceToken {
				return nil // Token already registered
			}
		}
		aps.deviceTokens[userID] = append(tokens, deviceToken)
	} else {
		aps.deviceTokens[userID] = []string{deviceToken}
	}

	return nil
}

// UnregisterDeviceToken removes a device token for a user
func (aps *APNSServiceImpl) UnregisterDeviceToken(userID, deviceToken string) error {
	if userID == "" || deviceToken == "" {
		return ErrInvalidDeviceToken
	}

	aps.mutex.Lock()
	defer aps.mutex.Unlock()

	if tokens, exists := aps.deviceTokens[userID]; exists {
		for i, token := range tokens {
			if token == deviceToken {
				aps.deviceTokens[userID] = append(tokens[:i], tokens[i+1:]...)
				return nil
			}
		}
		return ErrDeviceTokenNotFound
	}

	return ErrUserNotFound
}

// GetDeviceTokensForUser retrieves all device tokens for a user
func (aps *APNSServiceImpl) GetDeviceTokensForUser(userID string) ([]string, error) {
	if userID == "" {
		return nil, ErrUserNotFound
	}

	aps.mutex.RLock()
	defer aps.mutex.RUnlock()

	if tokens, exists := aps.deviceTokens[userID]; exists {
		return tokens, nil
	}

	return []string{}, nil
}
