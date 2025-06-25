# Stock Rating Microservice API Documentation

## Overview

This microservice provides a comprehensive API for managing stock ratings data. It includes functionality for creating, retrieving, and managing stock ratings, as well as an asynchronous data synchronization feature that downloads stock rating data from external APIs and stores it in CockroachDB. The service also includes advanced trading algorithms for finding optimal buy and sell opportunities.

## Base URL

```
http://localhost:8080
```

## Authentication

Currently, the API does not require authentication for most endpoints. However, the external API integration uses Bearer token authentication for accessing third-party data sources.

## Data Models

### Stock Rating Response (DTO)
```json
{
  "ticker": "AAPL",
  "target_from": "150.00",
  "target_to": "180.00",
  "company": "Apple Inc.",
  "action": "initiate",
  "brokerage": "Goldman Sachs",
  "rating_from": "buy",
  "rating_to": "strong_buy",
  "time": "2024-01-15T10:30:00Z"
}
```

### Job Response
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "status": "processing",
  "type": "external_api_sync",
  "progress": 500,
  "total_items": 1000,
  "error_message": null,
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:35:00Z",
  "completed_at": null
}
```

### Paginated Response
```json
{
  "data": [
    {
      "ticker": "AAPL",
      "target_from": "150.00",
      "target_to": "180.00",
      "company": "Apple Inc.",
      "action": "initiate",
      "brokerage": "Goldman Sachs",
      "rating_from": "buy",
      "rating_to": "strong_buy",
      "time": "2024-01-15T10:30:00Z"
    }
  ],
  "page": 1,
  "page_size": 20,
  "total_count": 150,
  "total_pages": 8,
  "has_next": true,
  "has_prev": false
}
```

### Trading Recommendation
```json
{
  "ticker": "AAPL",
  "buy_price": 150.00,
  "sell_price": 180.00,
  "max_profit": 30.00,
  "profit_percentage": 20.0,
  "buy_time": "2024-01-15T10:30:00Z",
  "sell_time": "2024-03-20T14:15:00Z",
  "buy_brokerage": "Goldman Sachs",
  "sell_brokerage": "Morgan Stanley",
  "buy_action": "initiate",
  "sell_action": "upgrade",
  "buy_rating": "buy",
  "sell_rating": "strong_buy",
  "buy_ticker": "AAPL",
  "sell_ticker": "AAPL",
  "total_data_points": 25,
  "date_range": {
    "start_date": "2024-01-15T10:30:00Z",
    "end_date": "2024-03-20T14:15:00Z"
  }
}
```

## Endpoints

### 1. Health Check

#### GET /api/hello
Simple health check endpoint to verify the service is running.

**Response:**
```json
{
  "message": "Hello World!"
}
```

**Status Codes:**
- `200 OK` - Service is running

---

### 2. External API Data Synchronization

#### GET /api/external/hello
**Asynchronous Data Import Endpoint**

This endpoint initiates an asynchronous job to download all stock rating data from the external API and store it in the database. This is the primary method for populating your stock ratings database with external data.

**How it works:**
1. Creates a new job with status "pending"
2. Returns immediately with the job ID
3. Processes data in the background using chunked processing (100 items per batch)
4. Updates job progress in real-time
5. Handles pagination automatically until all data is downloaded

**Request:**
```
GET /api/external/hello
```

**Response:**
```json
{
  "job_id": "550e8400-e29b-41d4-a716-446655440000",
  "status": "pending",
  "message": "Job created successfully. Use /api/jobs/{job_id} to check status."
}
```

**Status Codes:**
- `202 Accepted` - Job created successfully
- `500 Internal Server Error` - Failed to create job

**Job Processing Details:**
- **Timeout**: 30 minutes maximum processing time
- **Chunk Size**: 100 items per database batch
- **Pagination**: Automatically handles all pages from external API
- **Progress Tracking**: Real-time updates on items processed
- **Error Handling**: Detailed error messages for failed operations

**Job Status Values:**
- `pending` - Job created, waiting to start
- `processing` - Currently downloading and storing data
- `completed` - Successfully finished
- `failed` - Error occurred (check error_message for details)

---

### 3. Job Management

#### GET /api/jobs/{jobId}
**Job Status Check Endpoint**

Retrieve the current status and progress of an asynchronous job.

**Parameters:**
- `jobId` (path parameter) - UUID of the job to check

**Request:**
```
GET /api/jobs/550e8400-e29b-41d4-a716-446655440000
```

**Response:**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "status": "processing",
  "type": "external_api_sync",
  "progress": 500,
  "total_items": 1000,
  "error_message": null,
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:35:00Z",
  "completed_at": null
}
```

