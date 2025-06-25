# Truora Stock Rating Microservice

A high-performance Go microservice for managing stock ratings data with advanced trading algorithms and external API integration.

## 🚀 Features

- **Stock Rating Management**: Create, retrieve, and manage stock ratings with full CRUD operations
- **External API Integration**: Asynchronous data synchronization from external stock rating APIs
- **Advanced Trading Algorithms**: Find optimal buy/sell opportunities using sophisticated algorithms
- **RESTful API**: Clean, well-documented REST endpoints with proper error handling
- **Database Integration**: PostgreSQL/CockroachDB support with GORM ORM
- **Job Management**: Track asynchronous operations with real-time progress updates
- **Pagination**: Efficient data retrieval with configurable pagination
- **Comprehensive Logging**: Detailed logging for debugging and monitoring

## 🏗️ Architecture

The service follows a clean architecture pattern with clear separation of concerns:

```
├── cmd/api/                    # Application entry point
├── config/                     # Configuration files
├── internal/
│   ├── database/              # Database migrations and setup
│   ├── delivery/http/         # HTTP handlers and routing
│   ├── domain/                # Core business entities
│   ├── dto/                   # Data Transfer Objects
│   ├── repository/            # Data access layer
│   └── usecase/               # Business logic and services
├── migrations/                # Database migration files
└── docs/                      # API documentation
```

## 🛠️ Tech Stack

- **Language**: Go 1.19+
- **Framework**: Chi router for HTTP routing
- **Database**: PostgreSQL/CockroachDB with GORM
- **Migrations**: golang-migrate
- **Configuration**: YAML-based config with environment variable support
- **External APIs**: HTTP client with timeout and retry logic

## 📋 Prerequisites

- Go 1.19 or higher
- PostgreSQL or CockroachDB database
- Access to external stock rating API (configured in config.yml)

## 🚀 Quick Start

### 1. Clone the Repository

```bash
git clone <repository-url>
cd truora
```

### 2. Install Dependencies

```bash
go mod download
```

### 3. Configure the Application

Copy and modify the configuration file:

```bash
cp config/config.yml config/config.local.yml
```

Update the database and external API settings in `config/config.local.yml`:

```yaml
server:
  port: 8080
  host: "localhost"

database:
  host: "your-database-host"
  port: 26257
  user: "your-username"
  password: "your-password"
  dbname: "your-database"
  sslmode: "verify-full"

external_api:
  base_url: "https://your-external-api.com"
  timeout: 30
  token: "your-api-token"
```

### 4. Run Database Migrations

```bash
go run cmd/api/main.go migrate
```

### 5. Start the Service

```bash
go run cmd/api/main.go
```

The service will be available at `http://localhost:8080`

## 📚 API Documentation

### Core Endpoints

#### Health Check
```http
GET /api/hello
```

#### External Data Synchronization
```http
GET /api/external/hello
```
Initiates asynchronous import of stock rating data from external API.

#### Job Management
```http
GET /api/jobs/{jobId}
```
Check the status of asynchronous jobs.

### Stock Rating Endpoints

#### Get Paginated Stock Ratings
```http
GET /api/stock-ratings?page=1&page_size=20
```

#### Create Stock Rating
```http
POST /api/stock-ratings
Content-Type: application/json

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

#### Get Stock Rating by ID
```http
GET /api/stock-ratings/{id}
```

#### Get Stock Ratings by Ticker
```http
GET /api/stock-ratings/ticker/{ticker}
```

#### Get Latest Stock Rating by Ticker
```http
GET /api/stock-ratings/ticker/{ticker}/latest
```

#### Create Batch Stock Ratings
```http
POST /api/stock-ratings/batch
Content-Type: application/json

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
  }
]
```

### Trading Algorithm Endpoints

#### Get Best Time to Buy/Sell for Single Ticker
```http
GET /api/algorithms/best-time-to-buy-sell/{ticker}?start_date=2024-01-01&end_date=2024-12-31
```

#### Get Best Time to Buy/Sell for Multiple Tickers
```http
POST /api/algorithms/best-time-to-buy-sell/multiple
Content-Type: application/json

