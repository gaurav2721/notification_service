package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gaurav2721/notification-service/handlers"
	"github.com/gaurav2721/notification-service/routes"
	"github.com/gaurav2721/notification-service/services"
	"github.com/gin-gonic/gin"
)

// TestAPIKeyAuthentication demonstrates how to test the API key middleware
func TestAPIKeyAuthentication(t *testing.T) {
	// Set up test API key
	os.Setenv("API_KEY", "test-api-key-123")
	defer os.Unsetenv("API_KEY")

	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Initialize service container and handlers
	serviceContainer := services.NewServiceContainer()
	notificationHandler := handlers.NewNotificationHandler(
		serviceContainer.GetNotificationService(),
		serviceContainer.GetUserService(),
		serviceContainer.GetKafkaService(),
	)
	userHandler := handlers.NewUserHandler(serviceContainer.GetUserService())

	// Setup router
	router := gin.New()
	routes.SetupRoutes(router, notificationHandler, userHandler)

	// Test cases
	testCases := []struct {
		name           string
		authHeader     string
		expectedStatus int
		description    string
	}{
		{
			name:           "Valid API Key - Bearer format",
			authHeader:     "Bearer test-api-key-123",
			expectedStatus: http.StatusOK, // or whatever the actual endpoint returns
			description:    "Should allow access with valid API key in Bearer format",
		},
		{
			name:           "Valid API Key - ApiKey format",
			authHeader:     "ApiKey test-api-key-123",
			expectedStatus: http.StatusOK,
			description:    "Should allow access with valid API key in ApiKey format",
		},
		{
			name:           "Valid API Key - Direct format",
			authHeader:     "test-api-key-123",
			expectedStatus: http.StatusOK,
			description:    "Should allow access with valid API key in direct format",
		},
		{
			name:           "Invalid API Key",
			authHeader:     "Bearer wrong-api-key",
			expectedStatus: http.StatusUnauthorized,
			description:    "Should reject invalid API key",
		},
		{
			name:           "Missing API Key",
			authHeader:     "",
			expectedStatus: http.StatusUnauthorized,
			description:    "Should reject request without API key",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create test request
			req, _ := http.NewRequest("GET", "/api/v1/notifications/test-id", nil)
			if tc.authHeader != "" {
				req.Header.Set("Authorization", tc.authHeader)
			}

			// Create response recorder
			w := httptest.NewRecorder()

			// Serve the request
			router.ServeHTTP(w, req)

			// Check status code
			if w.Code != tc.expectedStatus {
				t.Errorf("Expected status %d, got %d: %s", tc.expectedStatus, w.Code, tc.description)
			}

			// For unauthorized requests, check error message
			if tc.expectedStatus == http.StatusUnauthorized {
				var response map[string]interface{}
				if err := json.Unmarshal(w.Body.Bytes(), &response); err == nil {
					if _, exists := response["error"]; !exists {
						t.Errorf("Expected error field in response for unauthorized request")
					}
				}
			}
		})
	}
}

// Example of how to make a real API call with API key
func ExampleAPICallWithKey() {
	// Set your API key
	apiKey := "your-secure-api-key-here"

	// Create request payload
	payload := map[string]interface{}{
		"type":      "email",
		"recipient": "user@example.com",
		"subject":   "Test Notification",
		"body":      "This is a test notification",
	}

	jsonPayload, _ := json.Marshal(payload)

	// Create HTTP request
	req, _ := http.NewRequest("POST", "http://localhost:8080/api/v1/notifications", bytes.NewBuffer(jsonPayload))
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error making request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	fmt.Printf("Response status: %s\n", resp.Status)
}

func main() {
	fmt.Println("API Key Authentication Test Examples")
	fmt.Println("Run 'go test' to execute the tests")
}
