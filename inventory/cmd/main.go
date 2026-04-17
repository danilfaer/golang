package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	inventoryV1API "github.com/danilfaer/golang/inventory/internal/api/inventory/v1"
	inventoryRepository "github.com/danilfaer/golang/inventory/internal/repository/part"
	inventoryService "github.com/danilfaer/golang/inventory/internal/service/part"
	sharedInterceptors "github.com/danilfaer/golang/shared/pkg/interceptors"
	inventoryV1 "github.com/danilfaer/golang/shared/pkg/proto/inventory/v1"
)

const grpsPort = 50051

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("предупреждение: не удалось загрузить .env: %v", err)
	}

	uri := os.Getenv("MONGO_DB_URI")
	if uri == "" {
		log.Fatal("MONGO_DB_URI не задан — без MongoDB сервис не стартует")
	}

	mongoClient, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("ошибка подключения к MongoDB: %v", err)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := mongoClient.Disconnect(ctx); err != nil {
			log.Printf("ошибка отключения от MongoDB: %v", err)
		}
	}()

	pingCtx, pingCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer pingCancel()
	if err := mongoClient.Ping(pingCtx, nil); err != nil {
		log.Fatalf("MongoDB недоступна (ping): %v", err)
	}
	log.Println("соединение с MongoDB установлено")

	dbName := os.Getenv("MONGO_INVENTORY_DB")
	if dbName == "" {
		dbName = "inventory-service"
	}
	partsColl := mongoClient.Database(dbName).Collection("parts")

	
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpsPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	defer func() {
		if cerr := lis.Close(); cerr != nil {
			log.Fatalf("failed to close listener: %v", cerr)
		}
	}()

	s := grpc.NewServer(
		grpc.UnaryInterceptor(sharedInterceptors.UnaryErrorInterceptor()),
	)

	repo := inventoryRepository.NewRepository(partsColl)
	service := inventoryService.NewService(repo)
	api := inventoryV1API.NewAPI(service)

	inventoryV1.RegisterInventoryServiceServer(s, api)
	reflection.Register(s)

	go func() {
		log.Printf("gRPS inventory listening on %d", grpsPort)
		err := s.Serve(lis)
		if err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down gRPC server...")
	s.GracefulStop()
	log.Println("gRPC server stopped")
}
