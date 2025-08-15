# User Management API Documentation

This document provides comprehensive documentation for the User Management API endpoints, including user operations and device management functionality.

## Table of Contents

- [Base URL](#base-url)
- [Authentication](#authentication)
- [User Management Endpoints](#user-management-endpoints)
- [Device Management Endpoints](#device-management-endpoints)
- [Error Responses](#error-responses)
- [Data Models](#data-models)

## Base URL

```
http://localhost:8080/api/v1
```

## Authentication

Currently, the API does not require authentication. All endpoints are publicly accessible.

## User Management Endpoints

### 1. Get All Users

Retrieves all active users in the system.

**Endpoint:** `GET /users`

**Response:**
```json
{
  "users": [
    {
      "id": "user-001",
      "email": "john.doe@company.com",
      "full_name": "John Doe",
      "slack_user_id": "U1234567890",
      "slack_channel": "#general",
      "phone_number": "+1-555-0101",
      "is_active": true,
      "created_at": "2024-01-15T10:30:00Z",
      "updated_at": "2024-01-15T10:30:00Z"
    }
  ],
  "count": 1
}
```

**Curl Example:**
```bash
curl -X GET http://localhost:8080/api/v1/users/ \
  -H "Content-Type: application/json"
```

### 2. Get User by ID

Retrieves a specific user by their ID.

**Endpoint:** `GET /users/{id}`

**Parameters:**
- `id` (path parameter): User ID

**Response:**
```json
{
  "id": "user-001",
  "email": "john.doe@company.com",
  "full_name": "John Doe",
  "slack_user_id": "U1234567890",
  "slack_channel": "#general",
  "phone_number": "+1-555-0101",
  "is_active": true,
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:30:00Z"
}
```

**Curl Example:**
```bash
curl -X GET http://localhost:8080/api/v1/users/user-001 \
  -H "Content-Type: application/json"
```

### 3. Create User

Creates a new user in the system.

**Endpoint:** `POST /users`

**Request Body:**
```json
{
  "email": "jane.smith@company.com",
  "full_name": "Jane Smith"
}
```

**Required Fields:**
- `email`: Valid email address
- `full_name`: User's full name

**Response:**
```json
{
  "id": "user-002",
  "email": "jane.smith@company.com",
  "full_name": "Jane Smith",
  "slack_user_id": "",
  "slack_channel": "",
  "phone_number": "",
  "is_active": true,
  "created_at": "2024-01-15T11:00:00Z",
  "updated_at": "2024-01-15T11:00:00Z"
}
```

**Curl Example:**
```bash
curl -X POST http://localhost:8080/api/v1/users/ \
  -H "Content-Type: application/json" \
  -d '{
    "email": "jane.smith@company.com",
    "full_name": "Jane Smith"
  }'
```

### 4. Update User

Updates an existing user's information.

**Endpoint:** `PUT /users/{id}`

**Parameters:**
- `id` (path parameter): User ID

**Request Body:**
```json
{
  "email": "jane.smith.updated@company.com",
  "full_name": "Jane Smith Updated",
  "slack_user_id": "U0987654321",
  "slack_channel": "#design",
  "phone_number": "+1-555-0102"
}
```

**Optional Fields:**
- `email`: New email address
- `full_name`: New full name
- `slack_user_id`: Slack user ID
- `slack_channel`: Slack channel
- `phone_number`: Phone number

**Response:**
```json
{
  "id": "user-002",
  "email": "jane.smith.updated@company.com",
  "full_name": "Jane Smith Updated",
  "slack_user_id": "U0987654321",
  "slack_channel": "#design",
  "phone_number": "+1-555-0102",
  "is_active": true,
  "created_at": "2024-01-15T11:00:00Z",
  "updated_at": "2024-01-15T11:30:00Z"
}
```

**Curl Example:**
```bash
curl -X PUT http://localhost:8080/api/v1/users/user-002 \
  -H "Content-Type: application/json" \
  -d '{
    "email": "jane.smith.updated@company.com",
    "full_name": "Jane Smith Updated",
    "slack_user_id": "U0987654321",
    "slack_channel": "#design",
    "phone_number": "+1-555-0102"
  }'
```

### 5. Delete User

Soft deletes a user (marks as inactive).

**Endpoint:** `DELETE /users/{id}`

**Parameters:**
- `id` (path parameter): User ID

**Response:**
```json
{
  "message": "User deleted successfully"
}
```

**Curl Example:**
```bash
curl -X DELETE http://localhost:8080/api/v1/users/user-002 \
  -H "Content-Type: application/json"
```

### 6. Get User Notification Info

Retrieves user information along with their active devices for notification purposes.

**Endpoint:** `GET /users/{id}/notification-info`

**Parameters:**
- `id` (path parameter): User ID

**Response:**
```json
{
  "id": "user-001",
  "email": "john.doe@company.com",
  "full_name": "John Doe",
  "slack_user_id": "U1234567890",
  "slack_channel": "#general",
  "phone_number": "+1-555-0101",
  "devices": [
    {
      "id": "device-001",
      "user_id": "user-001",
      "device_token": "ios_token_123456789",
      "device_type": "ios",
      "app_version": "1.2.3",
      "os_version": "iOS 16.0",
      "device_model": "iPhone 14",
      "is_active": true,
      "last_used_at": "2024-01-15T10:30:00Z",
      "created_at": "2024-01-15T10:00:00Z",
      "updated_at": "2024-01-15T10:30:00Z"
    }
  ]
}
```

**Curl Example:**
```bash
curl -X GET http://localhost:8080/api/v1/users/user-001/notification-info \
  -H "Content-Type: application/json"
```

## Device Management Endpoints

### 1. Register Device

Registers a new device for a user.

**Endpoint:** `POST /users/{id}/devices`

**Parameters:**
- `id` (path parameter): User ID

**Request Body:**
```json
{
  "device_token": "fcm_token_abc123",
  "device_type": "android",
  "app_version": "1.2.3",
  "os_version": "Android 13",
  "device_model": "Samsung Galaxy S23"
}
```

**Required Fields:**
- `device_token`: Device token for push notifications
- `device_type`: Device type ("ios", "android", "web")

**Optional Fields:**
- `app_version`: Application version
- `os_version`: Operating system version
- `device_model`: Device model

**Response:**
```json
{
  "id": "device-003",
  "user_id": "user-001",
  "device_token": "fcm_token_abc123",
  "device_type": "android",
  "app_version": "1.2.3",
  "os_version": "Android 13",
  "device_model": "Samsung Galaxy S23",
  "is_active": true,
  "last_used_at": "2024-01-15T12:00:00Z",
  "created_at": "2024-01-15T12:00:00Z",
  "updated_at": "2024-01-15T12:00:00Z"
}
```

**Curl Example:**
```bash
curl -X POST http://localhost:8080/api/v1/users/user-001/devices \
  -H "Content-Type: application/json" \
  -d '{
    "device_token": "fcm_token_abc123",
    "device_type": "android",
    "app_version": "1.2.3",
    "os_version": "Android 13",
    "device_model": "Samsung Galaxy S23"
  }'
```

### 2. Get User Devices

Retrieves all devices (active and inactive) for a user.

**Endpoint:** `GET /users/{id}/devices`

**Parameters:**
- `id` (path parameter): User ID

**Response:**
```json
{
  "user_id": "user-001",
  "devices": [
    {
      "id": "device-001",
      "user_id": "user-001",
      "device_token": "ios_token_123456789",
      "device_type": "ios",
      "app_version": "1.2.3",
      "os_version": "iOS 16.0",
      "device_model": "iPhone 14",
      "is_active": true,
      "last_used_at": "2024-01-15T10:30:00Z",
      "created_at": "2024-01-15T10:00:00Z",
      "updated_at": "2024-01-15T10:30:00Z"
    }
  ],
  "count": 1
}
```

**Curl Example:**
```bash
curl -X GET http://localhost:8080/api/v1/users/user-001/devices \
  -H "Content-Type: application/json"
```

### 3. Get Active User Devices

Retrieves only active devices for a user.

**Endpoint:** `GET /users/{id}/devices/active`

**Parameters:**
- `id` (path parameter): User ID

**Response:**
```json
{
  "user_id": "user-001",
  "devices": [
    {
      "id": "device-001",
      "user_id": "user-001",
      "device_token": "ios_token_123456789",
      "device_type": "ios",
      "app_version": "1.2.3",
      "os_version": "iOS 16.0",
      "device_model": "iPhone 14",
      "is_active": true,
      "last_used_at": "2024-01-15T10:30:00Z",
      "created_at": "2024-01-15T10:00:00Z",
      "updated_at": "2024-01-15T10:30:00Z"
    }
  ],
  "count": 1
}
```

**Curl Example:**
```bash
curl -X GET http://localhost:8080/api/v1/users/user-001/devices/active \
  -H "Content-Type: application/json"
```

### 4. Update Device Information

Updates device information (app version, OS version, device model).

**Endpoint:** `PUT /devices/{deviceId}`

**Parameters:**
- `deviceId` (path parameter): Device ID

**Request Body:**
```json
{
  "app_version": "1.3.0",
  "os_version": "iOS 17.0",
  "device_model": "iPhone 15"
}
```

**Optional Fields:**
- `app_version`: New application version
- `os_version`: New operating system version
- `device_model`: New device model

**Response:**
```json
{
  "message": "Device information updated successfully"
}
```

**Curl Example:**
```bash
curl -X PUT http://localhost:8080/api/v1/devices/device-001 \
  -H "Content-Type: application/json" \
  -d '{
    "app_version": "1.3.0",
    "os_version": "iOS 17.0",
    "device_model": "iPhone 15"
  }'
```

### 5. Remove Device

Completely removes a device from the system.

**Endpoint:** `DELETE /devices/{deviceId}`

**Parameters:**
- `deviceId` (path parameter): Device ID

**Response:**
```json
{
  "message": "Device removed successfully"
}
```

**Curl Example:**
```bash
curl -X DELETE http://localhost:8080/api/v1/devices/device-001 \
  -H "Content-Type: application/json"
```

### 6. Deactivate Device

Marks a device as inactive (soft delete).

**Endpoint:** `PATCH /devices/{deviceId}/deactivate`

**Parameters:**
- `deviceId` (path parameter): Device ID

**Response:**
```json
{
  "message": "Device deactivated successfully"
}
```

**Curl Example:**
```bash
curl -X PATCH http://localhost:8080/api/v1/devices/device-001/deactivate \
  -H "Content-Type: application/json"
```

### 7. Update Device Last Used

Updates the last used timestamp for a device.

**Endpoint:** `PATCH /devices/{deviceId}/last-used`

**Parameters:**
- `deviceId` (path parameter): Device ID

**Response:**
```json
{
  "message": "Device last used timestamp updated successfully"
}
```

**Curl Example:**
```bash
curl -X PATCH http://localhost:8080/api/v1/devices/device-001/last-used \
  -H "Content-Type: application/json"
```

## Error Responses

The API returns standard HTTP status codes and error messages in JSON format.

### Common Error Responses

**400 Bad Request:**
```json
{
  "error": "user ID is required"
}
```

**404 Not Found:**
```json
{
  "error": "user not found"
}
```

**500 Internal Server Error:**
```json
{
  "error": "failed to create user"
}
```

### HTTP Status Codes

- `200 OK`: Request successful
- `201 Created`: Resource created successfully
- `400 Bad Request`: Invalid request data
- `404 Not Found`: Resource not found
- `500 Internal Server Error`: Server error

## Data Models

### User Model

```json
{
  "id": "string",
  "email": "string",
  "full_name": "string",
  "slack_user_id": "string",
  "slack_channel": "string",
  "phone_number": "string",
  "is_active": "boolean",
  "created_at": "datetime",
  "updated_at": "datetime"
}
```

### UserDeviceInfo Model

```json
{
  "id": "string",
  "user_id": "string",
  "device_token": "string",
  "device_type": "string",
  "app_version": "string",
  "os_version": "string",
  "device_model": "string",
  "is_active": "boolean",
  "last_used_at": "datetime",
  "created_at": "datetime",
  "updated_at": "datetime"
}
```

### UserNotificationInfo Model

```json
{
  "id": "string",
  "email": "string",
  "full_name": "string",
  "slack_user_id": "string",
  "slack_channel": "string",
  "phone_number": "string",
  "devices": [
    {
      "id": "string",
      "user_id": "string",
      "device_token": "string",
      "device_type": "string",
      "app_version": "string",
      "os_version": "string",
      "device_model": "string",
      "is_active": "boolean",
      "last_used_at": "datetime",
      "created_at": "datetime",
      "updated_at": "datetime"
    }
  ]
}
```

## Testing the API

### Prerequisites

1. Ensure the notification service is running on `http://localhost:8080`
2. Make sure you have `curl` installed on your system
3. Optional: Install `jq` for pretty-printed JSON output (not required for basic testing)

### Quick Test Script

You can use the following bash script to test all endpoints:

```bash
#!/bin/bash

BASE_URL="http://localhost:8080/api/v1"

echo "Testing User Management API..."

# Create a user
echo "1. Creating a user..."
USER_RESPONSE=$(curl -s -X POST "$BASE_URL/users/" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test.user@company.com",
    "full_name": "Test User"
  }')

USER_ID=$(echo $USER_RESPONSE | grep -o '"id":"[^"]*"' | cut -d'"' -f4)
echo "Created user with ID: $USER_ID"

# Get all users
echo "2. Getting all users..."
curl -s -X GET "$BASE_URL/users/" | jq '.'

# Get specific user
echo "3. Getting specific user..."
curl -s -X GET "$BASE_URL/users/$USER_ID" | jq '.'

# Register a device
echo "4. Registering a device..."
DEVICE_RESPONSE=$(curl -s -X POST "$BASE_URL/users/$USER_ID/devices" \
  -H "Content-Type: application/json" \
  -d '{
    "device_token": "test_device_token_123",
    "device_type": "ios",
    "app_version": "1.0.0",
    "os_version": "iOS 16.0",
    "device_model": "iPhone 14"
  }')

DEVICE_ID=$(echo $DEVICE_RESPONSE | grep -o '"id":"[^"]*"' | cut -d'"' -f4)
echo "Registered device with ID: $DEVICE_ID"

# Get user devices
echo "5. Getting user devices..."
curl -s -X GET "$BASE_URL/users/$USER_ID/devices" | jq '.'

# Get user notification info
echo "6. Getting user notification info..."
curl -s -X GET "$BASE_URL/users/$USER_ID/notification-info" | jq '.'

# Update device info
echo "7. Updating device info..."
curl -s -X PUT "$BASE_URL/devices/$DEVICE_ID" \
  -H "Content-Type: application/json" \
  -d '{
    "app_version": "1.1.0",
    "os_version": "iOS 17.0"
  }' | jq '.'

# Update user
echo "8. Updating user..."
curl -s -X PUT "$BASE_URL/users/$USER_ID" \
  -H "Content-Type: application/json" \
  -d '{
    "full_name": "Test User Updated",
    "slack_user_id": "U1234567890"
  }' | jq '.'

# Deactivate device
echo "9. Deactivating device..."
curl -s -X PATCH "$BASE_URL/devices/$DEVICE_ID/deactivate" | jq '.'

# Remove device
echo "10. Removing device..."
curl -s -X DELETE "$BASE_URL/devices/$DEVICE_ID" | jq '.'

# Delete user
echo "11. Deleting user..."
curl -s -X DELETE "$BASE_URL/users/$USER_ID" | jq '.'

echo "API testing completed!"
```

### Individual Test Commands

If you prefer to test endpoints individually, here are some key commands:

```bash
# Test server is running
curl -X GET http://localhost:8080/api/v1/users/

# Create a test user
curl -X POST http://localhost:8080/api/v1/users/ \
  -H "Content-Type: application/json" \
  -d '{"email": "test@example.com", "full_name": "Test User"}'

# Get user notification info (replace user-001 with actual user ID)
curl -X GET http://localhost:8080/api/v1/users/user-001/notification-info

# Register a device (replace user-001 with actual user ID)
curl -X POST http://localhost:8080/api/v1/users/user-001/devices \
  -H "Content-Type: application/json" \
  -d '{"device_token": "test_token", "device_type": "android"}'
```

## Notes

- All timestamps are in ISO 8601 format (UTC)
- Device types supported: "ios", "android", "web"
- User deletion is soft delete (marks as inactive)
- Device tokens should be unique per device
- The API includes sample data for testing purposes
- **URL Patterns**: Some endpoints redirect automatically:
  - `/users` redirects to `/users/` (use trailing slash)
  - `/users/{id}/devices` works without trailing slash
  - `/users/{id}/notification-info` works without trailing slash
  - Device management endpoints work without trailing slashes 