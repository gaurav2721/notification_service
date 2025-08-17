This service is used to send email , slack and in-app(for ios(apple push notification) and android(firebase cloud messaging)) notifications .

### Quick Start

### Build and Run with Docker
1. `make docker-build`
2. `make docker-run`

### View Docker Container Output(Run this to check if the notification has been sent or not)
1. `make docker-exec`

### Quick Testing

For quick testing instructions and example API calls, please refer to [QUICK_TEST.md](QUICK_TEST.md).

### Detailed Testing

For comprehensive testing instructions, example API calls, and test data, please refer to [TEST.md](TEST.md). 

## Building and Running

For detailed instructions on building and running the notification service, please refer to [BUILD.md](BUILD.md).

## Api documentation 

For detailed API documentation, please refer to [API.md](API.md).

## Assumptions for this service

1. Interfaces for sending emails, Slack messages, and in-app notifications will be mocked in the first iteration to focus on building scalable service logic with features such as templates and scheduling.(If the .env does not have creds for the email, slack, apns and fcm , the information will be printed in a simple output/<service_name>.txt for eg output/email.txt)
2. When a customer raises a notification request, only user IDs will be provided as recipients. The service will retrieve other necessary details from the pre-stored user information.
3. Only text-based content will be supported for notifications in this iteration.
4. Each notification request raised by customer/user will be linked to only one notification type/channel in this iteration.
5. In-app notifications will be limited to mobile push notifications for iOS and Android in this iteration.
6. All the information is stored in memory , database persistence may be added later
7. User and UserDeviceInfo have been preloaded and the apis for these are disabled by default , since it is considered to be out of scope

## Development

Following things have been implemented 

1. Logging with different log levels for eg info, debug
2. Constants package 
3. Input validation has been done for the notification and template apis

Future enhancements 
1. Creating a requestId for observability 
2. Adding a CORS check for apis

Basic project structure explaining what each package is doing :

notification_service/
  main.go
  services/ -> creates a service container that basically has reference to all the external service objects and internal objects for eg email,slack,apns,fcm,user,consumer, notification_manager
  validation/ -> has the logic to validate inputs for notification and template apis

```
notification_service/
├── main.go                 # Application entry point
├── go.mod                  # Go module file
├── go.sum                  # Go module checksums
├── Makefile                # Build and deployment scripts
├── Dockerfile              # Docker configuration
├── docker-compose.yml      # Docker Compose configuration
├── env.example             # Environment variables template
├── README.md               # Project documentation
├── USER_README.md          # User-focused documentation
├── bin/                    # Compiled binaries
│   └── notification-service # Executable binary
├── docs/                   # Documentation files
├── examples/               # Example implementations
│   ├── user_service_integration.go      # User service integration example
│   └── service_architecture_example.go  # Service architecture example
├── models/                 # Data models and types
│   └── notification.go     # Notification data models
├── handlers/               # HTTP request handlers
│   └── notification_handlers.go  # Notification API handlers
├── routes/                 # API routing configuration
│   ├── routes.go           # Main routing setup
│   ├── notification_routes.go    # Notification endpoints
│   ├── user_routes.go      # User management endpoints
│   ├── template_routes.go  # Template management endpoints
│   ├── health_routes.go    # Health check endpoints
│   └── middleware/         # HTTP middleware
│       └── middleware.go   # Authentication and logging middleware
├── services/               # Core service implementations
│   ├── interfaces.go       # Service interfaces
│   ├── notification_service.go  # Main notification service
│   ├── email_service.go    # Email service implementation
│   ├── slack_service.go    # Slack service implementation
│   ├── inapp_service.go    # In-app service implementation
│   └── scheduler_service.go # Scheduler service implementation
├── notification_manager/   # Notification orchestration
│   ├── interface.go        # Manager interface
│   ├── notification_manager.go    # Main manager implementation
│   ├── notification_manager_test.go  # Manager tests
│   ├── config.go           # Manager configuration
│   └── errors.go           # Manager error definitions
├── external_services/      # External service integrations
│   ├── service_factory.go  # Service factory pattern
│   ├── service_provider.go # Service provider implementation
│   ├── email/              # Email service module
│   │   ├── interface.go    # Email service interface
│   │   ├── email_service.go # Email service implementation
│   │   ├── config.go       # Email configuration
│   │   └── errors.go       # Email error definitions
│   ├── slack/              # Slack service module
│   │   ├── interface.go    # Slack service interface
│   │   ├── slack_service.go # Slack service implementation
│   │   ├── config.go       # Slack configuration
│   │   └── errors.go       # Slack error definitions
│   ├── inapp/              # In-app notification module
│   │   ├── interface.go    # In-app service interface
│   │   ├── inapp_service.go # In-app service implementation
│   │   ├── config.go       # In-app configuration
│   │   └── errors.go       # In-app error definitions
│   └── user/               # User service module
│       ├── interface.go    # User service interface
│       ├── user_service.go # User service implementation
│       ├── user_service_test.go # User service tests
│       ├── config.go       # User service configuration
│       └── errors.go       # User service error definitions
└── helper/                 # Utility and helper functions
    └── scheduler/          # Scheduler utilities
        ├── interface.go    # Scheduler interface
        ├── scheduler.go    # Scheduler implementation
        ├── config.go       # Scheduler configuration
        └── errors.go       # Scheduler error definitions
```