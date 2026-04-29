package app

import (
	"context"

	inventoryV1API "github.com/danilfaer/golang/inventory/internal/api/inventory/v1"
	"github.com/danilfaer/golang/inventory/internal/config"
	"github.com/danilfaer/golang/inventory/internal/repository"
	inventoryRepository "github.com/danilfaer/golang/inventory/internal/repository/part"
	"github.com/danilfaer/golang/inventory/internal/service"
	inventoryService "github.com/danilfaer/golang/inventory/internal/service/part"
	inventoryV1 "github.com/danilfaer/golang/shared/pkg/proto/inventory/v1"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type diContainer struct {
	inventoryGRPCImpl inventoryV1.InventoryServiceServer

	inventoryService service.InventoryService

	partRepository repository.InventoryRepository

	mongoClient *mongo.Client
}

func newDIContainer() *diContainer {
	return &diContainer{}
}

// InitMongo подключает клиент MongoDB один раз; проверка доступности — в App (ping).
func (d *diContainer) InitMongo(ctx context.Context) error {
	if d.mongoClient != nil {
		return nil
	}
	client, err := mongo.Connect(options.Client().ApplyURI(config.AppConfig().Mongo.URI()))
	if err != nil {
		return err
	}
	d.mongoClient = client
	return nil
}

// MongoClient возвращает клиент после успешного InitMongo.
func (d *diContainer) MongoClient() *mongo.Client {
	return d.mongoClient
}

func (d *diContainer) partsCollection() *mongo.Collection {
	dbName := config.AppConfig().Mongo.Database()
	return d.mongoClient.Database(dbName).Collection("parts")
}

func (d *diContainer) PartRepository(ctx context.Context) repository.InventoryRepository {
	if d.partRepository == nil {
		d.partRepository = inventoryRepository.NewRepository(d.partsCollection())
	}
	return d.partRepository
}

func (d *diContainer) InventoryService(ctx context.Context) service.InventoryService {
	if d.inventoryService == nil {
		d.inventoryService = inventoryService.NewService(d.PartRepository(ctx))
	}
	return d.inventoryService
}

func (d *diContainer) InventoryGRPCAPI(ctx context.Context) inventoryV1.InventoryServiceServer {
	if d.inventoryGRPCImpl == nil {
		d.inventoryGRPCImpl = inventoryV1API.NewAPI(d.InventoryService(ctx))
	}
	return d.inventoryGRPCImpl
}
