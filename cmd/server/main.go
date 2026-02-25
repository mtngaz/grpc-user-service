package main

import (
	"context"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	pb "github.com/mtngaz/grpc-user-service/api"
	"github.com/mtngaz/grpc-user-service/internal/service"
	"github.com/mtngaz/grpc-user-service/internal/storage"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	http.Handle("/swagger/",
	http.StripPrefix("/swagger/",
		http.FileServer(http.Dir("./swagger")),
	),
)

	store := storage.NewRedisStore("localhost:6379")
	userService := service.NewUserService(store)

	grpcLis, _ := net.Listen("tcp", ":50051")
	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, userService)

	ctx := context.Background()
	mux := runtime.NewServeMux()

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	if err := pb.RegisterUserServiceHandlerFromEndpoint(ctx, mux, ":50051", opts); err != nil {
		log.Fatalf("failed to start gateway: %v", err)
	}

	go func() {
		log.Println("gRPC server running on :50051")
		grpcServer.Serve(grpcLis)
	}()

	log.Println("HTTP gateway running on :8080")
	http.ListenAndServe(":8080", mux)
}
