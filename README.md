# Order Service

## Overview
The Order Service is a microservice designed to handle order management in a distributed system. It's responsible for creating and managing orders, integrating with a Product Service, and implementing various backend functionalities using modern cloud-native architecture patterns.

## Technology Stack
- **Programming Language:** Go (Golang)
- **Web Framework:** Gin
- **Database:** PostgreSQL 15
- **Caching:** Redis 7
- **Message Broker:** RabbitMQ 3
- **ORM:** GORM
- **API Documentation:** OpenAPI/Swagger
- **Containerization:** Docker & Docker Compose
- **Database Migration:** golang-migrate

## Prerequisites
- Docker and Docker Compose
- Go 1.x (for local development)
- Make (optional, for using Makefile commands)

## How to Run

1. Clone the repository:
```bash
git clone <repository-url>
cd order-service
```

2. Start the services using Docker Compose:
```bash
docker-compose up -d
```

This will start:
- Order Service (port 8000)
- PostgreSQL (internal port 5432)
- Redis (internal port 6379)
- RabbitMQ (ports 5672, 15672)
- Database migrations will run automatically

## API Overview

### Endpoints

#### Create Order
- **Method:** POST
- **Path:** `/orders`
- **Description:** Create a new order
- **Request Body:**
```json
{
    "productId": "uuid-string",
    "quantity": 1
}
```
- **Response:** 201 Created
```json
{
    "id": "order-uuid",
    "productId": "product-uuid",
    "totalPrice": 99.99,
    "status": "pending",
    "createdAt": "2025-11-02T12:00:00Z"
}
```

#### Get Orders by Product ID
- **Method:** GET
- **Path:** `/orders/product/:productId`
- **Description:** Retrieve all orders for a specific product
- **Response:** 200 OK
```json
[
    {
        "id": "order-uuid",
        "productId": "product-uuid",
        "totalPrice": 99.99,
        "status": "pending",
        "createdAt": "2025-11-02T12:00:00Z"
    }
]
```

## Architecture

### Components
1. **API Layer (Handlers)**
   - Handles HTTP requests
   - Input validation
   - Response formatting

2. **Service Layer**
   - Business logic
   - Integration with external services
   - Transaction management

3. **Repository Layer**
   - Database operations
   - Data access patterns
   - CRUD operations

4. **External Services Integration**
   - Product Service communication
   - RabbitMQ message publishing
   - Redis caching

### Data Flow
1. Client makes HTTP request to Order Service
2. Request is validated by the handler
3. Service layer processes the business logic
4. Product information is fetched from Product Service
5. Order is stored in PostgreSQL
6. Cache is updated in Redis
7. Event is published to RabbitMQ
8. Response is sent back to client

### Infrastructure
- **PostgreSQL:** Primary data store for orders
- **Redis:** Caching layer for frequent queries
- **RabbitMQ:** Message broker for event-driven architecture
- **Docker:** Containerization for consistent development and deployment
- **Docker Compose:** Local development environment orchestration

## Environment Variables
```
APP_PORT=8000
PRODUCT_SERVICE_URL=http://product-service-app:3000
MOCK_PRODUCT_SERVICE=true
DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=order_db
REDIS_HOST=redis
REDIS_PORT=6379
RABBITMQ_URL=amqp://guest:guest@rabbitmq:5672/
```

## Testing
The service includes both unit tests and integration tests. Mock implementations are provided for external service dependencies to facilitate testing.

---
For more detailed information about specific components or development guidelines, please refer to the documentation in the `/docs` directory.