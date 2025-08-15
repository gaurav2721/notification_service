# User Service Documentation

## Overview

The User Service is a comprehensive service that manages user information required for sending email, Slack, and in-app notifications. It provides preloaded user data with all necessary contact information and notification preferences.

## Features

### Core Functionality
- **User Management**: CRUD operations for users
- **Preloaded Data**: 8 sample users with comprehensive information
- **Notification Preferences**: Per-user notification channel preferences
- **Quiet Hours**: Configurable quiet hours for each user
- **Team/Department Filtering**: Get users by team, department, or role
- **Thread-Safe**: Concurrent access support with read-write mutex

### User Information Includes
- **Basic Info**: ID, name, email, phone number
- **Slack Integration**: Slack user ID, channel, workspace
- **Organization**: Company, department, role, manager
- **Team/Project**: Team IDs, project IDs
- **Preferences**: Notification channel preferences, quiet hours, timezone

## User Model Structure

### User
```go
type User struct {
    ID                string           `json:"id"`
    Email             string           `json:"email"`
    FirstName         string           `json:"first_name"`
    LastName          string           `json:"last_name"`
    FullName          string           `json:"full_name"`
    SlackUserID       string           `json:"slack_user_id,omitempty"`
    SlackChannel      string           `json:"slack_channel,omitempty"`
    SlackWorkspace    string           `json:"slack_workspace,omitempty"`
    PhoneNumber       string           `json:"phone_number,omitempty"`
    Company           string           `json:"company,omitempty"`
    Department        string           `json:"department,omitempty"`
    Role              string           `json:"role,omitempty"`
    ManagerID         string           `json:"manager_id,omitempty"`
    TeamIDs           []string         `json:"team_ids,omitempty"`
    ProjectIDs        []string         `json:"project_ids,omitempty"`
    Preferences       UserPreferences  `json:"preferences"`
    IsActive          bool             `json:"is_active"`
    LastLoginAt       *time.Time       `json:"last_login_at,omitempty"`
    CreatedAt         time.Time        `json:"created_at"`
    UpdatedAt         time.Time        `json:"updated_at"`
}
```

### UserPreferences
```go
type UserPreferences struct {
    EmailEnabled     bool   `json:"email_enabled"`
    SlackEnabled     bool   `json:"slack_enabled"`
    InAppEnabled     bool   `json:"in_app_enabled"`
    MarketingEmails  bool   `json:"marketing_emails"`
    UrgentOnly       bool   `json:"urgent_only"`
    QuietHoursStart  int    `json:"quiet_hours_start"` // 24-hour format (0-23)
    QuietHoursEnd    int    `json:"quiet_hours_end"`   // 24-hour format (0-23)
    Timezone         string `json:"timezone"`
}
```

## Usage Examples

### Basic Usage

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/gaurav2721/notification-service/services"
)

func main() {
    // Create user service
    userService := services.NewUserService()
    ctx := context.Background()
    
    // Get user by ID
    user, err := userService.GetUserByID(ctx, "user-001")
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("User: %s (%s)\n", user.FullName, user.Email)
}
```

### Notification Integration

```go
// Get notification info for a user
notificationInfo, err := userService.GetUserNotificationInfo(ctx, "user-001")
if err != nil {
    log.Fatal(err)
}

// Get enabled notification channels
channels, err := userService.GetUserNotificationChannels(ctx, "user-001")
if err != nil {
    log.Fatal(err)
}

// Check if user is in quiet hours
inQuietHours, err := userService.IsUserInQuietHours(ctx, "user-001")
if err != nil {
    log.Fatal(err)
}

if inQuietHours {
    fmt.Println("User is in quiet hours, only send urgent notifications")
}
```

### Team and Department Queries

```go
// Get all users in Engineering department
engineers, err := userService.GetUsersByDepartment(ctx, "Engineering")
if err != nil {
    log.Fatal(err)
}

// Get all users in a specific team
teamMembers, err := userService.GetUsersByTeam(ctx, "team-eng")
if err != nil {
    log.Fatal(err)
}

// Get all users with a specific role
managers, err := userService.GetUsersByRole(ctx, "Engineering Manager")
if err != nil {
    log.Fatal(err)
}
```

### Bulk Operations

```go
// Get multiple users by IDs
userIDs := []string{"user-001", "user-002", "user-003"}
users, err := userService.GetUsersByIDs(ctx, userIDs)
if err != nil {
    log.Fatal(err)
}

