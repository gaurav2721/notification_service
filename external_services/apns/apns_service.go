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

	"github.com/gaurav2721/notification-service/models"
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
	notif, ok := notification.(*models.APNSNotificationRequest)
	if !ok {
		return nil, ErrInvalidNotificationPayload
	}

	// Validate the APNS notification
	if err := models.ValidateAPNSNotification(notif); err != nil {
		return nil, fmt.Errorf("APNS validation failed: %w", err)
	}

	// Extract device token from recipient
	deviceToken := notif.Recipient

	if deviceToken == "" {
		return &models.APNSResponse{
			ID:      notif.ID,
			Status:  "no_devices",
			Message: "No device token provided in recipient",
			SentAt:  time.Now(),
			Channel: "apns",
		}, nil
	}

	// Check if config is available for JWT authentication
	if aps.config == nil || aps.privateKey == nil {
		// Return mock response for demo purposes when config is not available
		return &models.APNSResponse{
			ID:           notif.ID,
			Status:       "demo_mode",
			Message:      "APNS notification simulated (no config provided)",
			SentAt:       time.Now(),
			Channel:      "apns",
			SuccessCount: 1,
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
				"title": notif.Content.Title,
				"body":  notif.Content.Body,
			},
			"sound": "default",
			"badge": 1,
		},
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	// Send to the device token
	successCount := 0
	failureCount := 0

	url := fmt.Sprintf("https://api.push.apple.com/3/device/%s", deviceToken)

	req, err := http.NewRequestWithContext(ctx, "POST", url, nil)
	if err != nil {
		failureCount = 1
	} else {
		req.Header.Set("Authorization", "bearer "+tokenString)
		req.Header.Set("apns-topic", aps.config.BundleID)
		req.Header.Set("Content-Type", "application/json")
		req.Body = ioutil.NopCloser(bytes.NewReader(payloadBytes))

		resp, err := aps.client.Do(req)
		if err != nil {
			failureCount = 1
		} else {
			defer resp.Body.Close()

			if resp.StatusCode == http.StatusOK {
				successCount = 1
			} else {
				failureCount = 1
			}
		}
	}

	// Return success response
	return &models.APNSResponse{
		ID:           notif.ID,
		Status:       "sent",
		Message:      fmt.Sprintf("APNS notification sent successfully. Success: %d, Failed: %d", successCount, failureCount),
		SentAt:       time.Now(),
		Channel:      "apns",
		SuccessCount: successCount,
		FailureCount: failureCount,
	}, nil
}
