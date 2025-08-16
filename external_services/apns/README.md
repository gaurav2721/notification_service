# Apple Push Notification Service (APNS)

This package provides a Go implementation for sending push notifications to Apple devices using the Apple Push Notification Service (APNS).

## Features

- Send push notifications to iOS, macOS, and tvOS devices
- Device token management (register/unregister)
- JWT-based authentication with Apple's APNS API
- Support for both sandbox and production environments
- Thread-safe operations
- Comprehensive error handling

## Configuration

The APNS service requires the following configuration:

```go
config := &APNSConfig{
    BundleID:       "com.yourcompany.yourapp",
    KeyID:          "YOUR_KEY_ID",
    TeamID:         "YOUR_TEAM_ID",
    PrivateKeyPath: "/path/to/your/private/key.p8",
    Environment:    "sandbox", // or "production"
    MaxRetries:     3,
    Timeout:        30, // seconds
}
```

### Required Credentials

1. **Bundle ID**: Your app's bundle identifier
2. **Key ID**: The key identifier from your APNS key
3. **Team ID**: Your Apple Developer Team ID
4. **Private Key**: The .p8 file containing your APNS authentication key

## Usage

### Basic Usage

```go
package main

import (
    "context"
    "log"
    
    "github.com/gaurav2721/notification-service/external_services/apns"
)

func main() {
    // Create APNS service
    service := apns.NewAPNSService()
    
    // Register device token for a user
    err := service.RegisterDeviceToken("user123", "device_token_here")
    if err != nil {
        log.Fatal(err)
    }
    
    // Send push notification
    notification := &struct {
        ID       string
        Type     string
        Content  map[string]interface{}
        Template *struct {
            ID   string
            Data map[string]interface{}
        }
        Recipients  []string
        ScheduledAt *time.Time
    }{
        ID:   "notification_123",
        Type: "alert",
        Content: map[string]interface{}{
            "title": "Hello!",
            "body":  "This is a test notification",
        },
        Recipients: []string{"user123"},
    }
    
    result, err := service.SendPushNotification(context.Background(), notification)
    if err != nil {
        log.Fatal(err)
    }
    
    log.Printf("Notification sent: %+v", result)
}
```

### Advanced Usage with Configuration

```go
config := &apns.APNSConfig{
    BundleID:       "com.example.app",
    KeyID:          "ABC123DEF4",
    TeamID:         "TEAM123456",
    PrivateKeyPath: "/path/to/AuthKey_ABC123DEF4.p8",
    Environment:    "production",
    MaxRetries:     3,
    Timeout:        30,
}

service, err := apns.NewAPNSServiceWithConfig(config)
if err != nil {
    log.Fatal(err)
}
```

## API Reference

### Methods

#### `SendPushNotification(ctx context.Context, notification interface{}) (interface{}, error)`
Sends a push notification to all registered devices of the specified recipients.

#### `RegisterDeviceToken(userID, deviceToken string) error`
Registers a device token for a user. Multiple tokens can be registered per user.

#### `UnregisterDeviceToken(userID, deviceToken string) error`
Removes a device token for a user.

#### `GetDeviceTokensForUser(userID string) ([]string, error)`
Retrieves all device tokens registered for a user.

## Error Handling

The service provides specific error types:

- `ErrAPNSSendFailed`: Failed to send push notification
- `ErrInvalidDeviceToken`: Invalid device token provided
- `ErrInvalidConfiguration`: Invalid service configuration
- `ErrDeviceTokenNotFound`: Device token not found
- `ErrUserNotFound`: User not found
- `ErrInvalidNotificationPayload`: Invalid notification payload

## Testing

Run the tests:

```bash
go test ./external_services/apns
```

## Environment Variables

For production deployment, consider using environment variables:

```bash
export APNS_BUNDLE_ID="com.yourcompany.yourapp"
export APNS_KEY_ID="YOUR_KEY_ID"
export APNS_TEAM_ID="YOUR_TEAM_ID"
export APNS_PRIVATE_KEY_PATH="/path/to/private/key.p8"
export APNS_ENVIRONMENT="production"
```

## Security Considerations

1. Keep your private key secure and never commit it to version control
2. Use environment variables for sensitive configuration
3. Regularly rotate your APNS keys
4. Monitor failed deliveries and handle token invalidation

## Troubleshooting

### Common Issues

1. **Invalid Device Token**: Ensure the device token is correctly formatted and not expired
2. **Authentication Errors**: Verify your Key ID, Team ID, and private key are correct
3. **Environment Mismatch**: Use sandbox for development and production for live apps
4. **Bundle ID Mismatch**: Ensure the bundle ID matches your app's configuration

### Debug Mode

Enable debug logging by setting the log level in your application:

```go
log.SetLevel(log.DebugLevel)
``` 