// Get notification info for multiple users
notificationInfos, err := userService.GetUsersNotificationInfo(ctx, userIDs)
if err != nil {
    log.Fatal(err)
}
```

## Preloaded User Data

The service comes with 8 preloaded users representing different roles and departments:

### User Details

| ID | Name | Email | Department | Role | Slack Channel |
|----|------|-------|------------|------|---------------|
| user-001 | John Doe | john.doe@company.com | Engineering | Senior Software Engineer | #general |
| user-002 | Jane Smith | jane.smith@company.com | Design | UI/UX Designer | #design |
| user-003 | Mike Johnson | mike.johnson@company.com | Marketing | Marketing Manager | #marketing |
| user-004 | Sarah Wilson | sarah.wilson@company.com | Sales | Sales Representative | #sales |
| user-005 | David Brown | david.brown@company.com | Engineering | Engineering Manager | #engineering |
| user-006 | Lisa Garcia | lisa.garcia@company.com | Marketing | Marketing Director | #marketing |
| user-007 | Robert Taylor | robert.taylor@company.com | Sales | Sales Manager | #sales |
| user-008 | Emma Davis | emma.davis@company.com | Executive | VP of Operations | #executives |

### Notification Preferences

Each user has different notification preferences:

- **user-001**: All channels enabled, quiet hours 10 PM - 8 AM
- **user-002**: Email + Slack enabled, in-app disabled, urgent only
- **user-003**: Email + In-app enabled, Slack disabled
- **user-004**: All channels enabled, quiet hours 8 PM - 8 AM
- **user-005**: All channels enabled, no quiet hours
- **user-006**: All channels enabled, quiet hours 10 PM - 7 AM
- **user-007**: All channels enabled, quiet hours 9 PM - 8 AM
- **user-008**: All channels enabled, no quiet hours (executive)

## API Reference

### UserService Interface

```go
type UserService interface {
    // Basic CRUD operations
    GetUserByID(ctx context.Context, userID string) (*models.User, error)
    GetUserByEmail(ctx context.Context, email string) (*models.User, error)
    GetUsersByIDs(ctx context.Context, userIDs []string) ([]*models.User, error)
    CreateUser(ctx context.Context, user *models.User) error
    UpdateUser(ctx context.Context, user *models.User) error
    DeleteUser(ctx context.Context, userID string) error
    
    // Query operations
    GetUsersByTeam(ctx context.Context, teamID string) ([]*models.User, error)
    GetUsersByDepartment(ctx context.Context, department string) ([]*models.User, error)
    GetUsersByRole(ctx context.Context, role string) ([]*models.User, error)
    GetActiveUsers(ctx context.Context) ([]*models.User, error)
    
    // Notification-specific operations
    GetUserNotificationInfo(ctx context.Context, userID string) (*models.UserNotificationInfo, error)
    GetUsersNotificationInfo(ctx context.Context, userIDs []string) ([]*models.UserNotificationInfo, error)
    IsUserInQuietHours(ctx context.Context, userID string) (bool, error)
    GetUserNotificationChannels(ctx context.Context, userID string) ([]string, error)
    
    // Utility operations
    ReloadUsers(ctx context.Context) error
}
```

## Error Handling

The service returns specific errors for different scenarios:

```go
var (
    ErrUserNotFound     = errors.New("user not found")
    ErrUserInactive     = errors.New("user is inactive")
    ErrUserAlreadyExists = errors.New("user already exists")
    ErrInvalidUserID    = errors.New("invalid user ID")
    ErrInvalidEmail     = errors.New("invalid email address")
)
```

## Integration with Notification System

The UserService is designed to work seamlessly with the notification system:

1. **Recipient Resolution**: Convert user IDs to actual contact information
2. **Channel Filtering**: Only send to enabled notification channels
3. **Quiet Hours**: Respect user quiet hours for non-urgent notifications
4. **Bulk Operations**: Efficiently handle multiple recipients

### Example Integration

```go
// In your notification service
func (s *NotificationService) SendNotificationToUsers(ctx context.Context, userIDs []string, notification *models.Notification) error {
    // Get user notification info
    userInfos, err := s.userService.GetUsersNotificationInfo(ctx, userIDs)
    if err != nil {
        return err
    }
    
    for _, userInfo := range userInfos {
        // Check quiet hours
        inQuietHours, err := s.userService.IsUserInQuietHours(ctx, userInfo.ID)
        if err != nil {
            continue
        }
        
        // Skip non-urgent notifications during quiet hours
        if inQuietHours && notification.Priority != models.UrgentPriority {
            continue
        }
        
        // Get user's preferred channels
        channels, err := s.userService.GetUserNotificationChannels(ctx, userInfo.ID)
        if err != nil {
            continue
        }
        
        // Send to each enabled channel
        for _, channel := range channels {
            switch channel {
            case "email":
                s.sendEmail(userInfo.Email, notification)
            case "slack":
                s.sendSlack(userInfo.SlackUserID, notification)
            case "in_app":
                s.sendInApp(userInfo.ID, notification)
            }
        }
    }
    
    return nil
}
```

## Testing

The service includes comprehensive unit tests covering all functionality:

```bash
# Run all user service tests
go test ./services -v -run "TestUserService"

# Run specific test
go test ./services -v -run "TestUserService_GetUserByID"
```

## Performance Considerations

- **In-Memory Storage**: Fast access for development and testing
- **Thread-Safe**: Concurrent read/write operations supported
- **Efficient Queries**: Optimized for common notification use cases
- **Bulk Operations**: Efficient handling of multiple users

## Future Enhancements

- Database persistence
- Caching layer
- User preference management API
- Integration with external user systems
- Advanced filtering and search capabilities 