**Status Codes:**
- `200 OK` - Job found and returned
- `400 Bad Request` - Invalid job ID format
- `404 Not Found` - Job not found
- `500 Internal Server Error` - Database error

**Progress Tracking:**
- `progress` - Number of items processed so far
- `total_items` - Total number of items to process
- `updated_at` - Last time the job status was updated

---

### 4. Stock Rating Management

#### GET /api/stock-ratings
**Get Paginated Stock Ratings**

Retrieve stock ratings with pagination support. Returns different stock ratings from various tickers.

**Parameters:**
- `page` (query parameter, optional) - Page number (default: 1)
- `page_size` (query parameter, optional) - Items per page, 1-100 (default: 20)

**Request:**
```
GET /api/stock-ratings?page=1&page_size=10
```

**Response:**
```json
{
  "data": [
    {
      "ticker": "AAPL",
      "target_from": "150.00",
      "target_to": "180.00",
      "company": "Apple Inc.",
      "action": "initiate",
      "brokerage": "Goldman Sachs",
      "rating_from": "buy",
      "rating_to": "strong_buy",
      "time": "2024-01-15T10:30:00Z"
    }
  ],
  "page": 1,
  "page_size": 10,
  "total_count": 150,
  "total_pages": 15,
  "has_next": true,
  "has_prev": false
}
```

**Status Codes:**
- `200 OK` - Stock ratings found and returned
- `400 Bad Request` - Invalid pagination parameters
- `500 Internal Server Error` - Database error

---

#### POST /api/stock-ratings
**Create Single Stock Rating**

Create a new stock rating record in the database.

**Request Body:**
```json
{
  "ticker": "AAPL",
  "target_from": "150.00",
  "target_to": "180.00",
  "company": "Apple Inc.",
  "action": "initiate",
  "brokerage": "Goldman Sachs",
  "rating_from": "buy",
  "rating_to": "strong_buy",
  "time": "2024-01-15T10:30:00Z"
}
```

**Response:**
```json
{
  "ticker": "AAPL",
  "target_from": "150.00",
  "target_to": "180.00",
  "company": "Apple Inc.",
  "action": "initiate",
  "brokerage": "Goldman Sachs",
  "rating_from": "buy",
  "rating_to": "strong_buy",
  "time": "2024-01-15T10:30:00Z"
}
```

**Status Codes:**
- `201 Created` - Stock rating created successfully
- `400 Bad Request` - Invalid request payload
- `500 Internal Server Error` - Database error

---

#### POST /api/stock-ratings/batch
**Create Multiple Stock Ratings**

Create multiple stock rating records in a single request. Useful for bulk data import.

**Request Body:**
```json
[
  {
    "ticker": "AAPL",
    "target_from": "150.00",
    "target_to": "180.00",
    "company": "Apple Inc.",
    "action": "initiate",
    "brokerage": "Goldman Sachs",
    "rating_from": "buy",
    "rating_to": "strong_buy",
    "time": "2024-01-15T10:30:00Z"
  },
  {
    "ticker": "GOOGL",
    "target_from": "2800.00",
    "target_to": "3200.00",
    "company": "Alphabet Inc.",
    "action": "upgrade",
    "brokerage": "Morgan Stanley",
    "rating_from": "hold",
    "rating_to": "buy",
    "time": "2024-01-15T11:00:00Z"
  }
]
```

**Response:**
```json
[
  {
    "ticker": "AAPL",
    "target_from": "150.00",
    "target_to": "180.00",
    "company": "Apple Inc.",
    "action": "initiate",
    "brokerage": "Goldman Sachs",
    "rating_from": "buy",
    "rating_to": "strong_buy",
    "time": "2024-01-15T10:30:00Z"
  },
  {
    "ticker": "GOOGL",
    "target_from": "2800.00",
    "target_to": "3200.00",
    "company": "Alphabet Inc.",
    "action": "upgrade",
    "brokerage": "Morgan Stanley",
    "rating_from": "hold",
    "rating_to": "buy",
    "time": "2024-01-15T11:00:00Z"
  }
]
```

