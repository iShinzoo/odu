# Real-Time Order Processing System  
### Go + gRPC + HTTP Gateway + Worker Pool + WebSocket + PostgreSQL

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

This project simulates a real-world microservice architecture suitable for scalable backend systems.

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

# Tech Stack

| Layer | Technology |
|--------|------------|
| Language | Go |
| Transport | gRPC |
| API Layer | HTTP (Chi Router) |
| Real-Time | WebSocket (Melody) |
| Concurrency | Goroutines + Channels |
| Database | PostgreSQL |
| Logging | Zap |
| Proto | Protocol Buffers |
| Architecture | Clean Architecture |

---

# Project Structure

```
order-system/
│
├── cmd/
│   ├── order-service/        # gRPC server
│   └── api-gateway/          # HTTP gateway
|   └── main.go               # standalone application for manual testing purpose
│
├── internal/
│   ├── config/               # Environment config loader
│   ├── db/                   # PostgreSQL connection
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
├── go.mod
└── go.sum
```

---

# System Design Concepts

## 1️⃣ Clean Architecture

- Business logic is independent of transport.
- Repository is an interface.
- Service depends only on abstraction.
- Transport (gRPC / HTTP) adapts to service layer.

---

## 2️⃣ Worker Pool Design

- Fixed 5 goroutines.
- Buffered channel for job queue.
- Context cancellation supported.
- Async processing of orders.

---

## 3️⃣ Real-Time Notification Flow

```
Order Created
      │
      ▼
Worker Processes
      │
      ▼
Status Updated in DB
      │
      ▼
WebSocket Hub Publishes Event
      │
      ▼
Subscribed Client Receives Update
```

---

# Setup Instructions

## 1️⃣ Install Requirements

- Go (1.21+)
- Docker
- PostgreSQL (via Docker)
- protoc
- protoc-gen-go
- protoc-gen-go-grpc

---

## 2️⃣ Start PostgreSQL (Docker)

```bash
docker run --name order-postgres \
-e POSTGRES_USER=admin \
-e POSTGRES_PASSWORD=admin123 \
-e POSTGRES_DB=orderdb \
-p 5432:5432 \
-d postgres
```

---

## 3️⃣ Run Migrations

```bash
migrate -path migrations \
-database "postgres://admin:admin123@localhost:5432/orderdb?sslmode=disable" up
```

---

## 4️⃣ Run gRPC Order Service

```bash
go run cmd/order-service/main.go
```

Runs on:
- gRPC → `:50051`
- WebSocket → `:8081`

---

## 5️⃣ Run API Gateway

```bash
go run cmd/api-gateway/main.go
```

Runs on:
- HTTP → `:8080`

---

# API Endpoints

## Create Order

POST:
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

GET:
```
http://localhost:8080/orders/{order_id}
```

---

# WebSocket Testing (Postman)

1. Open WebSocket in Postman  
2. Connect to:

```
ws://localhost:8081/ws
```

3. Send Order ID as message:

```
"ORDER_ID_HERE"
```

4. When worker processes order → real-time update received:

```json
{
  "order_id": "...",
  "status": "PROCESSED"
}
```

---

# Production Features Implemented

## gRPC Interceptors

- Logging interceptor
- Panic recovery interceptor
- Request ID tracing

---

## Context Propagation

- DB queries use `ExecContext`
- Worker supports cancellation
- gRPC calls have timeouts
- HTTP propagates request context

---

## Graceful Shutdown

Handles:
- SIGINT
- Worker cancellation
- gRPC GracefulStop()
- DB connection close

---

# UML DIAGRAMS

---

## 1️⃣ Component Diagram

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

## 2️⃣ Class Diagram (Domain Layer)

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

## 3️⃣ Sequence Diagram — Order Creation

```
Client → API Gateway → gRPC Service → DB
                                     ↓
                                 Worker Pool
                                     ↓
                               Update Status
                                     ↓
                               WebSocket Hub
                                     ↓
                                  Client
```

---

## 4️⃣ Concurrency Model

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
- In-memory WebSocket subscription map

To scale horizontally:
- Replace worker queue with Redis/Kafka
- Replace WebSocket hub with Redis Pub/Sub
- Deploy behind load balancer

---


# Future Improvements

- Prometheus metrics
- OpenTelemetry tracing
- JWT authentication
- Docker Compose full stack
- Kubernetes deployment
- Redis-backed distributed worker queue

---

# Author

Built as a backend systems learning project demonstrating production-ready architecture using Go.

---

# Contact
   - X: [@i_krsna4](https://x.com/i_krsna4)
   - LinkedIn: [@krishnathakur1](https://www.linkedin.com/in/krishnathakur1/)
