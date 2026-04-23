package main

import (
	"log"
	"net"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"


	userpb "user-service/gen/user"
)

func main(){

   srv:=grpc.NewServer()
   
   userpb.RegisterUserServiceServer(srv,NewUserServer())
       reflection.Register(srv)              // ← used here

   
   
   lis,err:= net.Listen("tcp",":50051")  //Just opens/creates the port — but no requests accepted yet

   if err!=nil{
	log.Fatalf("Listen: %v",err)
   }

   log.Println("Userservice listening on :50051")
   if err:=srv.Serve(lis);err!=nil{   //Hands the port to gRPC server — NOW requests start coming in
	log.Fatalf("serve: %v",err)
   }

}