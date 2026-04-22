# gRPC Microservices Project

This project demonstrates a microservices architecture using gRPC with two services: UserService and OrderService.

## Project Structure

```
grpc-microservices/
├── proto/
│   ├── user.proto       # User service definitions
│   └── order.proto      # Order service definitions
├── user-service/
│   ├── main.go          # Entry point for user service
│   ├── server.go        # User service implementation
│   └── gen/user/        # Generated protobuf code (generated)
├── order-service/
│   ├── main.go          # Entry point for order service
│   ├── server.go        # Order service implementation
│   ├── user_client.go   # gRPC client for UserService
│   └── gen/
│       ├── order/       # Generated from order.proto
│       └── user/        # Generated from user.proto
├── go.mod              # Go module file
└── README.md           # This file
```

## Services

### UserService (Port: 50051)
- **GetUser**: Retrieve a user by ID
- **CreateUser**: Create a new user
- **GetAllUsers**: Retrieve all users

### OrderService (Port: 50052)
- **GetOrder**: Retrieve an order by ID
- **CreateOrder**: Create a new order (validates user exists)
- **GetUserOrders**: Retrieve all orders for a user
- **UpdateOrderStatus**: Update order status

## Setup

### Prerequisites
- Go 1.21 or higher
- protoc (Protocol Buffer compiler)
- protoc-gen-go
- protoc-gen-go-grpc

### Installation

1. Install Go plugins:
```bash
go install github.com/golang/protobuf/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

2. Generate protobuf code:
```bash
# From the root directory
protoc --go_out=. --go-grpc_out=. proto/*.proto
```

3. Install dependencies:
```bash
go mod download
```

## Running the Services

### Terminal 1 - Start User Service
```bash
cd user-service
go run main.go server.go
# Output: User Service listening on :50051
```

### Terminal 2 - Start Order Service
```bash
cd order-service
go run main.go server.go user_client.go
# Output: Order Service listening on :50052
```

## Usage Example

You can use tools like `grpcurl` to test the services:

```bash
# Create a user
grpcurl -plaintext -d '{"name":"John Doe","email":"john@example.com","phone":"1234567890"}' \
  localhost:50051 user.UserService/CreateUser

# Get all users
grpcurl -plaintext localhost:50051 user.UserService/GetAllUsers

# Create an order (user_id should exist from created user)
grpcurl -plaintext -d '{"user_id":1,"product_name":"Laptop","quantity":1,"price":999.99}' \
  localhost:50052 order.OrderService/CreateOrder

# Get user orders
grpcurl -plaintext -d '{"user_id":1}' \
  localhost:50052 order.OrderService/GetUserOrders
```

## Key Features

- **Service-to-Service Communication**: OrderService calls UserService to validate user existence
- **Thread-Safe**: Uses mutex locks for concurrent access
- **Error Handling**: Proper error handling for missing resources
- **Scalable Architecture**: Each service runs independently on different ports

## Notes

- User and Order IDs are auto-incremented starting from 1
- Orders start with PENDING status
- The generated code directories (gen/) are not included in version control