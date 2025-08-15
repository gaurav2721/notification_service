# Template Manager Architecture

This document explains the new template manager architecture that abstracts all template-related logic into a dedicated component.

## Overview

The template manager is a dedicated component that handles all template-related operations including storage, validation, predefined templates, and template lifecycle management. It provides a clean separation of concerns and makes the codebase more maintainable.

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    Notification Service                     │
├─────────────────────────────────────────────────────────────┤
│  ┌─────────────────┐    ┌─────────────────────────────────┐ │
│  │   HTTP Handlers │    │    Notification Manager        │ │
│  │                 │    │                                 │ │
│  │ - CreateTemplate│    │ - SendNotification             │ │
│  │ - GetTemplates  │    │ - ScheduleNotification         │ │
│  │ - Use Templates │    │ - Template Validation          │ │
│  └─────────────────┘    └─────────────────────────────────┘ │
│           │                           │                     │
│           └───────────────────────────┼─────────────────────┤
│                                       │                     │
│  ┌─────────────────────────────────────────────────────────┐ │
│  │              Template Manager                          │ │
│  │                                                         │ │
│  │ - Template Storage & Retrieval                         │ │
│  │ - Predefined Templates Loading                         │ │
│  │ - Template Validation                                  │ │
│  │ - Template Lifecycle Management                        │ │
│  │ - Thread-Safe Operations                               │ │
│  └─────────────────────────────────────────────────────────┘ │
│                                       │                     │
│  ┌─────────────────────────────────────────────────────────┐ │
│  │                 Models                                 │ │
│  │                                                         │ │
│  │ - Template Structures                                  │ │
│  │ - Template Content Types                               │ │
│  │ - Error Definitions                                    │ │
│  │ - Predefined Templates                                 │ │
│  └─────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
```

## Components

### 1. Template Manager (`notification_manager/templates/template_manager.go`)

The core component that handles all template operations:

**Key Features:**
- **Template Storage**: In-memory storage with thread-safe access
- **Predefined Templates**: Automatic loading of predefined templates on startup
- **Template Validation**: Validates template content and required variables
- **Template Retrieval**: Multiple ways to retrieve templates (by ID, name, type)
- **Template Creation**: Creates new templates with proper validation
- **Template Statistics**: Provides counts and statistics

**Key Methods:**
```go
// Core operations
CreateTemplate(ctx context.Context, template *models.Template) (*models.TemplateResponse, error)
GetTemplateVersion(ctx context.Context, templateID string, version int) (*models.TemplateVersion, error)
GetPredefinedTemplates() []*models.Template

// Template retrieval
GetTemplateByID(templateID string) (*models.Template, error)
GetTemplateByName(name string) (*models.Template, error)
GetTemplatesByType(templateType models.NotificationType) []*models.Template

// Validation
ValidateTemplateData(templateID string, data map[string]interface{}) error

// Statistics
GetAllTemplates() []*models.Template
GetTemplateCount() int
GetPredefinedTemplateCount() int
```

### 2. Predefined Templates (`models/predefined_templates.go`)

Contains all predefined templates that are loaded on service startup:

**Included Templates:**
- **Email Templates**:
  - Welcome Email Template
  - Password Reset Template
  - Order Confirmation Template

- **Slack Templates**:
  - System Alert Template
  - Deployment Notification Template

- **In-App Templates**:
  - Order Status Update Template
  - Payment Reminder Template

**Helper Functions:**
```go
PredefinedTemplates() []*models.Template
GetTemplateByID(templateID string) *models.Template
GetTemplateByName(name string) *models.Template
GetTemplatesByType(templateType models.NotificationType) []*models.Template
```

### 3. Updated Notification Manager (`notification_manager/notification_manager.go`)

The notification manager now delegates all template operations to the template manager:

**Changes:**
- Removed direct template handling
- Added template manager dependency
- Delegates template operations to template manager
- Cleaner separation of concerns

### 4. Updated Handlers (`handlers/notification_handlers.go`)

HTTP handlers now use the template manager through the notification manager:

**New Endpoints:**
- `GET /api/v1/templates/predefined` - List all predefined templates
- `POST /api/v1/templates` - Create new templates
- `GET /api/v1/templates/{templateId}/versions/{version}` - Get specific template version

## Benefits of the New Architecture

### 1. **Separation of Concerns**
- Template logic is isolated in its own component
- Notification manager focuses on notification orchestration
- Handlers focus on HTTP request/response handling

### 2. **Reusability**
- Template manager can be used independently
- Easy to test template operations in isolation
- Can be extended for different storage backends

### 3. **Maintainability**
- Clear responsibility boundaries
- Easier to modify template logic without affecting other components
- Better code organization

### 4. **Thread Safety**
- Template manager provides thread-safe operations
- Concurrent access to templates is handled properly
- No race conditions in template operations

### 5. **Extensibility**
- Easy to add new template types
- Easy to add new template operations
- Easy to integrate with different storage systems

## Usage Examples

### 1. Using Template Manager Directly

```go
// Create template manager
templateManager := templates.NewTemplateManager()