{
  "tickers": ["AAPL", "GOOGL", "MSFT"],
  "start_date": "2024-01-01T00:00:00Z",
  "end_date": "2024-12-31T23:59:59Z"
}
```

#### Get Global Best Time to Buy/Sell
```http
GET /api/algorithms/best-time-to-buy-sell/global?start_date=2024-01-01&end_date=2024-12-31
```

## 🔧 Configuration

### Environment Variables

The service supports environment variables for configuration:

- `DB_HOST`: Database host
- `DB_PORT`: Database port
- `DB_USER`: Database username
- `DB_PASSWORD`: Database password
- `DB_NAME`: Database name
- `DB_SSLMODE`: Database SSL mode
- `EXTERNAL_API_BASE_URL`: External API base URL
- `EXTERNAL_API_TIMEOUT`: External API timeout in seconds
- `EXTERNAL_API_TOKEN`: External API authentication token

### Configuration File

The main configuration is in `config/config.yml`:

```yaml
server:
  port: 8080
  host: "localhost"

database:
  host: "your-database-host"
  port: 26257
  user: "your-username"
  password: "your-password"
  dbname: "your-database"
  sslmode: "verify-full"

external_api:
  base_url: "https://your-external-api.com"
  timeout: 30
  token: "your-api-token"
```

## 🗄️ Database Schema

### Stock Ratings Table
```sql
CREATE TABLE stock_ratings (
    id SERIAL PRIMARY KEY,
    ticker VARCHAR(10) NOT NULL,
    target_from DECIMAL(10,2),
    target_to DECIMAL(10,2),
    company VARCHAR(255),
    action VARCHAR(50),
    brokerage VARCHAR(255),
    rating_from VARCHAR(50),
    rating_to VARCHAR(50),
    time TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Jobs Table
```sql
CREATE TABLE jobs (
    id UUID PRIMARY KEY,
    status VARCHAR(20) NOT NULL,
    type VARCHAR(50) NOT NULL,
    progress INTEGER DEFAULT 0,
    total_items INTEGER DEFAULT 0,
    error_message TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP
);
```

## 🔄 Data Synchronization

The service includes an asynchronous data synchronization feature that:

1. **Creates Jobs**: Each sync operation creates a job with unique ID
2. **Chunked Processing**: Processes data in batches of 100 items
3. **Progress Tracking**: Real-time progress updates
4. **Error Handling**: Comprehensive error handling with detailed messages
5. **Pagination Support**: Automatically handles external API pagination

### Usage Example

```bash
# Start data synchronization
curl -X GET http://localhost:8080/api/external/hello

# Check job status
curl -X GET http://localhost:8080/api/jobs/{job-id}
```

## 🧪 Testing

Run the test suite:

```bash
go test ./...
```

Run tests with coverage:

```bash
go test -cover ./...
```

## 🚀 Deployment

### Docker Deployment

Create a Dockerfile:

```dockerfile
FROM golang:1.19-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main cmd/api/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
EXPOSE 8080
CMD ["./main"]
```

Build and run:

```bash
docker build -t truora-microservice .
docker run -p 8080:8080 truora-microservice
```

### Production Considerations

- Use environment variables for sensitive configuration
- Implement proper logging and monitoring
- Set up health checks and readiness probes
- Configure database connection pooling
- Implement rate limiting and authentication
- Set up proper backup and disaster recovery

## 📝 API Documentation

For detailed API documentation, see [docs/API_DOCUMENTATION.md](docs/API_DOCUMENTATION.md)

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🆘 Support

For support and questions:

- Create an issue in the repository
- Check the API documentation
- Review the configuration examples

## 🔄 Changelog

### Version 1.0.0
- Initial release
- Stock rating CRUD operations
- External API integration
- Trading algorithms
- Job management system
- Comprehensive API documentation 