package main

import (
	"log"
	"net"

	"go-sample-http/config"
	pb "go-sample-http/gen/item"
	"go-sample-http/internal/db"
	"go-sample-http/internal/item"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	cfg := config.Load()

	database, err := db.New(cfg.DBPath)
	if err != nil {
		log.Fatal("db init failed:", err)
	}
	defer database.Close()

	// HTTP 서버와 동일한 레이어 조립
	itemRepo := item.NewRepository(database)
	itemUseCase := item.NewUseCase(itemRepo)
	itemHandler := item.NewGRPCHandler(itemUseCase)

	lis, err := net.Listen("tcp", ":"+cfg.GRPCPort)
	if err != nil {
		log.Fatal("failed to listen:", err)
	}

	s := grpc.NewServer()
	pb.RegisterItemServiceServer(s, itemHandler)
	reflection.Register(s) // grpcurl 같은 도구로 탐색 가능

	log.Printf("gRPC server listening on :%s", cfg.GRPCPort)
	if err := s.Serve(lis); err != nil {
		log.Fatal("failed to serve:", err)
	}
}
