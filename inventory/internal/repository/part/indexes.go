package part

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// EnsurePartIndexes создаёт индексы для типовых запросов ListParts (идемпотентно).
func EnsurePartIndexes(ctx context.Context, coll *mongo.Collection) error {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// Keys только bson.D — map (bson.M) драйвером для индексов не допускается.
	idx := []mongo.IndexModel{
		{Keys: bson.D{{Key: "category", Value: 1}}},
		{Keys: bson.D{{Key: "manufacturer.country", Value: 1}}},
		{Keys: bson.D{{Key: "tags", Value: 1}}},
		{Keys: bson.D{{Key: "name", Value: 1}}},
	}
	_, err := coll.Indexes().CreateMany(ctx, idx)
	return err
}