**Status Codes:**
- `201 Created` - Stock ratings created successfully
- `400 Bad Request` - Invalid request payload
- `500 Internal Server Error` - Database error

---

#### GET /api/stock-ratings/{id}
**Get Stock Rating by ID**

Retrieve a specific stock rating by its database ID.

**Parameters:**
- `id` (path parameter) - Numeric ID of the stock rating

**Request:**
```
GET /api/stock-ratings/123
```

**Response:**
```json
{
  "ticker": "AAPL",
  "target_from": "150.00",
  "target_to": "180.00",
  "company": "Apple Inc.",
  "action": "initiate",
  "brokerage": "Goldman Sachs",
  "rating_from": "buy",
  "rating_to": "strong_buy",
  "time": "2024-01-15T10:30:00Z"
}
```

**Status Codes:**
- `200 OK` - Stock rating found and returned
- `400 Bad Request` - Invalid ID format
- `404 Not Found` - Stock rating not found
- `500 Internal Server Error` - Database error

---

#### GET /api/stock-ratings/ticker/{ticker}
**Get All Stock Ratings by Ticker**

Retrieve all stock ratings for a specific stock ticker symbol.

**Parameters:**
- `ticker` (path parameter) - Stock ticker symbol (e.g., AAPL, GOOGL)

**Request:**
```
GET /api/stock-ratings/ticker/AAPL
```

**Response:**
```json
[
  {
    "ticker": "AAPL",
    "target_from": "150.00",
    "target_to": "180.00",
    "company": "Apple Inc.",
    "action": "initiate",
    "brokerage": "Goldman Sachs",
    "rating_from": "buy",
    "rating_to": "strong_buy",
    "time": "2024-01-15T10:30:00Z"
  },
  {
    "ticker": "AAPL",
    "target_from": "160.00",
    "target_to": "190.00",
    "company": "Apple Inc.",
    "action": "upgrade",
    "brokerage": "JPMorgan",
    "rating_from": "hold",
    "rating_to": "buy",
    "time": "2024-01-16T09:15:00Z"
  }
]
```

**Status Codes:**
- `200 OK` - Stock ratings found and returned
- `400 Bad Request` - Ticker parameter missing
- `500 Internal Server Error` - Database error

---

#### GET /api/stock-ratings/ticker/{ticker}/latest
**Get Latest Stock Rating by Ticker**

Retrieve the most recent stock rating for a specific ticker symbol.

**Parameters:**
- `ticker` (path parameter) - Stock ticker symbol (e.g., AAPL, GOOGL)

**Request:**
```
GET /api/stock-ratings/ticker/AAPL/latest
```

**Response:**
```json
{
  "ticker": "AAPL",
  "target_from": "160.00",
  "target_to": "190.00",
  "company": "Apple Inc.",
  "action": "upgrade",
  "brokerage": "JPMorgan",
  "rating_from": "hold",
  "rating_to": "buy",
  "time": "2024-01-16T09:15:00Z"
}
```

**Status Codes:**
- `200 OK` - Latest stock rating found and returned
- `400 Bad Request` - Ticker parameter missing
- `404 Not Found` - No stock rating found for ticker
- `500 Internal Server Error` - Database error

---

### 5. Trading Algorithms

#### GET /api/algorithms/best-time-to-buy-sell/{ticker}
**Single Ticker Trading Analysis**

Analyzes a single ticker to find the optimal buy and sell points based on analyst target prices using the "Best Time to Buy and Sell Stock" algorithm.

**Parameters:**
- `ticker` (path parameter) - Stock ticker symbol
- `start_date` (query parameter, optional) - Start date for analysis (YYYY-MM-DD)
- `end_date` (query parameter, optional) - End date for analysis (YYYY-MM-DD)

**Request:**
```
GET /api/algorithms/best-time-to-buy-sell/AAPL?start_date=2024-01-01&end_date=2024-12-31
```

