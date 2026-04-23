package main

import (
  "context"
  "fmt"
  "sync"
  "time"

  orderpb "order-service/gen/order"
)

// Its job is to create and fetch orders, and before creating an order, it checks with UserService if the user is valid.

type Order struct {
  OrderID   string
  UserID    string
  ProductID string
  Quantity  int32
  Status    string
}

type OrderServer struct {
  orderpb.UnimplementedOrderServiceServer
  userClient *UserClient   // injected gRPC client
  mu         sync.RWMutex
  orders     map[string]*Order
}

func NewOrderServer(uc *UserClient) *OrderServer {
  return &OrderServer{
    userClient: uc,
    orders:     make(map[string]*Order),  // Initializes in-memory storage for orders,Empty map ready to store data
  }
}

//  CreateOrder
func (s *OrderServer) CreateOrder(
  ctx context.Context,
  req *orderpb.CreateOrderRequest,
) (*orderpb.CreateOrderResponse, error) {

  valid, username, err := s.userClient.ValidateUser(ctx, req.UserId)
  if err != nil {
    return nil, fmt.Errorf("user validation failed: %w", err)
  }
  if !valid {
    return &orderpb.CreateOrderResponse{
      Status: "failed: user not found",
    }, nil
  }

  //  Create order
  orderID := fmt.Sprintf("ord-%d", time.Now().UnixNano())
  order := &Order{
    OrderID:   orderID,
    UserID:    req.UserId,
    ProductID: req.ProductId,
    Quantity:  req.Quantity,
    Status:    "created",
  }

  s.mu.Lock()
  s.orders[orderID] = order
  s.mu.Unlock()

  _ = username // could log: "Order for " + username

  return &orderpb.CreateOrderResponse{
    OrderId: orderID,
    Status:  "created",
  }, nil
}

// GetOrder 
func (s *OrderServer) GetOrder(
  ctx context.Context,
  req *orderpb.GetOrderRequest,
) (*orderpb.GetOrderResponse, error) {

  s.mu.RLock()
  order, ok := s.orders[req.OrderId]
  s.mu.RUnlock()

  if !ok {
    return nil, fmt.Errorf("order %s not found", req.OrderId)
  }

  return &orderpb.GetOrderResponse{
    OrderId:   order.OrderID,
    UserId:    order.UserID,
    ProductId: order.ProductID,
    Quantity:  order.Quantity,
    Status:    order.Status,
  }, nil
}