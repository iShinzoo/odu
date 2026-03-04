# Real-Time Order Processing System  
### Go + gRPC + HTTP Gateway + Worker Pool + WebSocket + PostgreSQL

![Go](https://img.shields.io/badge/Go-1.25-blue)
![Docker](https://img.shields.io/badge/Docker-enabled-blue)
![gRPC](https://img.shields.io/badge/gRPC-microservice-green)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-database-blue)

A production-grade backend system built in Go demonstrating:

- Clean architecture
- gRPC microservice communication
- HTTP API Gateway
- Concurrent worker pool
- Real-time updates via WebSockets
- PostgreSQL persistence
- Interceptors (logging, recovery, request tracing)
- Graceful shutdown
- Context propagation
- Dockerized microservice deployment

This project simulates a real-world backend system architecture commonly used in scalable microservice environments.

---

# Architecture Overview

```
Client (HTTP / WebSocket)
        │
        ▼
API Gateway (HTTP :8080)
        │
        ▼
gRPC Order Service (:50051)
        │
        ├── PostgreSQL
        └── Worker Pool (5 goroutines)
                    │
                    ▼
              WebSocket Hub (:8081)
                    │
                    ▼
                 Client
```

---

# Containerized Architecture (Docker)

```
                +---------------------+
                |    API Gateway      |
                |      (HTTP)         |
                |      :8080          |
                +----------+----------+
                           |
                           v
                +---------------------+
                |     Order Service   |
                |      gRPC Server    |
                |      :50051         |
                |  WebSocket :8081    |
                +----------+----------+
                           |
                           v
                +---------------------+
                |     PostgreSQL      |
                |       :5432         |
                +---------------------+
```

All services run inside a **Docker network** and communicate using **service names**.

Example:

```
api-gateway → order-service:50051
order-service → postgres:5432
```

---

# Tech Stack

| Layer | Technology |
|------|-------------|
| Language | Go |
| Transport | gRPC |
| API Layer | HTTP (Chi Router) |
| Real-Time | WebSocket (Melody) |
| Concurrency | Goroutines + Channels |
| Database | PostgreSQL |
| Logging | Zap |
| Proto | Protocol Buffers |
| Containerization | Docker + Docker Compose |

---

# Project Structure

```
order-system/
│
├── cmd/
│   ├── order-service/        # gRPC + WebSocket service
│   │   └── Dockerfile
│   │
│   ├── api-gateway/          # HTTP gateway
│   │   └── Dockerfile
│   │
│   └── main.go               # standalone application for manual testing
│
├── internal/
│   ├── config/               # Environment configuration loader
│   ├── db/                   # PostgreSQL connection + retry logic
│   ├── order/                # Domain + service + repository
│   ├── worker/               # Worker pool implementation
│   ├── ws/                   # WebSocket hub
│   └── interceptor/          # gRPC interceptors
│
├── proto/
│   ├── order.proto
│   ├── order.pb.go
│   └── order_grpc.pb.go
│
├── migrations/               # DB schema migrations
│
├── pkg/
│   └── logger/               # Zap logger setup
│
├── docker-compose.yml        # container orchestration
│
├── go.mod
└── go.sum
```

---

# System Design Concepts

## Clean Architecture

- Business logic is independent of transport layer.
- Repository pattern isolates database access.
- Services depend on abstractions rather than implementations.
- Transport layers (HTTP/gRPC) simply adapt requests to the service layer.

---

## Worker Pool Design

- Fixed 5 goroutines
- Buffered job queue
- Context-aware processing
- Graceful shutdown support

---

## Real-Time Notification Flow

```
Order Created
      │
      ▼
Worker Processes Order
      │
      ▼
Order Status Updated in DB
      │
      ▼
WebSocket Hub Publishes Event
      │
      ▼
Subscribed Clients Receive Update
```

---

# Running the System

The system can be run in **two ways**:

1. Local Development (manual services)
2. Docker Deployment (recommended)

---

# Local Development Setup

## Install Requirements

- Go 1.21+
- PostgreSQL
- Docker
- protoc
- protoc-gen-go
- protoc-gen-go-grpc

---

## Start PostgreSQL

```
docker run --name order-postgres \
-e POSTGRES_USER=admin \
-e POSTGRES_PASSWORD=admin123 \
-e POSTGRES_DB=orderdb \
-p 5432:5432 \
-d postgres
```

---

## Run Migrations

```
migrate -path migrations \
-database "postgres://admin:admin123@localhost:5432/orderdb?sslmode=disable" up
```

---

## Run Order Service

```
go run cmd/order-service/main.go
```

Runs on:

```
gRPC      :50051
WebSocket :8081
```

---

## Run API Gateway

```
go run cmd/api-gateway/main.go
```

Runs on:

```
HTTP :8080
```

---

# Docker Deployment (Recommended)

The entire system can be started with **one command using Docker Compose**.

This launches:

- PostgreSQL container
- Order Service container
- API Gateway container

All connected through a shared Docker network.

---

## Start the System

From the project root:

```
docker compose up --build
```

Docker will:

1. Build Go binaries
2. Create containers
3. Start PostgreSQL
4. Start order-service
5. Start API gateway

---

## Stop the System

```
docker compose down
```

---

## View Running Containers

```
docker ps
```

Expected services:

```
order-postgres
order-service
api-gateway
```

---

# API Endpoints

## Create Order

POST

```
http://localhost:8080/orders
```

Body:

```json
{
  "user_id": "user-1",
  "amount": 100
}
```

---

## Get Order

GET

```
http://localhost:8080/orders/{order_id}
```

---

# WebSocket Testing (Postman)

1. Open Postman  
2. Create WebSocket request  

Connect to:

```
ws://localhost:8081/ws
```

Send the Order ID:

```
"ORDER_ID_HERE"
```

When the worker processes the order you will receive:

```json
{
  "order_id": "...",
  "status": "PROCESSED"
}
```

---

# Production Features

## gRPC Interceptors

- Logging interceptor
- Panic recovery interceptor
- Request ID tracing

---

## Context Propagation

- Database queries use context
- Worker pool supports cancellation
- HTTP requests propagate context to gRPC
- gRPC calls use timeouts

---

## Graceful Shutdown

The system safely handles:

- SIGINT / SIGTERM
- Worker shutdown
- gRPC GracefulStop
- Database connection close

---

# UML DIAGRAMS

---

## Component Diagram

```
+----------------+
|   API Gateway  |
+----------------+
        |
        v
+----------------+
|  gRPC Service  |
+----------------+
    |        |
    v        v
+--------+  +-------------+
| Worker |  | PostgreSQL  |
+--------+  +-------------+
    |
    v
+--------------+
| WebSocket Hub|
+--------------+
```

---

## Class Diagram

```
+------------------+
| OrderService     |
|------------------|
| +CreateOrder()   |
| +GetOrder()      |
| +UpdateStatus()  |
+------------------+
         |
         v
+----------------------+
| OrderRepository      |
|----------------------|
| +Create()            |
| +GetByID()           |
| +UpdateStatus()      |
+----------------------+
         |
         v
+----------------------------+
| PostgresOrderRepository    |
+----------------------------+
```

---

## Sequence Diagram

```
Client → API Gateway → gRPC Service → Database
                                     ↓
                                  Worker Pool
                                     ↓
                               WebSocket Hub
                                     ↓
                                  Client
```

---

## Concurrency Model

```
                +------------------+
                | Job Channel      |
                +------------------+
                 |   |   |   |   |
                 v   v   v   v   v
                 W1  W2  W3  W4  W5
```

---

# Scalability

Current system:

- Single instance
- In-memory job queue
- In-memory WebSocket subscriptions

Scaling strategy:

- Redis / Kafka job queue
- Redis Pub/Sub for WebSocket notifications
- Horizontal scaling behind load balancer

---

# Future Improvements

- Prometheus metrics
- OpenTelemetry tracing
- JWT authentication
- Kubernetes deployment
- Distributed worker queue
- Redis Pub/Sub for notifications

---

# Author

Built as a backend systems learning project demonstrating production-ready architecture using Go.

---

# Contact

X: https://x.com/i_krsna4  
LinkedIn: https://www.linkedin.com/in/krishnathakur1/