**Response:**
```json
{
  "ticker": "AAPL",
  "buy_price": 150.00,
  "sell_price": 180.00,
  "max_profit": 30.00,
  "profit_percentage": 20.0,
  "buy_time": "2024-01-15T10:30:00Z",
  "sell_time": "2024-03-20T14:15:00Z",
  "buy_brokerage": "Goldman Sachs",
  "sell_brokerage": "Morgan Stanley",
  "buy_action": "initiate",
  "sell_action": "upgrade",
  "buy_rating": "buy",
  "sell_rating": "strong_buy",
  "buy_ticker": "AAPL",
  "sell_ticker": "AAPL",
  "total_data_points": 25,
  "date_range": {
    "start_date": "2024-01-15T10:30:00Z",
    "end_date": "2024-03-20T14:15:00Z"
  }
}
```

**Status Codes:**
- `200 OK` - Analysis completed successfully
- `400 Bad Request` - Invalid parameters or date format
- `404 Not Found` - No data found for ticker
- `500 Internal Server Error` - Algorithm error

**Algorithm Details:**
- **Time Complexity**: O(n) where n is the number of price points
- **Data Source**: Uses `target_from` prices from analyst ratings
- **Minimum Data**: Requires at least 2 price points
- **Date Filtering**: Optional date range filtering for analysis

---

#### POST /api/algorithms/best-time-to-buy-sell/multiple
**Multiple Ticker Trading Analysis**

Analyzes multiple tickers and returns recommendations sorted by profit percentage.

**Request Body:**
```json
{
  "tickers": ["AAPL", "GOOGL", "MSFT", "TSLA"],
  "start_date": "2024-01-01T00:00:00Z",
  "end_date": "2024-12-31T23:59:59Z"
}
```

**Response:**
```json
[
  {
    "ticker": "TSLA",
    "buy_price": 200.00,
    "sell_price": 280.00,
    "max_profit": 80.00,
    "profit_percentage": 40.0,
    "buy_time": "2024-02-01T09:00:00Z",
    "sell_time": "2024-04-15T16:30:00Z",
    "buy_brokerage": "Deutsche Bank",
    "sell_brokerage": "Citigroup",
    "buy_action": "initiate",
    "sell_action": "upgrade",
    "buy_rating": "hold",
    "sell_rating": "buy",
    "buy_ticker": "TSLA",
    "sell_ticker": "TSLA",
    "total_data_points": 18,
    "date_range": {
      "start_date": "2024-02-01T09:00:00Z",
      "end_date": "2024-04-15T16:30:00Z"
    }
  },
  {
    "ticker": "AAPL",
    "buy_price": 150.00,
    "sell_price": 180.00,
    "max_profit": 30.00,
    "profit_percentage": 20.0,
    "buy_time": "2024-01-15T10:30:00Z",
    "sell_time": "2024-03-20T14:15:00Z",
    "buy_brokerage": "Goldman Sachs",
    "sell_brokerage": "Morgan Stanley",
    "buy_action": "initiate",
    "sell_action": "upgrade",
    "buy_rating": "buy",
    "sell_rating": "strong_buy",
    "buy_ticker": "AAPL",
    "sell_ticker": "AAPL",
    "total_data_points": 25,
    "date_range": {
      "start_date": "2024-01-15T10:30:00Z",
      "end_date": "2024-03-20T14:15:00Z"
    }
  }
]
```

**Status Codes:**
- `200 OK` - Analysis completed successfully
- `400 Bad Request` - Invalid request payload
- `500 Internal Server Error` - Algorithm error

---

#### GET /api/algorithms/best-time-to-buy-sell/global
**Global Market Trading Analysis**

Analyzes all stock ratings as if they belonged to the same ticker, providing a global market perspective across all available stocks.

**Parameters:**
- `start_date` (query parameter, optional) - Start date for analysis (YYYY-MM-DD)
- `end_date` (query parameter, optional) - End date for analysis (YYYY-MM-DD)

**Request:**
```
GET /api/algorithms/best-time-to-buy-sell/global?start_date=2024-01-01&end_date=2024-12-31
```

