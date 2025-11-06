# oolio-kart-challenge

A RESTful API  Server in Golang for ordering Products, featuring product management, order processing, and promotional code support.

## Features

- **Product Management**: List and retrieve electronic products (phones, tablets, laptops)
- **Order Processing**: Place orders with multiple items and promotional codes
- **Authentication**: API key-based authentication for order endpoints
- **Promotional Codes**: Support for discount coupons loaded from text files
- **CORS Support**: Cross-origin resource sharing enabled
- **Structured Logging**: Comprehensive request/response logging
- **Clean Architecture**: Domain-driven design with clear separation of concerns

## Prerequisites

- **Go 1.22+** (required for new HTTP routing features)
- **Git** for version control

## Project Structure

```
go-food-ordering-api/
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
├── internal/
│   ├── application/
│   │   └── services/            # Business logic layer
│   ├── domain/
│   │   ├── entities/            # Domain models
│   │   ├── errors/              # Domain errors
│   │   └── interfaces/          # Repository and service interfaces
│   └── infrastructure/
│       ├── http/
│       │   ├── handlers/        # HTTP request handlers
│       │   ├── middleware/      # HTTP middleware
│       │   └── router.go        # Route configuration
│       └── repositories/        # Data access layer
├── pkg/
│   └── logger/                  # Logging utilities
├── testdata/                    # Test coupon files
├── go.mod                       # Go module definition
└── README.md                    # This file
```

## Quick Start

### 1. Clone the Repository

```bash
git clone <repository-url>
cd go-food-ordering-api
```

### 2. Install Dependencies

```bash
go mod download
```

### 3. Set Environment Variables

Create a `.env` file or set environment variables:

```bash
# Server configuration
export PORT=8080
export API_KEY=your-secret-api-key

# Coupon files (optional - defaults provided)
export COUPON_FILES=testdata/couponbase1.txt,testdata/couponbase2.txt,testdata/couponbase3.txt
```

### 4. Run the Application

```bash
# Development mode
go run cmd/server/main.go

# Or build and run
go build -o server cmd/server/main.go
./serv

# Integration test 
go test ./internal -v -run TestOpenAPICompliance

# Common HTTP status codes:
- `200` - Success
- `400` - Bad Request (validation errors)
- `401` - Unauthorized (missing/invalid API key)
- `404` - Not Found (product not found)
- `500` - Internal Server Error