// Get predefined templates
predefinedTemplates := templateManager.GetPredefinedTemplates()

// Create custom template
customTemplate := &models.Template{
    Name: "Custom Template",
    Type: models.EmailNotification,
    Content: models.TemplateContent{
        Subject:   "Welcome to {{company}}",
        EmailBody: "Hello {{name}}, welcome!",
    },
    RequiredVariables: []string{"name", "company"},
}

response, err := templateManager.CreateTemplate(ctx, customTemplate)

// Validate template data
err = templateManager.ValidateTemplateData(templateID, data)
```

### 2. Using Through Notification Manager

```go
// Template operations are automatically handled
notification := &Notification{
    Type: "email",
    Template: &models.TemplateData{
        ID: "550e8400-e29b-41d4-a716-446655440000",
        Data: map[string]interface{}{
            "name": "John Doe",
            "platform": "Tuskira",
        },
    },
    Recipients: []string{"user-123"},
}

// Template validation happens automatically
response, err := notificationManager.SendNotification(ctx, notification)
```

### 3. Using Through HTTP API

```bash
# Get predefined templates
curl -X GET http://localhost:8080/api/v1/templates/predefined

# Create custom template
curl -X POST http://localhost:8080/api/v1/templates \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Custom Template",
    "type": "email",
    "content": {
      "subject": "Welcome to {{company}}",
      "email_body": "Hello {{name}}, welcome!"
    },
    "required_variables": ["name", "company"]
  }'

# Use template in notification
curl -X POST http://localhost:8080/api/v1/notifications \
  -H "Content-Type: application/json" \
  -d '{
    "type": "email",
    "template": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "data": {
        "name": "John Doe",
        "platform": "Tuskira"
      }
    },
    "recipients": ["user-123"]
  }'
```

## Future Enhancements

### 1. **Database Integration**
- Replace in-memory storage with persistent database
- Support for template versioning
- Template backup and restore

### 2. **Template Processing**
- Variable substitution logic
- Template rendering engine
- Support for complex template logic

### 3. **Template Analytics**
- Template usage tracking
- Performance metrics
- Template effectiveness analysis

### 4. **Template Management UI**
- Web interface for template management
- Template preview functionality
- Template approval workflow

### 5. **Template Categories**
- Template categorization and tagging
- Template search and filtering
- Template organization features

## Testing

The template manager includes comprehensive testing capabilities:

```go
// Test template creation
func TestCreateTemplate(t *testing.T) {
    tm := templates.NewTemplateManager()
    template := &models.Template{...}
    response, err := tm.CreateTemplate(ctx, template)
    // Assertions...
}

// Test template validation
func TestValidateTemplateData(t *testing.T) {
    tm := templates.NewTemplateManager()
    err := tm.ValidateTemplateData(templateID, data)
    // Assertions...
}

// Test predefined templates
func TestPredefinedTemplates(t *testing.T) {
    tm := templates.NewTemplateManager()
    templates := tm.GetPredefinedTemplates()
    // Assertions...
}
```

## Conclusion

The new template manager architecture provides a clean, maintainable, and extensible solution for template management. It separates concerns properly, provides thread-safe operations, and makes the codebase more organized and easier to work with. 