# Routes Package Documentation

## 🎯 Overview

The routes package organizes all API routes in a modular and maintainable structure. Each route type has its own dedicated file, making it easy to manage and extend the API endpoints.

## 📁 Package Structure

```
routes/
├── README.md                    # This documentation file
├── routes.go                    # Main routes setup
├── notification_routes.go       # Notification endpoints
├── template_routes.go          # Template endpoints
├── user_routes.go              # User endpoints
├── health_routes.go            # Health check endpoints
└── middleware/
    └── middleware.go           # Middleware functions
```

## 🔧 Files Overview

### **Main Routes File** (`routes.go`)
- **Purpose**: Main entry point for route configuration
- **Function**: `SetupRoutes()` - Configures all routes and middleware
- **Usage**: Called from `main.go` to set up the entire routing structure

### **Notification Routes** (`notification_routes.go`)
- **Purpose**: Handles notification-related endpoints
- **Endpoints**:
  - `POST /api/v1/notifications` - Send notification
  - `GET /api/v1/notifications/:id` - Get notification status

### **Template Routes** (`template_routes.go`)
- **Purpose**: Handles template management endpoints
- **Endpoints**:
  - `POST /api/v1/templates` - Create template
  - `GET /api/v1/templates/:id` - Get template
  - `PUT /api/v1/templates/:id` - Update template
  - `DELETE /api/v1/templates/:id` - Delete template

### **User Routes** (`user_routes.go`)
- **Purpose**: Handles user management endpoints (future implementation)
- **Endpoints**:
  - `GET /api/v1/users` - Get all users
  - `GET /api/v1/users/:id` - Get user by ID
  - `POST /api/v1/users` - Create user
  - `PUT /api/v1/users/:id` - Update user
  - `DELETE /api/v1/users/:id` - Delete user

### **Health Routes** (`health_routes.go`)
- **Purpose**: Handles health check and monitoring endpoints
- **Endpoints**:
  - `GET /health` - Health check

### **Middleware** (`middleware/middleware.go`)
- **Purpose**: Centralized middleware management
- **Functions**:
  - `SetupMiddleware()` - Configure all middleware
  - `CORS()` - Cross-origin resource sharing
  - `RequestID()` - Add unique request IDs

## 🚀 Usage

### **Basic Setup**
```go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/gaurav2721/notification-service/routes"
    "github.com/gaurav2721/notification-service/handlers"
)

func main() {
    router := gin.Default()
    handler := handlers.NewNotificationHandler(notificationService)
    
    // Setup all routes
    routes.SetupRoutes(router, handler)
    
    router.Run(":8080")
}
```

### **Adding New Routes**
To add new routes, follow this pattern:

1. **Create a new route file** (e.g., `analytics_routes.go`):
```go
package routes

import (
    "github.com/gin-gonic/gin"
    "github.com/gaurav2721/notification-service/handlers"
)

func SetupAnalyticsRoutes(api *gin.RouterGroup, handler *handlers.NotificationHandler) {
    analytics := api.Group("/analytics")
    {
        analytics.GET("/stats", handler.GetStats)
        analytics.GET("/reports", handler.GetReports)
    }
}
```

2. **Add to main routes file** (`routes.go`):
```go
func SetupRoutes(router *gin.Engine, handler *handlers.NotificationHandler) {
    // ... existing setup ...
    
    api := router.Group("/api/v1")
    {
        // ... existing routes ...
        SetupAnalyticsRoutes(api, handler)  // Add new routes
    }
}
```

## 📋 API Endpoints Summary

### **Health Check**
- `GET /health` - Service health status

### **Notifications**
- `POST /api/v1/notifications` - Send notification
- `GET /api/v1/notifications/:id` - Get notification status

### **Templates**
- `POST /api/v1/templates` - Create notification template
- `GET /api/v1/templates/:id` - Get template by ID
- `PUT /api/v1/templates/:id` - Update template
- `DELETE /api/v1/templates/:id` - Delete template

