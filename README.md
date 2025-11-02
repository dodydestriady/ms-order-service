# Order Service
The Order Service is a microservice designed to handle order management in a distributed system. It's responsible for creating and managing orders, integrating with a Product Service, and implementing various backend functionalities using modern cloud-native architecture patterns.

## Technology Stack
- **Programming Language:** Go (Golang)
- **Web Framework:** Gin
- **Database:** PostgreSQL 15
- **Caching:** Redis 7
- **Message Broker:** RabbitMQ 3
- **ORM:** GORM
- **Containerization:** Docker & Docker Compose
- **Database Migration:** golang-migrate

## Prerequisites
- Docker and Docker Compose
- Go 1.x (for local development)
- Make (optional, for using Makefile commands)

## How to Run
You can run this service in two ways:

1. All services
Visit the [runner repo](https://github.com/dodydestriady/ms-runner) to run all the services 

2. Standalone Mode
Run the Product Service only:
```
docker compose up -d --build
```
Makesure you run this if you want to run other services, and uncomment the networks shared-net in compose
```
docker network create shared-net
```

This will start:
- Order Service (port 8000)
- PostgreSQL (internal port 5432)
- Redis (internal port 6379)
- RabbitMQ (ports 5672, 15672)
- Database migrations will run automatically

## API Overview

### Endpoints
| Method | Endpoint | Description|
| :---: | :---:| :---:|
|POST|/orders|Create an order|
|GET|/orders/product/:productId|Show ordeer by product|

### Example Requests
Create Order
```
curl -X POST http://localhost:8000/orders \  
-H "Content-Type: application/json" \
-d '{"product_id": "1", "quantity": 100}'
```

Show Order by Product
```
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
