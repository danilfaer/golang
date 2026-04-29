package app

import (
	"context"
	"fmt"

	inventoryV1API "github.com/danilfaer/golang/inventory/internal/api/inventory/v1"
	"github.com/danilfaer/golang/inventory/internal/config"
	inventoryRepository "github.com/danilfaer/golang/inventory/internal/repository"
	partRepository "github.com/danilfaer/golang/inventory/internal/repository/part"
	inventoryService "github.com/danilfaer/golang/inventory/internal/service"
	partService "github.com/danilfaer/golang/inventory/internal/service/part"
	"github.com/danilfaer/golang/platform/pkg/closer"
	inventoryV1 "github.com/danilfaer/golang/shared/pkg/proto/inventory/v1"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

type diContainer struct {
	inventoryV1API inventoryV1.InventoryServiceServer

	inventoryService inventoryService.InventoryService

	inventoryRepository inventoryRepository.InventoryRepository

	mongoDBClient *mongo.Client
	mongoDBHandle *mongo.Database
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) InventoryV1API(ctx context.Context) inventoryV1.InventoryServiceServer {
	if d.inventoryV1API == nil {
		d.inventoryV1API = inventoryV1API.NewAPI(d.InventoryService(ctx))
	}
	return d.inventoryV1API
}

func (d *diContainer) InventoryService(ctx context.Context) inventoryService.InventoryService {
	if d.inventoryService == nil {
		d.inventoryService = partService.NewService(d.InventoryRepository(ctx))
	}
	return d.inventoryService
}

func (d *diContainer) InventoryRepository(ctx context.Context) inventoryRepository.InventoryRepository {
	if d.inventoryRepository == nil {
		// Коллекция parts — та же, что в сидах и индексах в part.NewRepository.
		d.inventoryRepository = partRepository.NewRepository(d.MongoDBHandle(ctx).Collection("parts"))
	}
	return d.inventoryRepository
}

func (d *diContainer) MongoDBClient(ctx context.Context) *mongo.Client {
	if d.mongoDBClient == nil {
		client, err := mongo.Connect(options.Client().ApplyURI(config.AppConfig().Mongo.URI()))
		if err != nil {
			panic(fmt.Sprintf("failed to connect to MongoDB: %s\n", err.Error()))
		}
		err = client.Ping(ctx, readpref.Primary())
		if err != nil {
			panic(fmt.Sprintf("failed to ping MongoDB: %s\n", err.Error()))
		}
		closer.AddNamed("MongoDB Client", func(ctx context.Context) error {
			return client.Disconnect(ctx)
		})
		d.mongoDBClient = client
	}
	return d.mongoDBClient
}

func (d *diContainer) MongoDBHandle(ctx context.Context) *mongo.Database {
	if d.mongoDBHandle == nil {
		d.mongoDBHandle = d.MongoDBClient(ctx).Database(config.AppConfig().Mongo.Database())
	}
	return d.mongoDBHandle
}