### **Users** (Future)
- `GET /api/v1/users` - Get all users
- `GET /api/v1/users/:id` - Get user by ID
- `POST /api/v1/users` - Create user
- `PUT /api/v1/users/:id` - Update user
- `DELETE /api/v1/users/:id` - Delete user

## 🔒 Middleware

### **Default Middleware**
- **Logging**: Request/response logging
- **Recovery**: Panic recovery
- **CORS**: Cross-origin resource sharing (optional)
- **RequestID**: Unique request identification (optional)

### **Adding Custom Middleware**
```go
// In middleware/middleware.go
func AuthMiddleware() gin.HandlerFunc {
    return gin.HandlerFunc(func(c *gin.Context) {
        // Authentication logic
        token := c.GetHeader("Authorization")
        if token == "" {
            c.JSON(401, gin.H{"error": "Unauthorized"})
            c.Abort()
            return
        }
        c.Next()
    })
}

// In routes.go
func SetupRoutes(router *gin.Engine, handler *handlers.NotificationHandler) {
    middleware.SetupMiddleware(router)
    
    // Add custom middleware to specific routes
    api := router.Group("/api/v1")
    api.Use(middleware.AuthMiddleware())  // Apply to all API routes
    
    // ... rest of route setup
}
```

## 🧪 Testing Routes

### **Testing Individual Route Files**
```bash
# Test specific route functionality
go test ./routes -v

# Test with specific route file
go test ./routes -run TestNotificationRoutes
```

### **Integration Testing**
```go
func TestNotificationRoutes(t *testing.T) {
    router := gin.New()
    handler := &MockNotificationHandler{}
    
    SetupNotificationRoutes(router.Group("/api/v1"), handler)
    
    // Test POST /api/v1/notifications
    w := httptest.NewRecorder()
    req, _ := http.NewRequest("POST", "/api/v1/notifications", strings.NewReader(`{
        "type": "email",
        "title": "Test",
        "message": "Test message",
        "recipients": ["test@example.com"]
    }`))
    req.Header.Set("Content-Type", "application/json")
    
    router.ServeHTTP(w, req)
    
    assert.Equal(t, 200, w.Code)
}
```

## 🔧 Configuration

### **Environment Variables**
The routes package respects the following environment variables:
- `PORT` - Server port (default: 8080)
- `GIN_MODE` - Gin mode (debug, release, test)

### **Route Configuration**
Routes can be configured through environment variables:
```bash
# Enable CORS
export ENABLE_CORS=true

# Enable request ID tracking
export ENABLE_REQUEST_ID=true

# API version
export API_VERSION=v1
```

## 📈 Benefits

### **1. Modularity**
- Each route type has its own file
- Easy to add/remove route groups
- Clear separation of concerns

### **2. Maintainability**
- Organized structure
- Easy to find specific endpoints
- Consistent patterns

### **3. Scalability**
- Easy to add new route groups
- Middleware can be applied selectively
- Version control friendly

### **4. Testing**
- Routes can be tested independently
- Mock handlers for testing
- Clear test structure

## 🔮 Future Enhancements

### **Planned Features**
1. **Rate Limiting**: Add rate limiting middleware
2. **Authentication**: JWT-based authentication
3. **API Versioning**: Support for multiple API versions
4. **Documentation**: Auto-generated API documentation
5. **Monitoring**: Request metrics and monitoring

### **Route Groups to Add**
- **Analytics Routes**: Usage statistics and reports
- **Webhook Routes**: Webhook management
- **Settings Routes**: Service configuration
- **Admin Routes**: Administrative functions

## 📝 Best Practices

### **Route Organization**
1. Group related endpoints together
2. Use consistent naming conventions
3. Keep route files focused and small
4. Document all endpoints

### **Middleware Usage**
1. Apply global middleware in `SetupMiddleware()`
2. Apply route-specific middleware in route setup functions
3. Use middleware for cross-cutting concerns
4. Keep middleware functions simple and focused

### **Error Handling**
1. Use consistent error response format
2. Log errors appropriately
3. Return appropriate HTTP status codes
4. Provide meaningful error messages

The routes package provides a solid foundation for building scalable and maintainable APIs! 