**Response:**
```json
{
  "ticker": "GLOBAL",
  "buy_price": 150.00,
  "sell_price": 280.00,
  "max_profit": 130.00,
  "profit_percentage": 86.67,
  "buy_time": "2024-01-15T10:30:00Z",
  "sell_time": "2024-04-20T14:15:00Z",
  "buy_brokerage": "Goldman Sachs",
  "sell_brokerage": "Morgan Stanley",
  "buy_action": "initiate",
  "sell_action": "upgrade",
  "buy_rating": "buy",
  "sell_rating": "strong_buy",
  "buy_ticker": "AAPL",
  "sell_ticker": "TSLA",
  "total_data_points": 1250,
  "date_range": {
    "start_date": "2024-01-15T10:30:00Z",
    "end_date": "2024-04-20T14:15:00Z"
  }
}
```

**Status Codes:**
- `200 OK` - Analysis completed successfully
- `400 Bad Request` - Invalid date format
- `500 Internal Server Error` - Algorithm error

**Global Analysis Features:**
- **Cross-Ticker Opportunities**: Can find opportunities where you buy one stock and sell another
- **Market-Wide Perspective**: Provides insights into the best overall trading opportunities
- **Unified Dataset**: Treats all stock ratings as one continuous price series
- **Large Dataset Handling**: Uses pagination to efficiently process large amounts of data

---

## Usage Examples

### Complete Workflow Example

1. **Start Data Import:**
```bash
curl -X GET http://localhost:8080/api/external/hello
```
Response:
```json
{
  "job_id": "550e8400-e29b-41d4-a716-446655440000",
  "status": "pending",
  "message": "Job created successfully. Use /api/jobs/{job_id} to check status."
}
```

2. **Check Job Progress:**
```bash
curl -X GET http://localhost:8080/api/jobs/550e8400-e29b-41d4-a716-446655440000
```
Response:
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "status": "processing",
  "type": "external_api_sync",
  "progress": 500,
  "total_items": 1000,
  "error_message": null,
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:35:00Z",
  "completed_at": null
}
```

3. **Query Stock Ratings:**
```bash
# Get paginated ratings
curl -X GET http://localhost:8080/api/stock-ratings?page=1&page_size=20

# Get all ratings for Apple
curl -X GET http://localhost:8080/api/stock-ratings/ticker/AAPL

# Get latest rating for Apple
curl -X GET http://localhost:8080/api/stock-ratings/ticker/AAPL/latest

# Get specific rating by ID
curl -X GET http://localhost:8080/api/stock-ratings/123
```

4. **Run Trading Analysis:**
```bash
# Analyze single ticker
curl -X GET "http://localhost:8080/api/algorithms/best-time-to-buy-sell/AAPL?start_date=2024-01-01&end_date=2024-12-31"

# Analyze multiple tickers
curl -X POST http://localhost:8080/api/algorithms/best-time-to-buy-sell/multiple \
  -H "Content-Type: application/json" \
  -d '{
    "tickers": ["AAPL", "GOOGL", "MSFT", "TSLA"],
    "start_date": "2024-01-01T00:00:00Z",
    "end_date": "2024-12-31T23:59:59Z"
  }'

# Global market analysis
curl -X GET "http://localhost:8080/api/algorithms/best-time-to-buy-sell/global?start_date=2024-01-01&end_date=2024-12-31"
```

### Manual Data Creation Example

```bash
# Create a single stock rating
curl -X POST http://localhost:8080/api/stock-ratings \
  -H "Content-Type: application/json" \
  -d '{
    "ticker": "TSLA",
    "target_from": "200.00",
    "target_to": "250.00",
    "company": "Tesla Inc.",
    "action": "initiate",
    "brokerage": "Deutsche Bank",
    "rating_from": "hold",
    "rating_to": "buy",
    "time": "2024-01-15T14:30:00Z"
  }'

# Create multiple stock ratings
curl -X POST http://localhost:8080/api/stock-ratings/batch \
  -H "Content-Type: application/json" \
  -d '[
    {
      "ticker": "MSFT",
      "target_from": "350.00",
      "target_to": "400.00",
      "company": "Microsoft Corporation",
      "action": "upgrade",
      "brokerage": "Citigroup",
      "rating_from": "buy",
      "rating_to": "strong_buy",
      "time": "2024-01-15T15:00:00Z"
    },
    {
      "ticker": "AMZN",
      "target_from": "140.00",
      "target_to": "160.00",
      "company": "Amazon.com Inc.",
      "action": "initiate",
      "brokerage": "Barclays",
      "rating_from": "hold",
      "rating_to": "buy",
      "time": "2024-01-15T15:30:00Z"
    }
  ]'
