# API Key Authentication

The notification service now supports API key authentication for all `/api/v1` endpoints.

## Configuration

Set the `API_KEY` environment variable in your `.env` file:

```bash
API_KEY=your-secure-api-key-here
```

## Usage

### Setting the API Key

You can provide the API key in the `Authorization` header using any of these formats:

1. **Bearer token format:**
   ```
   Authorization: Bearer your-api-key-here
   ```

2. **ApiKey format:**
   ```
   Authorization: ApiKey your-api-key-here
   ```

3. **Direct API key:**
   ```
   Authorization: your-api-key-here
   ```

### Example Requests

#### cURL Example
```bash
curl -X POST http://localhost:8080/api/v1/notifications \
  -H "Authorization: Bearer your-api-key-here" \
  -H "Content-Type: application/json" \
  -d '{
    "type": "email",
    "recipient": "user@example.com",
    "subject": "Test Notification",
    "body": "This is a test notification"
  }'
```

#### JavaScript/Fetch Example
```javascript
const response = await fetch('http://localhost:8080/api/v1/notifications', {
  method: 'POST',
  headers: {
    'Authorization': 'Bearer your-api-key-here',
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
    type: 'email',
    recipient: 'user@example.com',
    subject: 'Test Notification',
    body: 'This is a test notification'
  })
});
```

## Error Responses

### Missing API Key
```json
{
  "error": "API key is required",
  "message": "Please provide an API key in the Authorization header"
}
```

### Invalid API Key
```json
{
  "error": "Invalid API key",
  "message": "The provided API key is invalid"
}
```

## Development Mode

If no `API_KEY` environment variable is set, the middleware will allow all requests (useful for development). In production, always set a secure API key.

## Security Best Practices

1. **Use a strong, random API key** - Generate a secure random string (at least 32 characters)
2. **Keep your API key secret** - Never commit API keys to version control
3. **Rotate API keys regularly** - Change your API keys periodically
4. **Use HTTPS in production** - Always use HTTPS to protect API keys in transit
5. **Monitor API usage** - Log and monitor API key usage for security

## Protected Endpoints

All endpoints under `/api/v1/*` require API key authentication:

- `POST /api/v1/notifications` - Send notifications
- `GET /api/v1/notifications/:id` - Get notification status
- `POST /api/v1/templates` - Create notification templates
- `GET /api/v1/templates` - List notification templates
- `GET /api/v1/users` - List users (if enabled)
- `POST /api/v1/users` - Create users (if enabled)

## Health Check Endpoints

Health check endpoints (typically under `/health/*`) do not require API key authentication. 