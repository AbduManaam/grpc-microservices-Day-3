package main

import (
  "context"
  "fmt"

  "google.golang.org/grpc"
  "google.golang.org/grpc/credentials/insecure"

  // generated from user.proto
  userpb "order-service/gen/user"
)

// UserClient wraps the gRPC stub
type UserClient struct {              //created this struct to wrap the gRPC stub and provide a clean,
  client userpb.UserServiceClient     // controlled, and maintainable interface for calling the UserService.
}


func NewUserClient(addr string) (*UserClient, func(), error) {
  conn, err := grpc.Dial(
    addr,                                                      //opens a network
    grpc.WithTransportCredentials(insecure.NewCredentials()),  //“Create a connection to server using plain TCP, WITHOUT encryption (no TLS) [secure gRPC Connection]”
  )
  if err != nil {
    return nil, nil, fmt.Errorf("dial UserService: %w", err)
  }

  cleanup := func() { conn.Close() }
  stub := userpb.NewUserServiceClient(conn)

  return &UserClient{client: stub}, cleanup, nil
}

// ValidateUser calls UserService.ValidateUser via gRPC
func (u *UserClient) ValidateUser(ctx context.Context, userID string) (bool, string, error) {
  resp, err := u.client.ValidateUser(ctx, &userpb.ValidateUserRequest{
    UserId: userID,
  })
  if err != nil {
    return false, "", fmt.Errorf("ValidateUser rpc: %w", err)
  }
  return resp.Valid, resp.Username, nil
}