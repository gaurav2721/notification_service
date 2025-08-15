package kafka

import (
	"os"
	"testing"
	"time"
)

func TestNewKafkaService(t *testing.T) {
	// Test with default values
	service, err := NewKafkaService()
	if err != nil {
		t.Fatalf("Failed to create KafkaService: %v", err)
	}
	defer service.Close()

	// Verify all channels are created
	if service.GetEmailChannel() == nil {
		t.Error("Email channel should not be nil")
	}
	if service.GetSlackChannel() == nil {
		t.Error("Slack channel should not be nil")
	}
	if service.GetIOSPushNotificationChannel() == nil {
		t.Error("iOS push notification channel should not be nil")
	}
	if service.GetAndroidPushNotificationChannel() == nil {
		t.Error("Android push notification channel should not be nil")
	}
}

func TestKafkaServiceWithCustomBufferSizes(t *testing.T) {
	// Set custom buffer sizes
	os.Setenv("EMAIL_CHANNEL_BUFFER_SIZE", "50")
	os.Setenv("SLACK_CHANNEL_BUFFER_SIZE", "75")
	os.Setenv("IOS_PUSH_CHANNEL_BUFFER_SIZE", "25")
	os.Setenv("ANDROID_PUSH_CHANNEL_BUFFER_SIZE", "30")

	service, err := NewKafkaService()
	if err != nil {
		t.Fatalf("Failed to create KafkaService: %v", err)
	}
	defer service.Close()

	// Test channel capacity
	if cap(service.GetEmailChannel()) != 50 {
		t.Errorf("Expected email channel capacity 50, got %d", cap(service.GetEmailChannel()))
	}
	if cap(service.GetSlackChannel()) != 75 {
		t.Errorf("Expected slack channel capacity 75, got %d", cap(service.GetSlackChannel()))
	}
	if cap(service.GetIOSPushNotificationChannel()) != 25 {
		t.Errorf("Expected iOS push channel capacity 25, got %d", cap(service.GetIOSPushNotificationChannel()))
	}
	if cap(service.GetAndroidPushNotificationChannel()) != 30 {
		t.Errorf("Expected Android push channel capacity 30, got %d", cap(service.GetAndroidPushNotificationChannel()))
	}
}

func TestKafkaServiceChannelOperations(t *testing.T) {
	service, err := NewKafkaService()
	if err != nil {
		t.Fatalf("Failed to create KafkaService: %v", err)
	}
	defer service.Close()

	// Test email channel
	emailChannel := service.GetEmailChannel()
	emailMsg := "test email message"

	// Send message
	select {
	case emailChannel <- emailMsg:
		// Message sent successfully
	default:
		t.Error("Failed to send message to email channel")
	}

	// Receive message
	select {
	case receivedMsg := <-emailChannel:
		if receivedMsg != emailMsg {
			t.Errorf("Expected message %s, got %s", emailMsg, receivedMsg)
		}
	case <-time.After(time.Second):
		t.Error("Timeout waiting for message from email channel")
	}

	// Test slack channel
	slackChannel := service.GetSlackChannel()
	slackMsg := "test slack message"

	select {
	case slackChannel <- slackMsg:
		// Message sent successfully
	default:
		t.Error("Failed to send message to slack channel")
	}

	select {
	case receivedMsg := <-slackChannel:
		if receivedMsg != slackMsg {
			t.Errorf("Expected message %s, got %s", slackMsg, receivedMsg)
		}
	case <-time.After(time.Second):
		t.Error("Timeout waiting for message from slack channel")
	}
}

func TestKafkaServiceClose(t *testing.T) {
	service, err := NewKafkaService()
	if err != nil {
		t.Fatalf("Failed to create KafkaService: %v", err)
	}

	// Close the service
	service.Close()

	// Verify channels are closed
	emailChannel := service.GetEmailChannel()
	_, ok := <-emailChannel
	if ok {
		t.Error("Email channel should be closed")
	}

	slackChannel := service.GetSlackChannel()
	_, ok = <-slackChannel
	if ok {
		t.Error("Slack channel should be closed")
	}

	iosChannel := service.GetIOSPushNotificationChannel()
	_, ok = <-iosChannel
	if ok {
		t.Error("iOS push notification channel should be closed")
	}

	androidChannel := service.GetAndroidPushNotificationChannel()
	_, ok = <-androidChannel
	if ok {
		t.Error("Android push notification channel should be closed")
	}
}

func TestGetEnvAsInt(t *testing.T) {
	// Test with valid environment variable
	os.Setenv("TEST_VAR", "42")
	result := getEnvAsInt("TEST_VAR", 100)
	if result != 42 {
		t.Errorf("Expected 42, got %d", result)
	}

	// Test with missing environment variable
	result = getEnvAsInt("MISSING_VAR", 100)
	if result != 100 {
		t.Errorf("Expected 100, got %d", result)
	}

	// Test with invalid environment variable
	os.Setenv("INVALID_VAR", "not_a_number")
	result = getEnvAsInt("INVALID_VAR", 100)
	if result != 100 {
		t.Errorf("Expected 100, got %d", result)
	}
}
