#!/bin/bash

# Test script for notification API
echo "Starting notification service test..."

# Start the service in background
echo "Starting notification service..."
./bin/notification-service &
SERVICE_PID=$!

# Wait for service to start
sleep 3

# Test 1: Send email notification
echo "Test 1: Sending email notification..."
curl -X POST http://localhost:8080/api/v1/notifications \
  -H "Content-Type: application/json" \
  -d '{
    "type": "email",
    "content": {
      "subject": "Test Email",
      "email_body": "This is a test email notification"
    },
    "recipients": ["user-001", "user-002"]
  }'

echo -e "\n\n"

# Test 2: Send Slack notification
echo "Test 2: Sending Slack notification..."
curl -X POST http://localhost:8080/api/v1/notifications \
  -H "Content-Type: application/json" \
  -d '{
    "type": "slack",
    "content": {
      "text": "This is a test Slack notification"
    },
    "recipients": ["user-001"]
  }'

echo -e "\n\n"

# Test 3: Send iOS push notification
echo "Test 3: Sending iOS push notification..."
curl -X POST http://localhost:8080/api/v1/notifications \
  -H "Content-Type: application/json" \
  -d '{
    "type": "ios_push",
    "content": {
      "title": "Test Push",
      "body": "This is a test iOS push notification"
    },
    "recipients": ["user-003"]
  }'

echo -e "\n\n"

# Test 4: Send scheduled notification
echo "Test 4: Sending scheduled notification..."
curl -X POST http://localhost:8080/api/v1/notifications \
  -H "Content-Type: application/json" \
  -d '{
    "type": "email",
    "content": {
      "subject": "Scheduled Email",
      "email_body": "This is a scheduled email notification"
    },
    "recipients": ["user-001"],
    "scheduled_at": "2024-12-31T23:59:59Z"
  }'

echo -e "\n\n"

# Test 5: Get notification status
echo "Test 5: Getting notification status..."
curl -X GET http://localhost:8080/api/v1/notifications/123456789

echo -e "\n\n"

# Stop the service
echo "Stopping notification service..."
kill $SERVICE_PID

echo "Test completed!" 