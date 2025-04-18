# Insights
An API that ingests, stores, and analyzes login data for tenants

## Installation
[Setup Guide](Setup.MD)

## API Endpoints

### POST /api/login/new
Submit a login attempt 

#### Request Body
```json
{
  "tenant": "string",    // Required: Tenant identifier
  "user": "string",      // Required: Username that attempted the login
  "origin": "string",    // Required: IP address of the login
  "status": "string",    // Required: Either "success" or "failure"
  "timestamp": "string"  // Optional: ISO8601 datetime, defaults to current time
}
```

#### Responses
- `200 OK`: Login event already stored
- `201 Created`: Login event successfully stored
- `400 Bad Request`: Invalid payload or missing required fields
- `500 Internal Server Error`: Server-side error

### GET /api/login/suspicious
Retrieves origins with suspicious login activity based on failure thresholds

#### Authentication
Requires valid API key (X-API-Key)

#### Query Parameters
- `threshold` (optional): Minimum number of failures to be considered suspicious (default: 5)
- `minutes` (optional): Time window in minutes to look for failures (default: 3)
- `page` (optional): Page number for pagination (default: 1)
- `limit` (optional): Number of results per page (default: 10, max: 100)
- `order` (optional): Sort order, either "asc" or "desc" (default: "desc")

#### Response
```json
[
  {
    "origin": "string",  // IP address 
    "failCount": 0       // Number of failed login attempts
  }
]
```

#### Responses
- `200 OK`: Request successful (even if no results are found)
- `400 Bad Request`: Invalid query parameters
- `401 Unauthorized`: Missing or invalid API key
- `500 Internal Server Error`: Server-side error