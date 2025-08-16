# Firebase Cloud Messaging (FCM)

This package provides a Go implementation for sending push notifications to Android devices using Firebase Cloud Messaging (FCM).

## Features

- Send push notifications to Android devices
- Device token management (register/unregister)
- Batch processing for efficient delivery
- Support for both notification and data messages
- Thread-safe operations
- Comprehensive error handling
- Configurable batch sizes and timeouts

## Configuration

The FCM service requires the following configuration:

```go
config := &FCMConfig{
    ServerKey:  "YOUR_FCM_SERVER_KEY",
    ProjectID:  "your-project-id",
    MaxRetries: 3,
    Timeout:    30, // seconds
    BatchSize:  1000, // tokens per batch
}
```

### Required Credentials

1. **Server Key**: Your FCM server key from Firebase Console
2. **Project ID**: Your Firebase project identifier

## Usage

### Basic Usage

```go
package main

import (
    "context"
    "log"
    
    "github.com/gaurav2721/notification-service/external_services/fcm"
)

func main() {
    // Create FCM service
    service := fcm.NewFCMService()
    
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
            "data": map[string]interface{}{
                "key1": "value1",
                "key2": "value2",
            },
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
config := &fcm.FCMConfig{
    ServerKey:  "AIzaSyBxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
    ProjectID:  "my-firebase-project",
    MaxRetries: 3,
    Timeout:    30,
    BatchSize:  1000,
}

service, err := fcm.NewFCMServiceWithConfig(config)
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

## Message Types

### Notification Messages
These are automatically displayed by the system:

```go
Content: map[string]interface{}{
    "title": "Notification Title",
    "body":  "Notification Body",
}
```

### Data Messages
These are handled by your app:

```go
Content: map[string]interface{}{
    "title": "Data Message Title",
    "body":  "Data Message Body",
    "data": map[string]interface{}{
        "custom_key": "custom_value",
        "action":     "open_screen",
        "screen_id":  "123",
    },
}
```

## Error Handling

The service provides specific error types:

- `ErrFCMSendFailed`: Failed to send push notification
- `ErrInvalidDeviceToken`: Invalid device token provided
- `ErrInvalidConfiguration`: Invalid service configuration
- `ErrDeviceTokenNotFound`: Device token not found
- `ErrUserNotFound`: User not found
- `ErrInvalidNotificationPayload`: Invalid notification payload
- `ErrInvalidServerKey`: Invalid FCM server key

## Testing

Run the tests:

```bash
go test ./external_services/fcm
```

## Environment Variables

For production deployment, consider using environment variables:

```bash
export FCM_SERVER_KEY="YOUR_FCM_SERVER_KEY"
export FCM_PROJECT_ID="your-project-id"
export FCM_TIMEOUT="30"
export FCM_BATCH_SIZE="1000"
```

## Performance Considerations

### Batch Processing
The service automatically processes device tokens in batches to optimize performance:

- Default batch size: 1000 tokens
- Configurable via `BatchSize` in config
- Reduces API calls and improves throughput

### Rate Limiting
FCM has rate limits:
- 1000 requests per second per project
- 1 million messages per day per project
- Consider implementing retry logic for failed deliveries

## Security Considerations

1. Keep your FCM server key secure and never commit it to version control
2. Use environment variables for sensitive configuration
3. Regularly rotate your FCM server keys
4. Monitor failed deliveries and handle token invalidation
5. Implement proper authentication in your app

## Troubleshooting

### Common Issues

1. **Invalid Server Key**: Ensure your FCM server key is correct and active
2. **Invalid Device Token**: Ensure the device token is correctly formatted and not expired
3. **Rate Limiting**: Monitor your FCM usage and implement backoff strategies
4. **Network Issues**: Check your network connectivity and firewall settings

### Debug Mode

Enable debug logging by setting the log level in your application:

```go
log.SetLevel(log.DebugLevel)
```

### FCM Response Codes

Common FCM response codes:
- `200`: Success
- `400`: Bad Request (check payload format)
- `401`: Unauthorized (check server key)
- `500`: Internal Server Error (retry later)

## Integration with Firebase Console

1. Create a Firebase project in the Firebase Console
2. Add your Android app to the project
3. Download the `google-services.json` file
4. Get your server key from Project Settings > Cloud Messaging
5. Configure your Android app with the downloaded configuration

## Best Practices

1. **Token Management**: Always handle token refresh and invalidation
2. **Error Handling**: Implement proper error handling for failed deliveries
3. **Testing**: Test with both development and production FCM environments
4. **Monitoring**: Monitor delivery rates and user engagement
5. **Segmentation**: Use topics and user segments for targeted messaging 