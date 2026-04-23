package main

import (
  "log"
  "net"

  "google.golang.org/grpc"
  "google.golang.org/grpc/reflection"
  orderpb "order-service/gen/order"
)

func main() {
  // 1. Create UserService gRPC client
  userAddr := "localhost:50051"
  userClient, cleanup, err := NewUserClient(userAddr)
  if err != nil {
    log.Fatalf("connect to UserService: %v", err)
  }
  defer cleanup()

  log.Printf("Connected to UserService at %s", userAddr)

  // 2. Create OrderService server (inject user client)
  srv := grpc.NewServer()
  orderpb.RegisterOrderServiceServer(srv, NewOrderServer(userClient))
  reflection.Register(srv) 

  // 3. Listen and serve
  lis, err := net.Listen("tcp", ":50052")
  if err != nil {
    log.Fatalf("listen: %v", err)
  }

  log.Println("OrderService listening on :50052")
  if err := srv.Serve(lis); err != nil {
    log.Fatalf("serve: %v", err)
  }
}