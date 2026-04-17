package part

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/danilfaer/golang/inventory/internal/model"
	repoModel "github.com/danilfaer/golang/inventory/internal/repository/model"
)

func (r *repository) GetPart(ctx context.Context, uuid string) (*repoModel.Part, error) {
	var doc repoModel.Part
	err := r.coll.FindOne(ctx, bson.M{"_id": uuid}).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, model.ErrPartNotFound
		}
		return nil, err
	}
	return &doc, nil
}
