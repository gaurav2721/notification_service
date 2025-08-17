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

1. Interfaces for sending emails, slack messages, and in-app notifications will be mocked in the first iteration to focus on building scalable service logic(for eg using pub-sub, worker pool design pattern etc) with features such as templates and scheduling.(If the .env does not have creds for the email, slack, apns and fcm , the information will be printed in a simple output/<service_name>.txt for eg output/email.txt , output/slack.txt etc.)
2. When a customer raises a notification request, only user IDs will be provided as recipients. The service will retrieve other necessary details from the pre-stored user information.
3. Only text-based content will be supported for notifications in this iteration.
4. Each notification request raised by customer/user will be linked to only one notification type in this iteration.
5. In-app notifications will be limited to mobile push notifications for iOS and Android in this iteration.
6. All the information is stored in memory , database persistence will be added in further iterations
7. User and UserDeviceInfo have been preloaded and the apis for these are disabled by default , since it is considered to be out of scope for this iteration

## Development

Following things have been implemented 

Functional Requirement
1. For different notification types appropriate channel routing has been implemented(we have used buffered channels in golang as queues in the project)
2. Notification Scheduling – Two types:
    a) Immediate
    b) Scheduled for a later time
3. Notification Templates(Customization using parameters has been implemented) :
    a) Predefined notification templates are available 
    b) Ability for users to create their own templates and use them

Other things implemented are
1. Authentication via Api key
2. Input Validation for notification and template apis
3. Implemented logging with different levels 

Design Patterns implemented
1. Publisher/Subscriber
2. Worker Pool

Future enhancements 
1. Adding retries for enhancing reliability
2. Adding database persistence for enhancing reliability
3. Implementing event tracking 

Basic project structure explaining what each package is doing :

```
notification_service/
  main.go
  services/ -> creates a service container that basically has reference to all the external service objects and internal objects for eg email,slack,apns,fcm,user,consumer, notification_manager
  validation/ -> has the logic to validate inputs for notification and template apis
  routes/ -> defines all the routes for notification,user,templates
  notification_manager/ -> handles all the business logic for notifications for eg scheduling, templates, pushing to the appropriate channel
  models/ -> defines all the models
  logger/ -> sets up logger 
  handlers -> defines handlers for all the apis
  external_services/ -> has logic for all the services that notification_manager would require
    apns/ -> Apple Push Notification Service
    email/ -> email service
    fcm/ -> firebase cloud messaging 
    slack/ -> slack service
    user/ -> user service 
    kafka/ -> kafka service having apns,fcm,email and slack queue
    consumers/ -> consumer/workers that read from kafka queue and send notification via appropriate service for eg email,slack,apns,fcm service
```

Data Flow

```
Api -> Notification Manager -> Kafka -> Consumers -> Email/Slack/APNS/FCM Service
```