```

## Trading Algorithm Details

### Best Time to Buy and Sell Stock Algorithm

The service implements the classic "Best Time to Buy and Sell Stock" algorithm with the following characteristics:

**Algorithm Features:**
- **Time Complexity**: O(n) - optimal solution
- **Data Source**: Uses analyst target prices (`target_from` field)
- **Chronological Ordering**: Processes data in time sequence
- **Date Range Filtering**: Optional date range constraints
- **Cross-Ticker Analysis**: Global analysis treats all stocks as one dataset

**Algorithm Logic:**
1. **Data Collection**: Retrieves stock ratings for specified ticker(s)
2. **Price Extraction**: Uses `target_from` prices from analyst ratings
3. **Chronological Sorting**: Orders data by time for accurate analysis
4. **Profit Calculation**: Tracks minimum price and calculates potential profits
5. **Optimal Selection**: Finds the maximum profit opportunity

**Use Cases:**
- **Portfolio Optimization**: Find best trading opportunities
- **Market Timing**: Identify optimal entry/exit points
- **Cross-Sector Analysis**: Compare opportunities across industries
- **Risk Assessment**: Understand maximum potential profits

## Error Handling

All endpoints return consistent error responses:

```json
{
  "error": "Detailed error message"
}
```

Common error scenarios:
- **400 Bad Request**: Invalid request format, missing parameters, invalid date format
- **404 Not Found**: Resource not found, no data for ticker
- **500 Internal Server Error**: Server-side errors, database issues, algorithm errors

## Rate Limiting

Currently, no rate limiting is implemented. Consider implementing rate limiting for production use.

## Database Schema

The service uses CockroachDB with the following main tables:

### stock_ratings
- `id` (BIGSERIAL PRIMARY KEY)
- `ticker` (VARCHAR(10) NOT NULL)
- `target_from`, `target_to` (VARCHAR(50))
- `company` (VARCHAR(255))
- `action` (VARCHAR(50))
- `brokerage` (VARCHAR(255))
- `rating_from`, `rating_to` (VARCHAR(50))
- `time` (TIMESTAMP WITH TIME ZONE)
- `created_at`, `updated_at` (TIMESTAMP WITH TIME ZONE)

### jobs
- `id` (UUID PRIMARY KEY)
- `status` (VARCHAR(20) NOT NULL)
- `type` (VARCHAR(50) NOT NULL)
- `progress` (INTEGER DEFAULT 0)
- `total_items` (INTEGER DEFAULT 0)
- `error_message` (TEXT)
- `created_at`, `updated_at`, `completed_at` (TIMESTAMP WITH TIME ZONE)

## Configuration

The service requires a `config/config.yml` file with:

```yaml
database:
  host: "localhost"
  port: 26257
  user: "your_user"
  password: "your_password"
  dbname: "defaultdb"
  sslmode: "disable"

external_api:
  base_url: "https://api.example.com"
  timeout: 30
  token: "your_bearer_token"
```

## Monitoring and Logging

- All API requests are logged with middleware
- Job progress is tracked in real-time
- Error messages are detailed and actionable
- Database operations are wrapped in transactions for consistency
- Algorithm performance is monitored and logged

## Performance Considerations

- **Pagination**: Large datasets are paginated for efficient retrieval
- **Chunked Processing**: Batch operations for better performance
- **Indexing**: Database indexes on commonly queried fields
- **Async Processing**: Long-running operations are handled asynchronously
- **Caching**: Consider implementing caching for frequently accessed data

## Security Considerations

- **Input Validation**: All inputs are validated and sanitized
- **SQL Injection Prevention**: Uses parameterized queries
- **Error Information**: Error messages don't expose sensitive information
- **Rate Limiting**: Consider implementing rate limiting for production
- **Authentication**: Consider adding authentication for production use 