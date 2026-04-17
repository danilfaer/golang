package part

import (
	"context"
	"log"

	partUtils "github.com/danilfaer/golang/inventory/internal/repository/utils"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type repository struct {
	coll *mongo.Collection
}

// NewRepository создаёт репозиторий поверх коллекции MongoDB.
func NewRepository(coll *mongo.Collection) *repository {

	setupCtx := context.Background()
	if err := EnsurePartIndexes(setupCtx, coll); err != nil {
		log.Fatalf("ошибка создания индексов: %v", err)
	}
	if err := partUtils.SeedPartsIfEmpty(setupCtx, coll); err != nil {
		log.Fatalf("ошибка сида коллекции parts: %v", err)
	}

	return &repository{coll: coll}
}
