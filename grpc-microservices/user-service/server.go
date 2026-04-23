package main

import (
	"context"
	"fmt"
	"sync"
	"time"
	userpb "user-service/gen/user"
)

//------------------This code is creating a gRPC server that manages users in memory ----------------------------------

type User struct{
	UserId string
	UserName string
	UserEmail string
}

type UserServer struct{
	userpb.UnimplementedUserServiceServer
	mu sync.RWMutex
	user map[string]*User
}

func NewUserServer()*UserServer{
  return &UserServer{
    user: make(map[string]*User),
  }
}

//==========================================================================

func (s *UserServer)CreateUser(_ context.Context,req *userpb.CreateUserRequest)(*userpb.CreateUserResponse,error){
	
	userID:= fmt.Sprintf("user-%d",time.Now().UnixNano())
	user:=&User{
		UserId:userID,
		UserName:req.Username,
		UserEmail:req.Email,
	}
	s.mu.Lock()
	s.user[userID] = user
	s.mu.Unlock()

	return &userpb.CreateUserResponse{
		UserId: userID,
		Username: user.UserName,
	},nil
}

func(s *UserServer)GetUser(_ context.Context,req *userpb.GetUserRequest)(*userpb.GetUserResponse,error){

	s.mu.RLock()
	user,ok:= s.user[req.UserId]
	s.mu.RUnlock()

	if !ok{
		return nil,fmt.Errorf("user %s not found", req.UserId)
	}
	return &userpb.GetUserResponse{
		UserId: user.UserId,
		Username: user.UserName,
		Email: user.UserEmail,
	},nil
}
 
func(s *UserServer)ValidateUser(_ context.Context,req *userpb.ValidateUserRequest)(*userpb.ValidateUserResponse,error){

	s.mu.RLock()
	user,ok:= s.user[req.UserId]
	s.mu.RUnlock()

	if !ok{
		return &userpb.ValidateUserResponse{Valid:false},nil
	}

	return &userpb.ValidateUserResponse{
		Valid:true,
		Username:user.UserName,
		
	},nil
}
